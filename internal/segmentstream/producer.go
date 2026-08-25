package segmentstream

import "errors"

var ErrDecode = errors.New("segment decode failed")

func Start(parts []string) (<-chan string, <-chan error) {
	data := make(chan string)
	errs := make(chan error, 1)
	go func() {
		for _, p := range parts {
			if p == "bad" {
				errs <- ErrDecode
				return
			}
			data <- p
		}
		close(data)
		close(errs)
	}()
	return data, errs
}
