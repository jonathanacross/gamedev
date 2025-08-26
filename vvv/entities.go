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

func (le *LevelExit) HitRect() Rect {
	return le.Rect
}
