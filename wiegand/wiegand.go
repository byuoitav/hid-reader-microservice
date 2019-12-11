package wiegand

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/warthog618/gpio"
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

	err := gpio.Open()
	// Ignore ErrAlreadyOpen
	if err != nil && !errors.Is(err, gpio.ErrAlreadyOpen) {
		return fmt.Errorf("Error while trying to initialize GPIO pins: %w", err)
	}
	//defer gpio.Close()

	r.lastPulse = time.Now()
	r.buffer = make([]int, r.BufferSize)
	r.bitsCounted = 0

	data0 := gpio.NewPin(14)
	data1 := gpio.NewPin(15)
	data0.Input()
	data1.Input()

	err = data0.Watch(gpio.EdgeFalling, r.watchData0)
	if err != nil {
		return fmt.Errorf("Error while trying to watch data0 pin: %w", err)
	}
	data1.Watch(gpio.EdgeFalling, r.watchData1)
	if err != nil {
		return fmt.Errorf("Error while trying to watch data1 pin: %w", err)
	}

	// defer data0.Unwatch()
	// defer data1.Unwatch()

	// Catch ctrl+c plus kill commands
	catch := make(chan os.Signal, 1)
	signal.Notify(catch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-catch
		log.Println("Cleaning up...")
		data0.Unwatch()
		data1.Unwatch()
		gpio.Close()
		os.Exit(1)
	}()

	// Start watching the reader
	go r.watchForCard()

	return nil
}

func (r *Reader) watchData0(p *gpio.Pin) {
	r.lastPulse = time.Now()
	r.buffer[r.bitsCounted] = 0
	r.bitsCounted++
}

func (r *Reader) watchData1(p *gpio.Pin) {
	r.lastPulse = time.Now()
	r.buffer[r.bitsCounted] = 1
	r.bitsCounted++
}

func (r *Reader) watchForCard() {
	for {
		if time.Now().Sub(r.lastPulse) > (50*time.Millisecond) && r.bitsCounted > 0 {

			// Copy the buffer to send across the channel
			binCopy := make([]int, r.BufferSize)
			copy(binCopy, r.buffer)

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
