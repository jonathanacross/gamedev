package main

type Tile struct {
	BaseSprite
	solid bool
}

type Spike struct {
	BaseSprite
}

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

type Platform struct {
	BaseSprite
	Vx             float64
	Vy             float64
	startX, startY float64
	endX, endY     float64
}

type LevelExit struct {
	Rect
	ToLevel int
}

func (le LevelExit) HitBox() Rect {
	return le.Rect
}
