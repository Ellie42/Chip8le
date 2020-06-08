package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var keyMap = map[glfw.Key]uint32{
	glfw.Key1: 1 << 1,
	glfw.Key2: 1 << 2,
	glfw.Key3: 1 << 3,
	glfw.Key4: 1 << 0xC,
	glfw.KeyQ: 1 << 4,
	glfw.KeyW: 1 << 5,
	glfw.KeyE: 1 << 6,
	glfw.KeyR: 1 << 0xD,
	glfw.KeyA: 1 << 7,
	glfw.KeyS: 1 << 8,
	glfw.KeyD: 1 << 9,
	glfw.KeyF: 1 << 0xE,
	glfw.KeyZ: 1 << 0xA,
	glfw.KeyX: 1 << 0,
	glfw.KeyC: 1 << 0xB,
	glfw.KeyV: 1 << 0xF,
}

type Input struct {
	DownThisFrame uint32
	CurrentState  uint32

	Renderer          *Renderer
	StoreNextKeyPress bool
	StoredKey         int
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

	if action == glfw.Repeat {
		i.CurrentState |= flag
	}

	if action == glfw.Press {
		i.DownThisFrame |= flag
		i.CurrentState |= flag

		if i.StoreNextKeyPress {
			x := 0

			for f := flag; f > 1; f >>= 1 {
				x++
			}

			i.StoredKey = x
			i.StoreNextKeyPress = false
		}
	} else {
		i.CurrentState &= ^flag
	}
}

func (i *Input) Reset() {
	i.DownThisFrame = 0
}
