package handlers

import (
	"fmt"
	"time"

	"github.com/go-rpio"
	"github.com/labstack/echo"
)

// Beep activates the beep
func Beep(ctx echo.Context) error {
	fmt.Printf("BEEP\n")
	pin := rpio.Pin(5)
	pin.Output()
	pin.Low()

	return nil
}

// BeepOff turns off the beep
func BeepOff(ctx echo.Context) error {
	fmt.Printf("BEEP OFF\n")
	pin := rpio.Pin(5)
	pin.Output()
	pin.High()

	return nil
}

// Green activates a green light
func Green(ctx echo.Context) error {
	fmt.Printf("GREEN\n")
	pin := rpio.Pin(11)
	pin.Output()
	pin.Low()

	return nil
}

// GreenOff deactivates a green light
func GreenOff(ctx echo.Context) error {
	fmt.Printf("GREEN OFF\n")
	pin := rpio.Pin(11)
	pin.Output()
	pin.High()

	return nil
}

// Red activates a red light
func Red(ctx echo.Context) error {
	fmt.Printf("RED\n")
	pin := rpio.Pin(6)
	pin.Output()
	pin.Low()

	return nil
}

// RedOff deactivates a red light
func RedOff(ctx echo.Context) error {
	fmt.Printf("RED OFF\n")
	pin := rpio.Pin(6)
	pin.Output()
	pin.High()

	return nil
}

// Hold .
func Hold(ctx echo.Context) error {
	fmt.Printf("HOLD\n")
	pin := rpio.Pin(13)
	pin.Output()
	pin.Low()

	return nil
}

// HoldOff .
func HoldOff(ctx echo.Context) error {
	fmt.Printf("HOLD OFF\n")
	pin := rpio.Pin(13)
	pin.Output()
	pin.High()

	return nil
}

// StartListen .
func StartListen() {
	go ReadIn()
	//	go tempRead()
}

func tempRead() {
	pin := rpio.Pin(6)
	pin.Input()
	for {
		var bytes []int
		for {
			if len(bytes) >= 10 {
				break
			}
			res := pin.Read()
			if res == rpio.Low {
				bytes = append(bytes, 1)
			}
		}
		fmt.Printf("Bytes: %v\n", bytes)
	}
}

// ReadIn .
func ReadIn() {
	pin0 := rpio.Pin(6)
	pin0.Input()
	pin1 := rpio.Pin(10)
	pin1.Input()
	for {
		r0 := pin0.Read()
		r1 := pin1.Read()
		if r0 == rpio.Low || r1 == rpio.Low {
			if r1 == rpio.Low {
				var bytes []int
				if r0 == rpio.Low {
					bytes = append(bytes, 0)
				} else {
					bytes = append(bytes, 1)
				}

				read0 := false
				read1 := false
				for {
					if len(bytes) >= 130 {
						break
					}
					r0 = pin0.Read()
					r1 = pin1.Read()
					if r0 == rpio.High {
						read0 = false
					}
					if r1 == rpio.High {
						read1 = false
					}
					if r0 == rpio.Low && !read0 {
						bytes = append(bytes, 0)
						read0 = true
					}
					if r1 == rpio.Low && !read1 {
						bytes = append(bytes, 1)
						read1 = true
					}
				}
				fmt.Printf("Bytes: %v\n", bytes)
				time.Sleep(3 * time.Second)
			}
		}
	}
}
