package main

type Animation struct {
	first        int
	last         int
	frame        int
	speed        int
	frameCounter int
}

func NewAnimation(first int, last int, speed int) *Animation {
	return &Animation{
		first:        first,
		last:         last,
		frame:        first,
		speed:        speed,
		frameCounter: speed,
	}
}

func (a *Animation) Update() {
	a.frameCounter--
	if a.frameCounter <= 0 {
		a.frameCounter = a.speed
		a.frame++
		if a.frame > a.last {
			a.frame = a.first
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}
