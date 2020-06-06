package internal

import (
	"encoding/binary"
	"time"
)

type Engine struct {
	Memory
	Timers
	Input
	Execution
	Game

	ProgramLoaded bool
	Renderer      *Renderer
	Done          chan bool
	IsStopped     bool
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

func NewEngine(renderer *Renderer) *Engine {
	return &Engine{
		Renderer: renderer,
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

	e.InputFlags = 0

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
	cycles := uint(2)

	for {
		op := binary.BigEndian.Uint16(e.Heap[e.ProgramCounter : e.ProgramCounter+2])
		opcycles := uint(2)

		if cycles < opcycles {
			break
		}

		e.ExecCommand(op)

		cycles -= opcycles

		e.ProgramCounter += 2

		if e.IsStopped {
			return
		}
	}

	e.SoundTimer--
	e.DelayTimer--
}

func (e *Engine) Run() {
	frameTimer := time.NewTicker((1000 / 60) * time.Millisecond)
	e.Done = make(chan bool)

	for {
		select {
		case <-e.Done:
			return
		case <-frameTimer.C:
			if e.Renderer.Window.ShouldClose() {
				return
			}

			if !e.IsStopped {
				e.Tick()
			}

			e.Renderer.RenderFrame(&e.Pixels)
		}
	}
}

func (e *Engine) Stop() {
	e.IsStopped = true
}
