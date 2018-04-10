package input

import (
	"github.com/faiface/pixel"
)

// PlayerInput organizes the input collected from the Player
type PlayerInput struct {
	LeftClicked   bool
	RightClicked  bool
	ClickPosition pixel.Vec
}
