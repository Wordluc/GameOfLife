package input

import "GameOfLife/common"

type Input interface {
	GetMousePointer() common.Vec[int32]
	IsKeyPressed(key any) bool
}
