package core

import "GameOfLife/common"

type Character struct {
	pathFollower
}

func NewCharacter(pos common.Vec[int32]) *Character {
	return &Character{
		pathFollower: pathFollower{agentCore: newAgentCore(pos)},
	}
}

func isCharacter(a Agent) bool {
	_, ok := a.(*Character)
	return ok
}
