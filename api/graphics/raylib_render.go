package graphics

import (
	"GameOfLife/common"
	"GameOfLife/core"
	"errors"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRender struct {
	SizeCell   common.Vec[int32]
	PeopleSeed int8
}

func (r *RaylibRender) DrawCell(x, y int32, w *core.World) error {
	var color rl.Color
	pos := common.Vec[int32]{X: x, Y: y}
	cell, err := w.CellMap.GetCell(pos)
	if err != nil {
		return err
	}
	switch cell.GetType() {
	case core.GRASS:
		color = rl.Green
	case core.STONE:
		color = rl.Gray
	case core.HOUSE:
		color = rl.Blue
	case core.WHEAT_FIELD:
		color = rl.Yellow
	case core.MINE:
		color = rl.Black
	case core.FOREST:
		color = rl.Brown
	case core.WATER:
		color = rl.DarkBlue
	case core.DOCK:
		color = rl.DarkBrown
	case core.DEBUG:
		color = rl.White
	default:
		return errors.New("Error type")
	}

	cellX := r.SizeCell.X * x
	cellY := r.SizeCell.Y * y
	var people int
	var playable bool
	for _, nation := range w.Nations {
		people += len(nation.GetAgentsAt(pos, func(a *core.Agent) bool { return a.Status != core.DEAD }))
		if len(nation.GetCharactersAt(pos)) != 0 {
			playable = true
		}
	}
	zombie := len(w.Zombies.GetAgentsAt(pos, nil))
	rl.DrawRectangle(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, color)
	rl.DrawText(fmt.Sprint(cell.VirtualNPopulation), cellX, cellY, 3, rl.Red)
	rl.DrawText(fmt.Sprint(people), cellX+10, cellY, 3, rl.Red)
	if people != 0 {
		rl.DrawText(fmt.Sprint(people), cellX+r.SizeCell.X/2, cellY+r.SizeCell.Y/2-15, 4, rl.Black)
		rl.DrawCircle(cellX+r.SizeCell.X/2+5, cellY+r.SizeCell.Y/2+5, 10, rl.Blue)
	}

	if zombie != 0 {
		rl.DrawText(fmt.Sprint(zombie), cellX+r.SizeCell.X/2, cellY+r.SizeCell.Y/2-15, 4, rl.Black)
		rl.DrawCircle(cellX+r.SizeCell.X/2+5, cellY+r.SizeCell.Y/2+5, 10, rl.Red)
	}

	if playable {
		rl.DrawCircle(cellX+r.SizeCell.X/2+5, cellY+r.SizeCell.Y/2+5, 10, rl.Orange)
	}
	return nil
}
