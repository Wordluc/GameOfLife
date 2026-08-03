package graphics

import "GameOfLife/core"

type Render interface {
	DrawCell(int32, int32, core.ICell) error
	TickPeopleSeed()
}
