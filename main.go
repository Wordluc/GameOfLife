package main

import (
	"GameOfLife/api/graphics"
	"GameOfLife/common"
	"GameOfLife/core"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SIZE_CELL int32 = 10

func main() {
	var render graphics.Render = new(graphics.RaylibRender{
		SizeCell: common.Vec[int32]{X: SIZE_CELL, Y: SIZE_CELL},
	})

	mapSize := common.Vec[int32]{X: 50, Y: 50}
	w := core.NewWorld(mapSize)
	rl.InitWindow(mapSize.X*SIZE_CELL, mapSize.Y*SIZE_CELL, "ciao")
	w.Map.FeedMap(func(x, y int32) *core.BaseCell { return new(core.NewGrassCell()) })
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
			c.AppendPerson(w.NewPerson(core.GRASS))
			c.Touch()
		}
		if rl.IsKeyPressed(rl.KeyH) {
			p := rl.GetMousePosition()
			w.AddBlock(core.HOUSE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 5, Y: 3})
		}
		if rl.IsKeyPressed(rl.KeyR) {
			p := rl.GetMousePosition()
			w.AddBlock(core.STONE, common.Vec[int32]{X: int32(p.X) / SIZE_CELL, Y: int32(p.Y) / SIZE_CELL}, common.Vec[int32]{X: 5, Y: 5})
		}
		if rl.IsKeyPressed(rl.KeyP) {
			cells, _ := w.GetCellBlock(core.HOUSE)
			for i := range cells {
				cells[i].AppendPerson(w.NewPerson(core.GRASS))
			}
		}

		w.Map.ForEach(func(x, y int32, cell *core.BaseCell) error {
			origin, _ := w.Map.GetCell(common.Vec[int32]{X: x, Y: y})
			if origin.IsTouch() {
				return nil
			}
			if origin.GetPeopleNumber() < 2 {
				return nil
			}
			xOffset, yOffset := rand.Int31n(3), rand.Int31n(3)
			newP := common.Vec[int32]{X: x + xOffset - 1, Y: y + yOffset - 1}
			c, err := w.Map.GetCell(newP)
			if err != nil {
				return err
			}
			if c.IsTouch() {
				return nil
			}
			p := origin.PopPerson(func(p core.Person) bool { return true })
			if p != nil {
				origin.Touch()
				c.AppendPerson(p)
				c.Touch()
			}
			return nil
		})

	}
}
