package core

import "GameOfLife/common"

type Character struct {
	Pawn
}

func NewCharacter(idNation ID_NATION, pos common.Vec[int32]) Character {
	CURRENT_ID_AGENT++
	return Character{
		Pawn: Pawn{},
	}
}
