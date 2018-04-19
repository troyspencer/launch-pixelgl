package debris

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// TODO detection of bounds of character within debris,
// and the subsequent transfer of momentum and latching on

// RectangularDebris describes a rotatable rectangle with a center Position
type RectangularDebris struct {
	Color    color.Color
	Position pixel.Vec
	Rotation float64
	Bounds   pixel.Rect
	Mass     float64
}

func (rectDebris *RectangularDebris) getRotatedCorners() [4]pixel.Vec {
	var rotatedCorners [4]pixel.Vec
	var corners [4]pixel.Vec
	corners[0] = pixel.V(rectDebris.Bounds.Min.X, rectDebris.Bounds.Min.Y)
	corners[1] = pixel.V(rectDebris.Bounds.Min.X, rectDebris.Bounds.Max.Y)
	corners[2] = pixel.V(rectDebris.Bounds.Max.X, rectDebris.Bounds.Max.Y)
	corners[3] = pixel.V(rectDebris.Bounds.Max.X, rectDebris.Bounds.Min.Y)
	for index, corner := range corners {
		rotatedCorners[index] = rectDebris.getRotatedCorner(corner)
	}
	return rotatedCorners
}

func (rectDebris *RectangularDebris) getRotatedCorner(corner pixel.Vec) pixel.Vec {
	var rotatedPoint pixel.Vec
	rotatedPoint.X = math.Cos(rectDebris.Rotation)*(corner.X-rectDebris.Position.X) -
		math.Sin(rectDebris.Rotation)*(corner.Y-rectDebris.Position.Y) +
		rectDebris.Position.X
	rotatedPoint.Y = math.Sin(rectDebris.Rotation)*(corner.X-rectDebris.Position.X) +
		math.Cos(rectDebris.Rotation)*(corner.Y-rectDebris.Position.Y) +
		rectDebris.Position.Y
	return rotatedPoint
}

//Field is a collection of debris that are drawn together
type Field struct {
	debrisPieces [10]RectangularDebris
	imdraw.IMDraw
}

//DrawSingleRectangularDebris draws a single piece of rectangular debris in the Field
func DrawSingleRectangularDebris(rectDebris *RectangularDebris, im *imdraw.IMDraw) {
	im.Color = rectDebris.Color
	rotatedCorners := rectDebris.getRotatedCorners()
	for _, rotatedCorner := range rotatedCorners {
		im.Push(rotatedCorner)
	}
	im.Polygon(0)
}
