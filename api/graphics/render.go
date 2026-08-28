package graphics

import (
	"GameOfLife/core"
)

type Render interface {
	DrawCell(x, y int32, w *core.World) error
}
