package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/troyspencer/launch/character"
	"github.com/troyspencer/launch/timestep"
	"github.com/troyspencer/launch/view"
	"github.com/troyspencer/launch/view/camera"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Launch",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	gameview := &view.GameView{
		Window: win,
		Camera: camera.New(),
		Canvas: pixelgl.NewCanvas(pixel.R(-1000, -1000, 1000, 1000)),
	}

	playerCharacter := character.New()
	ts := timestep.New()

	for !win.Closed() {
		// calculate timestep
		ts.CalculateDelta()

		// adjust camera
		gameview.Camera.Follow(playerCharacter.Position.Center(), ts.Delta)
		gameview.Canvas.SetMatrix(gameview.Camera.Matrix)

		// get input
		input := gameview.GetInput()

		// draw character to canvas
		gameview.Canvas.Clear(colornames.Black)
		playerCharacter.Update(ts.Delta, input, gameview)
		playerCharacter.Body.Draw(gameview.Canvas)

		// stretch canvas to window
		gameview.FillWindowWithCanvas()
	}
}

func main() {
	pixelgl.Run(run)
}
