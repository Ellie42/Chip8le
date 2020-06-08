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
	Pixels      []uint
	ResolutionX uint8
	ResolutionY uint8
}

type Execution struct {
	ProgramCounter uint
	MemoryPointer  uint16
}

type UintStack struct{
	stack []uint
	index int
}

func (s *UintStack) Push(i uint){
	s.stack[s.index] = i
	s.index++
}

func (s *UintStack) Pop() uint{
	s.index--
	i := s.stack[s.index]
	s.stack[s.index] = 0
	return i
}

type Memory struct {
	Heap      []byte
	Registers []byte
	Stack     *UintStack
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
	e.Stack = &UintStack{
		stack: make([]uint, 256),
		index: 0,
	}

	copy(e.Heap, textSprites)

	e.Input.DownThisFrame = 0

	e.MemoryPointer = 0
	e.ProgramCounter = 0x200
	e.ProgramLoaded = false

	e.IsStopped = false

	e.Pixels = make([]uint, int(e.ResolutionX)*int(e.ResolutionY))
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
			break
		}

		e.ExecCommand(op)

		e.ProgramCounter += 2

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
		//hexNum := make([]byte,4)
		//binary.BigEndian.PutUint32(hexNum, e.Input.CurrentState)
		//fmt.Printf("%s\n", hex.EncodeToString(hexNum))

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
