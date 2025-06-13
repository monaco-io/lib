package xrate

import (
	"golang.org/x/time/rate"
)

func NewLimit(limit int) *rate.Limiter {
	if limit <= 0 {
		return nil
	}
	return rate.NewLimiter(rate.Limit(limit), limit)
}
