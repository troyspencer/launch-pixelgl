package timestep

import "time"

// Timestep simplifies the calculation of the Delta since the last frame of the game
type Timestep struct {
	now   time.Time
	Delta float64
}

// New creates a new Timestep with a 0 Delta and now as the current time
func New() *Timestep {
	ts := new(Timestep)
	ts.now = time.Now()
	ts.Delta = 0
	return ts
}

// CalculateDelta calculates the difference in time since the last call of CalculateDelta
func (ts *Timestep) CalculateDelta() {
	ts.Delta = time.Since(ts.now).Seconds()
	ts.now = time.Now()
}
