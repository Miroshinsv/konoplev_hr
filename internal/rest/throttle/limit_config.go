package throttle

import (
	"time"

	"golang.org/x/time/rate"
)

// LimitConfig allows to store rate.Limiter config.
type LimitConfig struct {
	Every time.Duration
	Burst int
}

func (c LimitConfig) Build() *rate.Limiter {
	return rate.NewLimiter(rate.Every(c.Every), c.Burst)
}
