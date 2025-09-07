package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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

type HelicopterMonster struct {
	BaseSprite
	spriteSheet *GridTileSet
	animation   *Animation
	low         float64
	high        float64
	delta       float64
	horiz       bool
}

func (m *HelicopterMonster) Update() {
	m.animation.Update()

	if m.horiz {
		m.Location.X += m.delta
		m.hitbox = m.hitbox.Offset(m.delta, 0)
		if (m.Location.X < m.low && m.delta < 0) || (m.Location.X > m.high && m.delta > 0) {
			m.delta = -m.delta
		}
	} else {
		m.Location.Y += m.delta
		m.hitbox = m.hitbox.Offset(0, m.delta)
		if (m.Location.Y < m.low && m.delta < 0) || (m.Location.Y > m.high && m.delta > 0) {
			m.delta = -m.delta
		}
	}
}

func (m *HelicopterMonster) Draw(screen *ebiten.Image) {
	currSpriteFrame := m.animation.Frame()
	m.srcRect = m.spriteSheet.Rect(currSpriteFrame)
	m.BaseSprite.Draw(screen)
}

type BreakingFloorState int

const (
	Intact = iota
	Break25
	BreaK50
	BreaK75
	Broken
)

type BreakingFloor struct {
	BaseSprite
	spriteSheet *GridTileSet
	state       BreakingFloorState
	breakTimer  *Timer
}

func (b *BreakingFloor) Draw(screen *ebiten.Image) {
	currSpriteFrame := int(b.state)
	b.srcRect = b.spriteSheet.Rect(currSpriteFrame)
	b.BaseSprite.Draw(screen)
}

func (b *BreakingFloor) Update() {
	if b.state == Intact {
		return
	}

	b.breakTimer.Update()
	if b.breakTimer.IsReady() {
		switch b.state {
		case Break25:
			b.state = BreaK50
		case BreaK50:
			b.state = BreaK75
		case BreaK75:
			b.state = Broken
		case Broken:
			b.state = Intact
		default:
			b.state = Intact
		}
		b.breakTimer.Reset()
	}
}

func (b *BreakingFloor) IsSolid() bool {
	return b.state != Broken
}

func (b *BreakingFloor) StartBreak() {
	if b.state != Intact {
		return
	}
	b.state = Break25
	b.breakTimer.Reset()
}

// Called to keep the floor broken if the player is still on it.
func (b *BreakingFloor) KeepBroken() {
	if b.state != Broken {
		return
	}
	b.breakTimer.Reset()
}
