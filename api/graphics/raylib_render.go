package graphics

import (
	"GameOfLife/common"
	"GameOfLife/core"
	"errors"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRender struct {
	SizeCell   common.Vec[int32]
	PeopleSeed int8
}

func (r *RaylibRender) DrawCell(x, y int32, c core.ICell) error {
	var color rl.Color
	switch c.GetType() {
	case core.GRASS:
		color = rl.Green
	case core.STONE:
		color = rl.Gray
	case core.HOUSE:
		color = rl.Brown
	default:
		return errors.New("Error type")
	}

	cellX := r.SizeCell.X * x
	cellY := r.SizeCell.Y * y
	rl.DrawRectangle(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, color)

	drawPeopleDots(cellX, cellY, r.SizeCell.X, r.SizeCell.Y, c.GetPeopleNumber(), r.PeopleSeed)
	return nil
}

func (r *RaylibRender) TickPeopleSeed() {
	r.PeopleSeed = int8(time.Now().Unix())
}

func drawPeopleDots(x, y, w, h, people int32, seed int8) {
	if people <= 0 {
		return
	}
	rand := rand.New(rand.NewSource(int64(int32(seed) * x * y)))
	for range people/10 + 1 {
		xOffset, yOffset := rand.Int31n(w-3), rand.Int31n(h-3)
		rl.DrawRectangle(xOffset+x, yOffset+y, 3, 3, rl.Black)
	}
}
