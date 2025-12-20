package enemy

import (
	"math/rand"
	"roguerpg/core"
)

const (
	HeartDropOneInN = 1
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
