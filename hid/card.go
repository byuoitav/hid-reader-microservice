package hid

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrBadFormat = errors.New("The given binary was malformatted")

// GetCardID returns the card ID from the HID formatted binary string passed in
func GetCardID(binary string) (string, error) {

	// We expect 48 bits here
	if len(binary) < 48 {
		return "", ErrBadFormat
	}

	// Convert the binary string to an int
	id, err := strconv.ParseUint(binary[24:47], 2, 32)
	if err != nil {
		return "", fmt.Errorf("Error while parsing binary: %w", err)
	}

	// Pad to length 6 with leading 0's
	return fmt.Sprintf("%06d", id), nil

}
