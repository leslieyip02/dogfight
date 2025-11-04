package entities

import (
	"math/rand"
)

// An AbilityFlag is a number where each bit corresponds to a particular
// powerup ability. Multiple abilities can be active at the same time.
type AbilityFlag uint32

const (
	MultishotAbilityFlag AbilityFlag = 1 << 1 // shoots 3 bullets in parallel
	WideBeamAbilityFlag  AbilityFlag = 1 << 2 // shoots 1 bullet
	ShieldAbilityFlag    AbilityFlag = 1 << 3 // protects against 1 collision
)

func isAbilityActive(flags AbilityFlag, ability AbilityFlag) bool {
	return (flags & ability) != 0
}

func newRandomAbility() AbilityFlag {
	return AbilityFlag(1 << (1 + rand.Intn(3)))
}
