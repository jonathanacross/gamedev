package core

import "time"

type Timer struct {
	current  time.Duration
	duration time.Duration
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		current:  0,
		duration: duration,
	}
}

func (t *Timer) Update() {
	if t.current < t.duration {
		t.current += 16 * time.Millisecond // approximate 60 FPS
	}
}

func (t *Timer) IsReady() bool {
	return t.current >= t.duration
}

func (t *Timer) Reset() {
	t.current = 0
}
