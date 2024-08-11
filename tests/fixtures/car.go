package fixtures

type Car struct {
	Wheels [4]*Wheel
	Engine *Engine
}

type Engine struct {
	started bool
}

func (e *Engine) Start() {
	if e.started {
		panic("engine should be stopped")
	}
	e.started = true
}

func (e *Engine) Stop() {
	if !e.started {
		panic("engine should be started")
	}
	e.started = false
}

type WheelPosition int

const (
	WheelFrontLeft WheelPosition = iota
	WheelFrontRight
	WheelBackLeft
	WheelBackRight
)

func NewWheel(p WheelPosition) *Wheel {
	return &Wheel{position: p}
}

type Wheel struct {
	position WheelPosition
}

func (w *Wheel) Position() WheelPosition {
	return w.position
}
