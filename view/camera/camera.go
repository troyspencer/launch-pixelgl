package camera

import (
	"math"

	"github.com/faiface/pixel"
)

// Camera is a Vec and a Matrix
type Camera struct {
	pos pixel.Vec
	pixel.Matrix
}

// New creates a new camera
func New() *Camera {
	cam := new(Camera)
	cam.pos = pixel.ZV
	cam.Matrix = pixel.IM
	return cam
}

// Follow lerps the camera position towards a position
func (cam *Camera) Follow(position pixel.Vec, dt float64) {
	cam.pos = pixel.Lerp(cam.pos, position, 1-math.Pow(1.0/128, dt))
	cam.Matrix = pixel.IM.Moved(cam.pos.Scaled(-1))
}
