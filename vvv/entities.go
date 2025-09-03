package main

import "github.com/hajimehoshi/ebiten/v2"

type Tile struct {
	BaseSprite
	solid bool
}

type Spike struct {
	BaseSprite
}

func (s *Spike) Update() {}

type Checkpoint struct {
	BaseSprite
	spriteSheet *GridTileSet
	Active      bool
	Id          int
	LevelNum    int
}

func (c *Checkpoint) SetActive(active bool) {
	c.Active = active
	if active {
		c.srcRect = c.spriteSheet.Rect(0)
	} else {
		c.srcRect = c.spriteSheet.Rect(1)
	}
}

func (c *Checkpoint) Update() {}

type Crystal struct {
	BaseSprite
	Collected bool
}

func (c *Crystal) Update() {}

func (c *Crystal) Draw(screen *ebiten.Image) {
	if !c.Collected {
		c.BaseSprite.Draw(screen)
	}
}

type Platform struct {
	BaseSprite
	low   float64
	high  float64
	delta float64
	horiz bool
}

func (p *Platform) Update() {
	if p.horiz {
		p.Location.X += p.delta
		p.hitbox = p.hitbox.Offset(p.delta, 0)
		if (p.Location.X < p.low && p.delta < 0) || (p.Location.X > p.high && p.delta > 0) {
			p.delta = -p.delta
		}
	} else {
		p.Location.Y += p.delta
		p.hitbox = p.hitbox.Offset(0, p.delta)
		if (p.Location.Y < p.low && p.delta < 0) || (p.Location.Y > p.high && p.delta > 0) {
			p.delta = -p.delta
		}
	}
}

type LevelExit struct {
	Rect
	ToLevel int
}

func (le LevelExit) HitBox() Rect {
	return le.Rect
}

func (le LevelExit) Update() {}
