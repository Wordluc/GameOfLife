package core

import "GameOfLife/common"

type Character struct {
	pos   common.Vec[int32]
	paths *common.Queue[common.Vec[int32]]
}

func NewCharacter(pos common.Vec[int32]) Character {
	return Character{pos: pos}
}
