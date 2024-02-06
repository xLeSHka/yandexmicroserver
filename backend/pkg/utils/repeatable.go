package repeatable

import "time"

func DoWithTries(fn func() error, attemt int, delay time.Duration) (err error) {
	for attemt > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemt--
			continue
		}
		return nil
	}
	return
}
