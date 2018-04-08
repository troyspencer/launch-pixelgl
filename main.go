package main

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type playerInput struct {
	leftClicked   bool
	rightClicked  bool
	clickPosition pixel.Vec
}

func (view *gameView) getInput() *playerInput {
	input := new(playerInput)
	if view.Window.JustPressed(pixelgl.MouseButtonLeft) {
		input.clickPosition = view.camera.matrix.Unproject(view.Window.MousePosition())
		input.leftClicked = true
	}
	if view.Window.JustPressed(pixelgl.MouseButtonRight) {
		input.clickPosition = view.camera.matrix.Unproject(view.Window.MousePosition())
		input.rightClicked = true
	}
	return input
}

type character struct {
	launchSpeed float64

	body     *imdraw.IMDraw
	position pixel.Rect
	vel      pixel.Vec
}

func newCharacter() *character {
	char := new(character)
	char.launchSpeed = 64
	char.position = pixel.R(-10, -10, 10, 10)
	char.createBody()
	return char
}

func (char *character) drawBody() {
	char.body.Clear()
	char.body.Reset()
	char.body.Push(pixel.ZV)
	char.body.Circle(40, 0)
	char.body.SetMatrix(pixel.IM.Moved(char.position.Center()))
}

func (char *character) createBody() {
	char.body = imdraw.New(nil)
	char.body.Color = colornames.Gray
}

func (char *character) update(dt float64, input *playerInput, view *gameView) {
	if input.leftClicked {
		char.vel = input.clickPosition.Sub(view.Window.Bounds().Center()).Unit().Scaled(char.launchSpeed)
	}
	if input.rightClicked {
		char.vel = pixel.ZV
	}

	// TODO: collision detection
	// check hitting canvas bounds
	if char.position.Max.X > view.Canvas.Bounds().W() {
		char.vel = pixel.ZV
	}

	// apply velocity to position
	char.position = char.position.Moved(char.vel.Scaled(dt))

	// recalculate body
	char.drawBody()
}

type gameView struct {
	*pixelgl.Window
	*pixelgl.Canvas
	*camera
}

func (view *gameView) fillWindowWithCanvas() {
	view.Window.Clear(colornames.Black)
	view.Window.SetMatrix(pixel.IM.Scaled(pixel.ZV,
		math.Min(
			view.Window.Bounds().W()/view.Canvas.Bounds().W(),
			view.Window.Bounds().H()/view.Canvas.Bounds().H(),
		),
	).Moved(view.Window.Bounds().Center()))
	view.Canvas.Draw(view.Window, pixel.IM.Moved(view.Canvas.Bounds().Center()))
	view.Window.Update()
}

type camera struct {
	pos    pixel.Vec
	matrix pixel.Matrix
}

func newCamera() *camera {
	cam := new(camera)
	cam.pos = pixel.ZV
	cam.matrix = pixel.IM
	return cam
}

func (cam *camera) followCharacter(characterPos pixel.Vec, dt float64) {
	cam.pos = pixel.Lerp(cam.pos, characterPos, 1-math.Pow(1.0/128, dt))
	cam.matrix = pixel.IM.Moved(cam.pos.Scaled(-1))
}

type timestep struct {
	now   time.Time
	delta float64
}

func newTimestep() *timestep {
	ts := new(timestep)
	ts.now = time.Now()
	ts.delta = 0
	return ts
}

func (ts *timestep) calc() {
	ts.delta = time.Since(ts.now).Seconds()
	ts.now = time.Now()
}

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

	view := &gameView{
		Window: win,
		camera: newCamera(),
		Canvas: pixelgl.NewCanvas(pixel.R(-1000, -1000, 1000, 1000)),
	}

	playerCharacter := newCharacter()
	ts := newTimestep()

	for !win.Closed() {
		// calculate timestep
		ts.calc()

		// adjust camera
		view.camera.followCharacter(playerCharacter.position.Center(), ts.delta)
		view.Canvas.SetMatrix(view.camera.matrix)

		// get input
		input := view.getInput()

		// draw character to canvas
		view.Canvas.Clear(colornames.Black)
		playerCharacter.update(ts.delta, input, view)
		playerCharacter.body.Draw(view.Canvas)

		// stretch canvas to window
		view.fillWindowWithCanvas()
	}
}

func main() {
	pixelgl.Run(run)
}
