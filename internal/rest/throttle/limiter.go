package throttle

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limitBuilder interface {
	Build() *rate.Limiter
}

// Limiter allows to rate limit requests by user id.
type Limiter struct {
	conf limitBuilder

	// TODO: run background gc
	users map[uint]*rate.Limiter
	m     sync.Mutex
}

func NewLimiter(conf limitBuilder) *Limiter {
	return &Limiter{
		conf:  conf,
		users: make(map[uint]*rate.Limiter),
	}
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise add a new entry to the map.
func (l *Limiter) getLimiter(userID uint) *rate.Limiter {
	l.m.Lock()
	defer l.m.Unlock()

	limiter, ok := l.users[userID]
	if ok {
		return limiter
	}

	limiter = l.conf.Build()
	l.users[userID] = limiter

	return limiter
}

func (l *Limiter) Allow(userID uint) bool {
	return l.getLimiter(userID).Allow()
}

func (l *Limiter) AllowInfo(userID uint) (allow bool, retry time.Duration) {
	lim := l.getLimiter(userID)
	now := time.Now()

	// take 1 token for reservation
	r := lim.ReserveN(now, 1)

	// if reserved for future
	delay := r.DelayFrom(now)
	if delay > 0 {
		// cancel reservation
		r.CancelAt(now)

		// return the delay
		return false, delay
	}

	// allowed to do action now
	return true, 0
}
