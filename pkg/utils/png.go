package utils

import (
	"bytes"
	"errors"
)

func ValidatePNG(data []byte) error {
	if len(data) < 8 {
		return errors.New("invalid png")
	}

	pngHeader := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !bytes.Equal(data[:8], pngHeader) {
		return errors.New("not png")
	}
	return nil
}
