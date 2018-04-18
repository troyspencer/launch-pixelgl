package character

import (
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/troyspencer/launch/view"
	"golang.org/x/image/colornames"
)

type positionTimeElement struct {
	pixel.Rect
	time.Time
}

type positionBuffer struct {
	queue [10]positionTimeElement
	index int
}

func (buffer *positionBuffer) initializeBuffer(position pixel.Rect) {
	for i := 0; i < len(buffer.queue); i++ {
		buffer.addPosition(position)
	}
}

func (buffer *positionBuffer) addPosition(position pixel.Rect) {
	buffer.queue[buffer.index] = positionTimeElement{position, time.Now()}
	buffer.index++
	if buffer.index > len(buffer.queue)-1 {
		buffer.index = 0
	}
}

func (buffer *positionBuffer) getPosition(age int) positionTimeElement {
	relativeIndex := int(math.Mod(float64(buffer.index-age), float64(len(buffer.queue))))
	if relativeIndex < 0 {
		relativeIndex += len(buffer.queue)
	}
	return buffer.queue[relativeIndex]
}

//Character is a circle controlled by the player
type Character struct {
	launchSpeed float64

	Body     *imdraw.IMDraw
	Position pixel.Rect
	vel      pixel.Vec

	lastDrawTime time.Time
	drawInterval float64
	positionBuffer
}

//New creates a new Character
func New(initialPosition pixel.Vec) *Character {
	char := new(Character)
	char.launchSpeed = 400
	char.drawInterval = 0.07
	char.Position = pixel.R(-20, -20, 20, 20)
	if initialPosition != pixel.ZV {
		char.Position = char.Position.Moved(initialPosition)
	}

	char.createBody()
	return char
}

func (char *Character) createBody() {
	char.Body = imdraw.New(nil)
	char.Body.Color = colornames.Gray
	char.positionBuffer.initializeBuffer(char.Position)
}

//drawBody Redraws the IMDraw shape with the new center position
func (char *Character) drawBody() {
	char.Body.Clear()
	char.Body.Reset()
	char.drawBodyQueue()
	char.Body.SetMatrix(pixel.IM.Moved(char.Position.Center()))
}

func (char *Character) drawSingleBody() {
	char.Body.Push(pixel.ZV)
	char.Body.Circle(40, 0)
}

func (char *Character) drawBodyQueue() {
	if time.Since(char.lastDrawTime).Seconds() > char.drawInterval {
		char.lastDrawTime = time.Now()
		char.positionBuffer.addPosition(char.Position)
	}

	// draw queue oldest first
	bufferLength := len(char.positionBuffer.queue)
	for i := bufferLength - 1; i > 1; i-- {
		char.drawBodyQueueElement(i)
	}

	// always draw current position
	char.drawBodyElement(char.Position, color.RGBA{255, 255, 255, 255}, 40)
}

func (char *Character) drawBodyQueueElement(index int) {
	// calculate multiplier based on time
	bufferLength := len(char.positionBuffer.queue)
	visibilityTime := float64(bufferLength) * char.drawInterval
	timeMultiplier := time.Since(char.positionBuffer.getPosition(index).Time).Seconds() / visibilityTime
	if timeMultiplier > 1 {
		timeMultiplier = 1
	}
	colorValue := uint8(255 - 255*timeMultiplier)
	newColor := color.RGBA{colorValue, colorValue, colorValue, 255}
	char.drawBodyElement(char.positionBuffer.getPosition(index).Rect, newColor, float64(40-40*timeMultiplier))
}

func (char *Character) drawBodyElement(position pixel.Rect, bodyColor color.Color, size float64) {
	char.Body.Color = bodyColor
	char.Body.Push(position.Center().Sub(char.Position.Center()))
	char.Body.Circle(size, 0)
}

//Update takes in player input to change its velocity and update its position
func (char *Character) Update(dt float64, gameview *view.GameView) {
	playerInput := gameview.GetInput()

	if playerInput.LeftClicked {
		char.vel = playerInput.ClickPosition.Sub(char.Position.Center()).Unit().Scaled(char.launchSpeed)
	}
	if playerInput.RightClicked {
		char.vel = pixel.ZV
	}

	newPosition := char.Position.Moved(char.vel.Scaled(dt))

	// check hitting canvas bounds
	switch {
	case newPosition.Max.X > gameview.Canvas.Bounds().Max.X:
		char.vel = pixel.V(-char.vel.X, char.vel.Y)
	case newPosition.Min.X < gameview.Canvas.Bounds().Min.X:
		char.vel = pixel.V(-char.vel.X, char.vel.Y)
	case newPosition.Max.Y > gameview.Canvas.Bounds().Max.Y:
		char.vel = pixel.V(char.vel.X, -char.vel.Y)
	case newPosition.Min.Y < gameview.Canvas.Bounds().Min.Y:
		char.vel = pixel.V(char.vel.X, -char.vel.Y)
	default:
		char.Position = newPosition
	}

	// recalculate body
	char.drawBody()
}
