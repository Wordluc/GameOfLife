package main

import (
	"GameOfLife/api/graphics"
	"GameOfLife/common"
	"GameOfLife/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SIZE_CELL int32 = 5

func main() {
	var render graphics.Render = graphics.RaylibRender{
		SizeCell: common.Vec[int32]{X: SIZE_CELL, Y: SIZE_CELL},
	}
	mapSize := common.Vec[int32]{X: 100, Y: 100}
	rl.InitWindow(mapSize.X*SIZE_CELL, mapSize.Y*SIZE_CELL, "ciao")
	m := core.NewMap(mapSize)
	m.FeedMap(func(x, y int32) core.ICell {
		return new(core.NewGrassCell())
	})
	for !rl.WindowShouldClose() {
		rl.ClearBackground(rl.White)
		rl.BeginDrawing()
		m.ForEach(func(x, y int32, cell core.ICell) error {
			return render.DrawCell(x, y, cell)
		})
		rl.EndDrawing()
		if rl.IsKeyPressed(rl.KeyF) {
			p := rl.GetMousePosition()
			m.SetRawCell(new(core.NewStoneCell()), common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL})
		}
	}
}
