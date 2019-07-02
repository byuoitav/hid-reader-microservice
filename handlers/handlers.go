package handlers

import (
	"fmt"

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
}

// ReadIn .
func ReadIn() {
	pin1 := rpio.Pin(9)
	pin1.Input()
	res1 := pin1.Read()
	oldRes1 := res1

	pin2 := rpio.Pin(10)
	pin2.Input()
	res2 := pin2.Read()
	oldRes2 := res2

	for {
		if oldRes1 != res1 {
			fmt.Printf("Res1: %v", res1)
		}
		oldRes1 = res1
		res1 = pin1.Read()

		if oldRes2 != res2 {
			fmt.Printf("Res2: %v", res2)
		}
		oldRes2 = res2
		res2 = pin2.Read()
	}
}
