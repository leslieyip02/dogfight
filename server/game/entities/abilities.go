package entities

import (
	"math/rand"
)

type AbilityFlag int

const (
	MultishotAbilityFlag AbilityFlag = 1 << 1
	WideBeamAbilityFlag  AbilityFlag = 1 << 2
	ShieldAbilityFlag    AbilityFlag = 1 << 3
)

func isAbilityActive(flags AbilityFlag, ability AbilityFlag) bool {
	return (flags & ability) != 0
}

func newRandomAbility() AbilityFlag {
	return AbilityFlag(1 << (1 + rand.Intn(3)))
}
