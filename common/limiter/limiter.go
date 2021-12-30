package limiter

import (
	"time"

	"github.com/Scrin/RuuviBridge/parser"
)

type Limiter struct {
	interval time.Duration
	macs     map[string]int64
}

func New(minInterval time.Duration) Limiter {
	return Limiter{
		interval: minInterval,
		macs:     make(map[string]int64),
	}
}

func (l Limiter) Check(m parser.Measurement) bool {
	now := time.Now().UnixNano()
	t := l.macs[m.Mac]
	if t+l.interval.Nanoseconds() > now {
		return false
	}
	l.macs[m.Mac] = now
	return true
}
