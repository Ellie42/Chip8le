package internal

import (
	"encoding/binary"
	"time"
)

type Engine struct {
	Memory
	Timers
	Execution
	Game

	ProgramLoaded      bool
	Renderer           *Renderer
	Input              *Input
	Done               chan bool
	IsStopped          bool
	WaitForInput       bool
	InputStoreRegister uint16
}

type Game struct {
	Pixels      []bool
	ResolutionX uint8
	ResolutionY uint8
}

type Execution struct {
	ProgramCounter uint
	MemoryPointer  uint16
	StackPointer   uint8
}

type Memory struct {
	Heap      []byte
	Registers []byte
	Stack     []byte
}

type Timers struct {
	DelayTimer byte
	SoundTimer byte
}

func NewEngine(renderer *Renderer, input *Input) *Engine {
	return &Engine{
		Renderer: renderer,
		Input:    input,
		Game: Game{
			ResolutionX: 64,
			ResolutionY: 32,
		},
	}
}

func (e *Engine) Init() {
	e.SoundTimer = 0
	e.DelayTimer = 0

	e.Heap = make([]byte, 4096)
	e.Registers = make([]byte, 16)
	e.Stack = make([]byte, 256)

	copy(e.Heap, textSprites)

	e.Input.DownThisFrame = 0

	e.MemoryPointer = 0
	e.ProgramCounter = 0x200
	e.ProgramLoaded = false

	e.IsStopped = false

	e.Pixels = make([]bool, int(e.ResolutionX)*int(e.ResolutionY))
}

func (e *Engine) LoadProgram(program *Program) {
	program.Load(&e.Heap, e.ProgramCounter)
	e.ProgramLoaded = true
}

func (e *Engine) Tick() {
	cycles := uint(8)

	for {
		op := binary.BigEndian.Uint16(e.Heap[e.ProgramCounter : e.ProgramCounter+2])
		opcycles := uint(2)

		if cycles < opcycles {
			break
		}

		if e.WaitForInput {
			if e.Input.DownThisFrame == 0 {
				break
			}

			e.Registers[e.InputStoreRegister] = byte(e.Input.StoredKey)
			e.WaitForInput = false
		}

		e.ProgramCounter += 2

		e.ExecCommand(op)

		cycles -= opcycles

		if e.IsStopped {
			break
		}
	}

	e.SoundTimer--
	e.DelayTimer--
	e.Input.Reset()
}

func (e *Engine) Run() {
	lastFrame := time.Now()
	e.Done = make(chan bool)

	for {
		if e.Renderer.Window.ShouldClose() {
			return
		}

		if !e.IsStopped {
			e.Tick()
		}

		e.Renderer.RenderFrame(&e.Pixels)
		now := time.Now()

		nextFrame := lastFrame.Add(16666666)

		time.Sleep(nextFrame.Sub(now))

		lastFrame = now
	}
}

func (e *Engine) Stop() {
	e.IsStopped = true
}
