package exportadapter

import (
	"errors"
	"fmt"
)

var ErrDenied = errors.New("export denied")

func Normalize(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("downstream: %w", err)
}
