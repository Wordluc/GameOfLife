package core

import "GameOfLife/common"

type Character struct {
	Agent
}

func NewCharacter(idNation ID_NATION, pos common.Vec[int32]) Character {
	return Character{
		newAgent("", idNation, pos),
	}
}
