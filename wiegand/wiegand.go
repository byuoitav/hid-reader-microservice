package wiegand

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

// Reader represents the configuration necessary to watch a card reader
type Reader struct {
	Data0Pin    int
	Data1Pin    int
	BufferSize  int
	Notifier    chan string
	lastPulse   time.Time
	buffer      []int
	bitsCounted int
}

// Setup sets up all of the necessary pieces and interrupts to listen to the
// configured reader
func (r *Reader) Setup() error {

	// Initialize host drivers
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("Error while trying to initialize GPIO pins: %w", err)
	}

	// Initialize the buffer
	r.lastPulse = time.Now()
	r.buffer = make([]int, r.BufferSize)
	r.bitsCounted = 0

	// Setup the pins
	data0 := gpioreg.ByName("14")
	data1 := gpioreg.ByName("15")

	err = data0.In(gpio.PullUp, gpio.FallingEdge)
	if err != nil {
		return fmt.Errorf("Error while setting up data0 pin: %w", err)
	}
	err = data1.In(gpio.PullUp, gpio.FallingEdge)
	if err != nil {
		return fmt.Errorf("Error while setting up data1 pin: %w", err)
	}

	// Watch for falling edge on Data0 Pin
	go func() {
		for {
			data0.WaitForEdge(-1)
			r.lastPulse = time.Now()
			r.buffer[r.bitsCounted%r.BufferSize] = 0
			//r.buffer[r.bitsCounted] = 0
			r.bitsCounted++
		}
	}()

	// Watch for falling edge on Data1 Pin
	go func() {
		for {
			data1.WaitForEdge(-1)
			r.lastPulse = time.Now()
			r.buffer[r.bitsCounted%r.BufferSize] = 1
			//r.buffer[r.bitsCounted] = 1
			r.bitsCounted++
		}
	}()

	// Start watching the reader
	go r.watchForCard()

	return nil
}

// watchForCard continuously watches to see if a card read has finished and when
// it has then it will send the the card binary down the notifier channel
func (r *Reader) watchForCard() {
	for {
		if time.Now().Sub(r.lastPulse) > (50*time.Millisecond) && r.bitsCounted > 0 {

			// Copy the buffer to send across the channel
			// binCopy := make([]int, r.BufferSize)
			// copy(binCopy, r.buffer)

			binCopy := make([]int, r.bitsCounted)
			copy(binCopy, r.buffer[:r.bitsCounted]) // Only copy the counted bits

			go r.sendCardBinary(binCopy, r.bitsCounted)

			// Clear out the buffer
			r.bitsCounted = 0
			for i := 0; i < r.BufferSize; i++ {
				r.buffer[i] = 0
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

}

// sendCardBinary sends the given integer slice as a string of integers (binary
// string) to the given notifier channel
func (r *Reader) sendCardBinary(bin []int, numBits int) {

	// Turn the int slice into a string
	buf := strings.Builder{}
	buf.Grow(numBits)

	for i := 0; i < numBits; i++ {
		buf.WriteString(strconv.Itoa(bin[i]))
	}

	// Send the string down the channel
	r.Notifier <- buf.String()

}
