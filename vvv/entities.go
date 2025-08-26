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
	Active   bool
	Id       int
	LevelNum int
}

func (c *Checkpoint) SetActive(active bool) {
	// TODO: using magic numbers here is flaky.  Maybe
	// read in separate images/ids when creating a
	// checkpoint, and then switch between them based on the state?
	c.Active = active
	if active {
		c.srcRect = c.spriteSheet.Rect(23)
	} else {
		c.srcRect = c.spriteSheet.Rect(24)
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
