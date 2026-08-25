package graphics

import (
	"GameOfLife/common"
	"GameOfLife/core"
	"errors"
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRender struct {
	SizeCell   common.Vec[int32]
	PeopleSeed int8
}

func (r *RaylibRender) DrawCell(x, y int32, w *core.World) error {
	var color rl.Color
	pos := common.Vec[int32]{X: x, Y: y}
	cell, err := w.Map.GetCell(pos)
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
	for _, nation := range w.Nation {
		people += len(nation.GetAgentsIdInsideCellf(pos, func(a *core.Agent) bool { return a.Status != core.DEAD }))
	}
	zombie := len(w.Zombies.GetAgentsIdInsideCell(pos))
	rl.DrawRectangle(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, color)
	rl.DrawText(fmt.Sprint(cell.VirtualNPopulation), cellX, cellY, 3, rl.Red)
	rl.DrawText(fmt.Sprint(people), cellX+10, cellY, 3, rl.Red)

	drawPeopleDots(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, people, r.PeopleSeed, rl.Red)
	drawPeopleDots(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, zombie, r.PeopleSeed, rl.Black)
	return nil
}

func (r *RaylibRender) TickPeopleAnimation() {
	r.PeopleSeed = int8(time.Now().Unix())
}

func drawPeopleDots(x, y, w, h int32, people int, seed int8, color rl.Color) {
	if people <= 0 {
		return
	}
	rand := rand.New(rand.NewSource(int64(int32(seed) * x * y)))
	for range people {
		xOffset, yOffset := rand.Int31n(w-3), rand.Int31n(h-3)
		rl.DrawRectangle(xOffset+x, yOffset+y, 3, 3, color)
	}
}
