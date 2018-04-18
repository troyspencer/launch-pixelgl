package view

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/troyspencer/launch/input"
	"github.com/troyspencer/launch/view/camera"
	"golang.org/x/image/colornames"
)

//GameView organizes the canvas, window, and camera to operate them smoothly together
type GameView struct {
	*pixelgl.Window
	*pixelgl.Canvas
	*camera.Camera
}

// FillWindowWithCanvas stretches the Canvas to fit the bounds of the window
func (gameview *GameView) FillWindowWithCanvas() {
	gameview.Window.Clear(colornames.Black)
	gameview.Window.SetMatrix(pixel.IM.Scaled(pixel.ZV,
		math.Min(
			gameview.Window.Bounds().W()/gameview.Canvas.Bounds().W(),
			gameview.Window.Bounds().H()/gameview.Canvas.Bounds().H(),
		),
	).Moved(gameview.Window.Bounds().Center()))
	gameview.Canvas.Draw(gameview.Window, pixel.IM.Moved(gameview.Canvas.Bounds().Center()))
	gameview.Window.Update()
}

func (gameview *GameView) DrawCanvasToWindow() {
	gameview.Window.Clear(colornames.Black)
	gameview.Canvas.Draw(gameview.Window, pixel.IM.Moved(gameview.Canvas.Bounds().Center()))
	gameview.Window.Update()
}

func (gameview *GameView) GetInput() *input.PlayerInput {
	playerInput := new(input.PlayerInput)

	// check for escape
	gameview.Window.SetClosed(gameview.Window.JustPressed(pixelgl.KeyEscape) || gameview.Window.JustPressed(pixelgl.KeyQ))

	if gameview.Window.Pressed(pixelgl.MouseButtonLeft) {
		playerInput.ClickPosition = gameview.Camera.Matrix.Unproject(gameview.Window.MousePosition())
		playerInput.LeftClicked = true
	}
	if gameview.Window.JustPressed(pixelgl.MouseButtonRight) {
		playerInput.ClickPosition = gameview.Camera.Matrix.Unproject(gameview.Window.MousePosition())
		playerInput.RightClicked = true
	}
	return playerInput
}
