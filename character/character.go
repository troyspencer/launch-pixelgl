package character

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/troyspencer/launch/input"
	"github.com/troyspencer/launch/view"
	"golang.org/x/image/colornames"
)

//Character is a circle controlled by the player
type Character struct {
	launchSpeed float64

	Body     *imdraw.IMDraw
	Position pixel.Rect
	vel      pixel.Vec
}

//New creates a new Character
func New() *Character {
	char := new(Character)
	char.launchSpeed = 64
	char.Position = pixel.R(-10, -10, 10, 10)
	char.createBody()
	return char
}

//drawBody Redraws the IMDraw shape with the new center position
func (char *Character) drawBody() {
	char.Body.Clear()
	char.Body.Reset()
	char.Body.Push(pixel.ZV)
	char.Body.Circle(40, 0)
	char.Body.SetMatrix(pixel.IM.Moved(char.Position.Center()))
}

func (char *Character) createBody() {
	char.Body = imdraw.New(nil)
	char.Body.Color = colornames.Gray
}

//Update takes in player input to change its velocity and update its position
func (char *Character) Update(dt float64, playerInput *input.PlayerInput, gameview *view.GameView) {
	if playerInput.LeftClicked {
		char.vel = playerInput.ClickPosition.Sub(gameview.Window.Bounds().Center()).Unit().Scaled(char.launchSpeed)
	}
	if playerInput.RightClicked {
		char.vel = pixel.ZV
	}

	// TODO: collision detection
	// check hitting canvas bounds
	if char.Position.Max.X > gameview.Canvas.Bounds().W() {
		char.vel = pixel.ZV
	}

	// apply velocity to position
	char.Position = char.Position.Moved(char.vel.Scaled(dt))

	// recalculate body
	char.drawBody()
}
