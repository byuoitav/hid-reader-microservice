package handlers

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/stianeikeland/go-rpio"
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
