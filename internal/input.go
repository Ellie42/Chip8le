package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var keyMap = map[glfw.Key]uint32{
	glfw.Key1: 1 << 0,
	glfw.Key2: 1 << 1,
	glfw.Key3: 1 << 2,
	glfw.Key4: 1 << 3,
	glfw.KeyQ: 1 << 4,
	glfw.KeyW: 1 << 5,
	glfw.KeyE: 1 << 6,
	glfw.KeyR: 1 << 7,
	glfw.KeyA: 1 << 8,
	glfw.KeyS: 1 << 9,
	glfw.KeyD: 1 << 10,
	glfw.KeyF: 1 << 11,
	glfw.KeyZ: 1 << 12,
	glfw.KeyX: 1 << 13,
	glfw.KeyC: 1 << 14,
	glfw.KeyV: 1 << 15,
}

type Input struct {
	Flags     uint32
	LastFlags uint32

	Renderer *Renderer
}

func NewInput(r *Renderer) *Input {
	return &Input{
		Renderer: r,
	}
}

func (i *Input) Init() {
	i.Renderer.Window.SetKeyCallback(i.OnKeyChange)
}

func (i *Input) OnKeyChange(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	flag, ok := keyMap[key]

	if !ok {
		return
	}

	if action == glfw.Press || action == glfw.Repeat {
		i.Flags |= flag
	} else {
		i.Flags |= ^flag
	}
}
