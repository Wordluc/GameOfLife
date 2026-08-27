package main

import (
	"GameOfLife/api/graphics"
	"GameOfLife/common"
	"GameOfLife/core"
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SIZE_CELL int32 = 36 * 2

func approssimare(a float32, t int) float32 {
	return float32(math.Floor(float64(a*float32(math.Pow10(t)))) / math.Pow10(t))
}

var fullscreen bool = false

func main() {
	var render graphics.Render = new(graphics.RaylibRender{
		SizeCell: common.Vec[int32]{X: SIZE_CELL, Y: SIZE_CELL},
	})

	totalMap := common.Vec[int32]{X: 16 * 5, Y: 9 * 5}
	visibleMap := common.Vec[int32]{X: 16 * 1, Y: 9 * 1}
	visibleWindow := visibleMap.Clone().MultScalar(SIZE_CELL)
	rlVisibleWindow := rl.Rectangle{X: 0, Y: 0, Width: float32(visibleWindow.X), Height: -float32(visibleWindow.Y)}
	w := core.NewWorld(totalMap)

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(visibleWindow.X, visibleWindow.Y, "ciao")
	dpiScale := rl.GetWindowScaleDPI()
	if fullscreen {
		rl.ToggleFullscreen()
	}
	w.GenerateMap()
	var timer float32 = 0.0
	var currentJob core.Job = core.FARMER
	var err error
	rl.SetTargetFPS(60)
	camera := rl.Camera2D{}
	camera.Zoom = 1
	var texture rl.RenderTexture2D = rl.LoadRenderTexture(visibleWindow.X, visibleWindow.Y)

	virtualW := float32(visibleWindow.X)
	virtualH := float32(visibleWindow.Y)

	//FINAL WIDTH RESOLUTION
	screenW := float32(rl.GetScreenWidth())
	//FINAL HEIGHT RESOLUTION
	screenH := float32(rl.GetScreenHeight())

	scale := float32(math.Min(
		float64(screenW/virtualW),
		float64(screenH/virtualH),
	))

	destW := virtualW * scale
	destH := virtualH * scale

	destX := (screenW - destW) / 2
	destY := (screenH - destH) / 2

	rlFinalWindow := rl.Rectangle{
		X:      destX,
		Y:      destY,
		Width:  destW,
		Height: destH,
	}
	getPosMouse := func() common.Vec[int32] {
		p := rl.GetMousePosition()
		if fullscreen {
			p.X *= dpiScale.X
			p.Y *= dpiScale.Y
		}

		texX := (p.X - destX) * (virtualW / destW)
		texY := (p.Y - destY) * (virtualH / destH)

		if texX < 0 || texX >= virtualW || texY < 0 || texY >= virtualH {
			return common.Vec[int32]{X: -1, Y: -1}
		}

		worldPos := rl.GetScreenToWorld2D(rl.Vector2{X: texX, Y: texY}, camera)

		return common.Vec[int32]{
			X: int32(worldPos.X) / SIZE_CELL,
			Y: int32(worldPos.Y) / SIZE_CELL,
		}
	}
	for !rl.WindowShouldClose() {
		core.TOUCH_MOVE_PERSON_ID = rand.Int()

		{
			rl.BeginTextureMode(texture)
			{
				rl.BeginMode2D(camera)
				render.TickPeopleAnimation()
				err = w.Map.ForEach(func(x, y int32) error {
					return render.DrawCell(x, y, w)
				})
				if err != nil {
					panic(err)
				}
				//				w.Zombies.BfsMap.ForEach(func(x, y int32) error {
				//					cellX := SIZE_CELL * x
				//					cellY := SIZE_CELL * y
				//					n, _ := w.Zombies.BfsMap.GetCell(common.Vec[int32]{X: x, Y: y})
				//					if n != nil {
				//						rl.DrawText(fmt.Sprint(*n), cellX+SIZE_CELL/2, cellY+SIZE_CELL/2, 3, rl.Red)
				//					}
				//					return nil
				//				})
				rl.EndMode2D()
			}
			rl.EndTextureMode()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawTexturePro(texture.Texture, rlVisibleWindow, rlFinalWindow, rl.Vector2{}, 0, rl.White)
		rl.EndDrawing()

		if rl.IsKeyPressed(rl.KeyG) {
			err = w.AddBlock(core.WHEAT_FIELD, getPosMouse(), common.Vec[int32]{X: 3, Y: 3})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyM) {
			err = w.AddBlock(core.MINE, getPosMouse(), common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyH) {
			err = w.AddBlock(core.HOUSE, getPosMouse(), common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyF) {
			err = w.AddBlock(core.FOREST, getPosMouse(), common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyR) {
			err = w.AddBlock(core.STONE, getPosMouse(), common.Vec[int32]{X: 1, Y: 1})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyW) {
			err = w.AddBlock(core.WATER, getPosMouse(), common.Vec[int32]{X: 3, Y: 3})
			if err != nil {
				println(err.Error())
			}
		}
		if rl.IsKeyPressed(rl.KeyD) {
			err = w.AddBlock(core.DOCK, getPosMouse(), common.Vec[int32]{X: 1, Y: 1})
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
			cells, err := w.GetCellsByType(core.HOUSE)
			if err != nil {
				println(err.Error())
			} else {
				for i := range cells {
					w.NewPerson(core.Job(currentJob), cells[i].GetPos(), 0)
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyZ) {
			w.Zombies.NewZombie(getPosMouse())
		}
		timer += rl.GetFrameTime()
		w.RefreshZombieVision()
		w.PerformPathFinding()
		if timer >= 0.5 {
			err = w.ZombieEat()
			if err != nil {
				panic(err)
			}
			w.HarvestingSimulation()
			w.StarvingSimulation()
			err := w.MovementSimulation()
			if err != nil {
				panic(err)
			}
			timer -= 0.5
		}
		if 1/rl.GetFrameTime() < 50 {
			fmt.Printf("fps:%v\n", 1/rl.GetFrameTime())
		}

	}
	rl.UnloadRenderTexture(texture)
}
