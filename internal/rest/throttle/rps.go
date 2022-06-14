package throttle

import "time"

func RPS(val float64) time.Duration {
	return time.Second / time.Duration(val)
}
