package main

import (
	"GameOfLife/api/graphics"
	"GameOfLife/common"
	"GameOfLife/core"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SIZE_CELL int32 = 30

func main() {
	var render graphics.Render = new(graphics.RaylibRender{
		SizeCell: common.Vec[int32]{X: SIZE_CELL, Y: SIZE_CELL},
	})

	mapSize := common.Vec[int32]{X: 40, Y: 25}
	w := core.NewWorld(mapSize)
	rl.InitWindow(mapSize.X*SIZE_CELL, mapSize.Y*SIZE_CELL, "ciao")
	w.GenerateMap()
	var timer float32 = 0.0
	var currentJob core.Job = core.FARMER
	for !rl.WindowShouldClose() {
		core.TOUCH_ID = rand.Int()
		rl.ClearBackground(rl.White)
		rl.BeginDrawing()
		render.TickPeopleAnimation()
		w.Map.ForEach(func(x, y int32, cell *core.BaseCell) error {
			return render.DrawCell(x, y, cell)
		})
		rl.EndDrawing()

		if rl.IsKeyPressed(rl.KeyG) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.WHEAT_FIELD, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyM) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.MINE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyH) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.HOUSE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyF) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.FOREST, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyR) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.STONE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 5, Y: 5})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyW) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.WATER, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 3, Y: 3})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyD) {
			p := rl.GetMousePosition()
			err := w.AddBlock(core.DOCK, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyZero) {
			currentJob = core.FARMER
			println("Farmer")
		}
		if rl.IsKeyPressed(rl.KeyOne) {
			currentJob = core.MINER
			println("Miner")
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			currentJob = core.WOODCUTTER
			println("Forest")
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			currentJob = core.FISHERMAN
			println("Fischerman")
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
			timer -= 0.5 // or timer = 0, but -= preserves overshoot accuracy
		}
	}
}
