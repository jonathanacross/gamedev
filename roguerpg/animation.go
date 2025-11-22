package main

import (
	"math/rand"
)

type Animation struct {
	frames       []int
	frameIndex   int
	speed        int
	frameCounter int
	isFinished   bool
	looping      bool
}

func NewAnimation(frames []int, speed int, looping bool) *Animation {
	return &Animation{
		frames:       frames,
		frameIndex:   0,
		speed:        speed,
		frameCounter: speed,
		isFinished:   false,
		looping:      looping,
	}
}

func (a *Animation) Update() {
	if a.isFinished && !a.looping {
		return
	}

	a.frameCounter--
	if a.frameCounter <= 0 {
		a.frameCounter = a.speed
		a.frameIndex++
		if a.frameIndex >= len(a.frames) {
			if a.looping {
				a.frameIndex = 0
			} else {
				a.frameIndex = len(a.frames) - 1
				a.isFinished = true
			}
		}
	}
}

func (a *Animation) SetRandomFrame() {
	a.frameIndex = rand.Intn(len(a.frames))
	a.frameCounter = rand.Intn(a.speed)
}

func (a *Animation) Frame() int {
	return a.frames[a.frameIndex]
}

func (a *Animation) Reset() {
	a.frameIndex = 0
	a.frameCounter = a.speed
	a.isFinished = false
}

func (a *Animation) IsFinished() bool {
	return a.isFinished
}
