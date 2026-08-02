package graphics

import (
	"GameOfLife/common"
	"GameOfLife/core"
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRender struct {
	SizeCell common.Vec[int32]
}

func (r RaylibRender) DrawCell(x, y int32, c core.ICell) error {
	var color rl.Color
	switch c.GetType() {
	case core.GRASS:
		color = rl.Green
	case core.STONE:
		color = rl.Gray
	default:
		return errors.New("Error type")
	}
	rl.DrawRectangle(r.SizeCell.X*x, r.SizeCell.Y*y, r.SizeCell.X, r.SizeCell.Y, color)
	return nil
}
