package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/troyspencer/launch/character"
	"github.com/troyspencer/launch/debris"
	"github.com/troyspencer/launch/timestep"
	"github.com/troyspencer/launch/view"
	"github.com/troyspencer/launch/view/camera"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "Launch",
		Bounds:      pixel.R(0, 0, 1920, 1080),
		VSync:       true,
		Undecorated: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	gameview := &view.GameView{
		Window: win,
		Camera: camera.New(),
		Canvas: pixelgl.NewCanvas(pixel.R(0, 0, 1920, 1080)),
	}

	playerCharacter := character.New(gameview.Window.Bounds().Center())
	im := imdraw.New(nil)
	debrisChunk := &debris.RectangularDebris{
		Color:    colornames.Blue,
		Position: pixel.V(400, 800),
		Rotation: math.Pi / 3,
		Bounds:   pixel.R(0, 0, 100, 200),
		Mass:     10,
	}
	debris.DrawSingleRectangularDebris(debrisChunk, im)
	ts := timestep.New()

	for !win.Closed() {
		// calculate timestep
		ts.CalculateDelta()

		// adjust camera
		//gameview.Camera.Follow(playerCharacter.Position.Center(), ts.Delta)
		//gameview.Canvas.SetMatrix(gameview.Camera.Matrix)

		// draw character to canvas
		gameview.Canvas.Clear(colornames.Black)
		playerCharacter.Update(ts.Delta, gameview)
		playerCharacter.Body.Draw(gameview.Canvas)
		im.Draw(gameview.Canvas)

		// stretch canvas to window
		//gameview.FillWindowWithCanvas()

		gameview.DrawCanvasToWindow()
	}
}

func main() {
	pixelgl.Run(run)
}
