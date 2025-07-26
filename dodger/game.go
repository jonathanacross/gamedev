package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	score            int
}

var gameInstance = NewGame()

func (g *Game) CheckCollisions() {
	// meteor / bullet collisions
	for i, m := range g.meteors {
		for j, b := range g.bullets {
			if m.HitRect().Intersects(b.HitRect()) {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
				g.score++
			}
		}
	}

	// meteor / player collision
	for _, m := range g.meteors {
		if m.HitRect().Intersects(g.player.HitRect()) {
			g.Reset()
		}
	}
}

func (g *Game) Update() error {
	g.player.Update()

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		m := NewMeteor()
		g.meteors = append(g.meteors, m)
	}

	g.meteors = slices.DeleteFunc(g.meteors, func(m *Meteor) bool {
		return m.HasFallenOffscreen()
	})

	g.bullets = slices.DeleteFunc(g.bullets, func(b *Bullet) bool {
		return b.HasFallenOffscreen()
	})

	for _, m := range g.meteors {
		m.Update()
	}

	for _, b := range g.bullets {
		b.Update()
	}

	g.CheckCollisions()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	g.player.Draw(screen)

	text.Draw(screen, fmt.Sprintf("%06d", g.score), ScoreFont, 40, 40, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewGame() *Game {
	return &Game{
		player:           NewPlayer(),
		meteorSpawnTimer: NewTimer(500 * time.Millisecond),
		meteors:          []*Meteor{},
		bullets:          []*Bullet{},
		score:            0,
	}
}

func (g *Game) Reset() {
	g.player = NewPlayer()
	g.meteors = []*Meteor{}
	g.bullets = []*Bullet{}
	g.score = 0

}
