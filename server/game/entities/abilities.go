package entities

import "math/rand"

type AbilityFlag int

const (
	MultishotAbilityFlag AbilityFlag = 1 << 1
	WideBeamAbilityFlag  AbilityFlag = 1 << 2
)

func isAbilityActive(flags AbilityFlag, ability AbilityFlag) bool {
	return (flags & ability) != 0
}

func NewRandomAbility() AbilityFlag {
	return AbilityFlag(1 << (1 + rand.Intn(2)))
}
