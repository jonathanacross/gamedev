package main

type Animation struct {
	first        int
	last         int
	frame        int
	speed        int
	frameCounter int
	isFinished   bool
	looping      bool
}

func NewAnimation(first int, last int, speed int, looping bool) *Animation {
	return &Animation{
		first:        first,
		last:         last,
		frame:        first,
		speed:        speed,
		frameCounter: speed,
		isFinished:   false,
		looping:      looping,
	}
}

func (a *Animation) Update() {
	a.frameCounter--
	if a.frameCounter <= 0 {
		a.frameCounter = a.speed
		if a.frame < a.last {
			a.frame++
		} else {
			a.frame = a.first
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}

func (a *Animation) Reset() {
	a.frame = a.first
	a.frameCounter = a.speed
	a.isFinished = false
}

func (a *Animation) IsFinished() bool {
	return a.frame == a.last && a.frameCounter == a.speed
}
