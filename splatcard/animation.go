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
	if a.isFinished {
		return
	}

	a.frameCounter--
	if a.frameCounter <= 0 {
		a.frameCounter = a.speed
		a.frame++
		if a.frame > a.last {
			if a.looping {
				a.frame = a.first
			} else {
				a.frame = a.last
				a.isFinished = true
			}
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
	return a.isFinished
}
