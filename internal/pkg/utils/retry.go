package utils

import "time"

func Retry(num int, duration time.Duration, fn func() error) error {

	var err error
	for i := 0; i < num; i++ {
		if err = fn(); err == nil {
			return nil
		}

		time.Sleep(duration)
	}

	return err
}
