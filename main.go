package main

import (
	"GameOfLife/api/graphics"
	"GameOfLife/common"
	"GameOfLife/core"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SIZE_CELL int32 = 20

func main() {
	var render graphics.Render = new(graphics.RaylibRender{
		SizeCell: common.Vec[int32]{X: SIZE_CELL, Y: SIZE_CELL},
	})

	mapSize := common.Vec[int32]{X: 50, Y: 30}
	w := core.NewWorld(mapSize)
	rl.InitWindow(mapSize.X*SIZE_CELL, mapSize.Y*SIZE_CELL, "ciao")
	w.GenerateMap()
	var timer float32 = 0.0
	currentJob := 0
	for !rl.WindowShouldClose() {
		core.TOUCH_ID = rand.Int()
		rl.ClearBackground(rl.White)
		rl.BeginDrawing()
		render.TickPeopleAnimation()
		w.Map.ForEach(func(x, y int32, cell *core.BaseCell) error {
			return render.DrawCell(x, y, cell)
		})
		rl.EndDrawing()
		if rl.IsKeyPressed(rl.KeyF) {
			p := rl.GetMousePosition()
			c, _ := w.Map.GetCell(common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL})
			w.NewPerson(core.GRASS, c)
			c.Touch()
		}
		if rl.IsKeyPressed(rl.KeyQ) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.WHEAT, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 3, Y: 3})
			if err != nil {
				panic(err)
			}
		}
		if rl.IsKeyPressed(rl.KeyM) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.MINE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				panic(err)
			}
		}
		if rl.IsKeyPressed(rl.KeyH) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.HOUSE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				panic(err)
			}
		}
		if rl.IsKeyPressed(rl.KeyR) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.STONE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 5, Y: 5})
			if err != nil {
				panic(err)
			}
		}
		if rl.IsKeyPressed(rl.KeyZero) {
			currentJob = int(core.FARMER)
			println("Farmer")
		}
		if rl.IsKeyPressed(rl.KeyOne) {
			currentJob = int(core.MINER)
			println("Miner")
		}
		if rl.IsKeyPressed(rl.KeyP) {
			cells, _ := w.GetCellBlock(core.HOUSE)
			for i := range cells {
				w.NewPerson(core.Job(currentJob), cells[i])
			}
		}
		timer += rl.GetFrameTime()
		if timer >= 0.5 {
			err := w.MovementSimulation()
			if err != nil {
				panic(err)
			}
			timer = 0 // or timer = 0, but -= preserves overshoot accuracy
		}
	}
}
