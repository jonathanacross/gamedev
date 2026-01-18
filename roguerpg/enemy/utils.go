package enemy

import (
	"math"
	"math/rand"
	"roguerpg/core"
)

const (
	HeartDropOneInN = 5
)

func maybeDropHeart(loc core.Location) []core.Action {
	if rand.Intn(HeartDropOneInN) == 0 {
		return []core.Action{
			{
				Type:     core.ActionDropHeart,
				Location: loc,
			},
		}
	}
	return []core.Action{}
}

type GridAlignment int

const (
	GridAlignmentNone = iota
	GridAlignmentHorizontal
	GridAlignmentVertical
)

func getGridAlignment(a core.Character, b core.Character) GridAlignment {
	dist := core.Vector(a.Location()).Minus(core.Vector(b.Location()))
	if math.Abs(dist.Y) < float64(core.TileSize)/2 {
		return GridAlignmentHorizontal
	} else if math.Abs(dist.X) < float64(core.TileSize)/2 {
		return GridAlignmentVertical
	} else {
		return GridAlignmentNone
	}
}

func AreAligned(a core.Character, b core.Character) bool {
	alignment := getGridAlignment(a, b)
	return alignment == GridAlignmentHorizontal || alignment == GridAlignmentVertical
}

// Gets the attack direction (L/R/U/D) as a vector from a to b
func GetAttackVector(a core.Character, b core.Character) core.Vector {
	dist := core.Vector(b.Location()).Minus(core.Vector(a.Location()))
	alignment := getGridAlignment(a, b)
	switch alignment {
	case GridAlignmentHorizontal:
		if dist.X < 0 {
			return core.Vector{X: -1, Y: 0}
		} else {
			return core.Vector{X: 1, Y: 0}
		}
	case GridAlignmentVertical:
		if dist.Y < 0 {
			return core.Vector{X: 0, Y: -1}
		} else {
			return core.Vector{X: 0, Y: 1}
		}
	default:
		return core.Vector{X: 0, Y: 0}
	}
}
