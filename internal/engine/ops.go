package engine

import (
	"fmt"
	"git.agehadev.com/elliebelly/chip8le/internal"
)

type OpCode uint8

const (
	//0NNN 	Call
	//00E0 	Display
	//00EE 	Flow
	CallClearReturnOps OpCode = iota
	Goto
	Subroutine
	SkipEqualConst
	SkipNotEqualConst
	SkipEqual
	SetConst
	AddConst

	//8XY0 	Assign
	//8XY1 	BitOp
	//8XY2 	BitOp
	//8XY3 	BitOp
	//8XY4	Math
	//8XY5	Math
	//8XY6 	BitOp
	//8XY7 	Math
	//8XYE 	BitOp
	Math
	SkipNotEqual
	SetAddress
	JumpOffset
	Random
	DrawAt

	//EX9E 	KeyOp
	//EXA1 	KeyOp
	//FX07 	Timer
	//FX0A 	KeyOp
	InputOps

	//FX15 	Timer
	//FX18 	Sound
	//FX1E 	MEM
	//FX29 	MEM
	//FX33 	BCD
	//FX55 	MEM
	//FX65 	MEM
	MemoryOps
)

func (e *Engine) ExecCommand(op uint16) {
	var err error

	internal.Logger.Info(fmt.Sprintf("exec op: 0x%x", op))

	switch OpCode(op & 0xFF00 >> 12) {
	case CallClearReturnOps:
		err = execCallClearReturnOps(e, op)
	case Goto:
		err = execGoto(e, op)
	case Subroutine:
		err = execSubroutine(e, op)
	case SkipEqualConst:
		err = execSkipEqualConst(e, op)
	case SkipNotEqualConst:
		err = execSkipNotEqualConst(e, op)
	case SkipEqual:
		err = execSkipEqual(e, op)
	case SetConst:
		err = execSetConst(e, op)
	case AddConst:
		err = execAddConst(e, op)
	case Math:
		err = execMath(e, op)
	case SkipNotEqual:
		err = execSkipNotEqual(e, op)
	case SetAddress:
		err = execSetAddress(e, op)
	case JumpOffset:
		err = execJumpOffset(e, op)
	case Random:
		err = execRandom(e, op)
	case DrawAt:
		err = execDrawAt(e, op)
	case InputOps:
		err = execInputOps(e, op)
	case MemoryOps:
		err = execMemoryOps(e, op)
	default:
		panic(fmt.Sprintf("op out of range: 0x%x", op))
	}

	if err != nil {
		panic(err)
	}
}

func execCallClearReturnOps(engine *Engine, op uint16) error {
	panic("Not Implemented: CallClearReturnOps")
}

func execGoto(engine *Engine, op uint16) error {
	engine.ProgramCounter = uint(op&0x0FFF) - 2

	return nil
}

func execSubroutine(engine *Engine, op uint16) error {
	panic("Not Implemented: Subroutine")
}

func execSkipEqualConst(engine *Engine, op uint16) error {
	if engine.Registers[op&0x0F00>>8] == byte(op&0x00FF) {
		engine.ProgramCounter += 2
	}

	return nil
}

func execSkipNotEqualConst(engine *Engine, op uint16) error {
	if engine.Registers[op&0x0F00>>8] != byte(op&0x00FF) {
		engine.ProgramCounter += 2
	}

	return nil
}

func execSkipEqual(engine *Engine, op uint16) error {
	if engine.Registers[op&0x0F00>>8] == engine.Registers[op&0x00F0>>4] {
		engine.ProgramCounter += 2
	}

	return nil
}

func execSetConst(engine *Engine, op uint16) error {
	engine.Registers[op&0x0F00>>8] = byte(op & 0x00FF)

	return nil
}

func execAddConst(engine *Engine, op uint16) error {
	engine.Registers[op&0x0F00>>8] = engine.Registers[op&0x0F00>>8] + byte(op&0x00FF)

	return nil
}

func execMath(engine *Engine, op uint16) error {
	a := op & 0x0F00 >> 8
	b := op & 0x00F0 >> 4

	switch byte(op & 0x000F) {
	case 0:
		engine.Registers[a] = engine.Registers[b]
	case 1:
		engine.Registers[a] = engine.Registers[a] | engine.Registers[b]
	case 2:
		engine.Registers[a] = engine.Registers[a] & engine.Registers[b]
	case 3:
		engine.Registers[a] = engine.Registers[a] ^ engine.Registers[b]
	case 4:
		result := int(engine.Registers[a]) + int(engine.Registers[b])

		cf := 0

		if result > 0xFF {
			cf = 1
		}

		engine.Registers[0xF] = byte(cf)
		engine.Registers[a] = byte(result)
	case 5:
		result := int(engine.Registers[a]) - int(engine.Registers[b])

		nbf := 0

		if result > 0 {
			nbf = 1
		}

		engine.Registers[0xF] = byte(nbf)
		engine.Registers[a] = byte(result)
	case 6:
		engine.Registers[0xF] = engine.Registers[a] & 0x0001
		engine.Registers[a] >>= 1
	case 7:
		result := int(engine.Registers[b]) - int(engine.Registers[a])

		nbf := 0

		if result > 0 {
			nbf = 1
		}

		engine.Registers[0xF] = byte(nbf)
		engine.Registers[a] = byte(result)
	case 0xE:
		engine.Registers[0xF] = (engine.Registers[a] & 0x80) >> 4
		engine.Registers[a] <<= 1
	}

	return nil
}

func execSkipNotEqual(engine *Engine, op uint16) error {
	if engine.Registers[op&0x0F00>>8] != engine.Registers[op&0x00F0>>4] {
		engine.ProgramCounter += 2
	}

	return nil
}

func execSetAddress(engine *Engine, op uint16) error {
	engine.MemoryPointer = op & 0x0FFF

	return nil
}

func execJumpOffset(engine *Engine, op uint16) error {
	panic("Not Implemented: JumpOffset")
}

func execRandom(engine *Engine, op uint16) error {
	panic("Not Implemented: Random")
}

func execDrawAt(engine *Engine, op uint16) error {
	x := int(op & 0x0F00 >> 8)
	y := int(op & 0x00F0 >> 4)
	mp := engine.MemoryPointer

	for row := 0; row < int(op&0x000F); row++ {
		spritePixel := engine.Heap[mp]

		for i := 0; i < 8; i++ {
			index := (y+row)*int(engine.ResolutionX) + x + i
			mask := byte(0x80 >> i)

			//pixel := engine.Pixels[index]
			spritePixelOn := false

			if spritePixel&mask > 0 {
				spritePixelOn = true
			}

			engine.Pixels[index] = spritePixelOn
		}

		mp++
	}

	return nil
}

func execInputOps(engine *Engine, op uint16) error {
	panic("Not Implemented: InputOps")
}

func execMemoryOps(engine *Engine, op uint16) error {
	x := op & 0x0F00 >> 8

	switch op & 0xF0FF {
	case 0xF007:
		engine.Registers[x] = engine.DelayTimer
	case 0xF00A:
		//TODO wait for key
		panic("nope")
	case 0xF015:
		engine.DelayTimer = engine.Registers[x]
	case 0xF018:
		engine.SoundTimer = engine.Registers[x]
	case 0xF01E:
		engine.MemoryPointer += uint16(engine.Registers[x])
	case 0xF029:
		//TODO set memory pointer to text sprite location
		panic("nope")
	case 0xF033:
		num := engine.Registers[x]
		engine.Heap[engine.MemoryPointer+2] = num % 10
		engine.Heap[engine.MemoryPointer+1] = (num / 10) % 10
		engine.Heap[engine.MemoryPointer] = (num / 10 / 10) % 10
	case 0xF055:
		for i := 0; i <= int(x); i++ {
			engine.Heap[int(engine.MemoryPointer)+i] = engine.Registers[int(x)+i]
		}
	case 0xF065:
		for i := 0; i <= int(x); i++ {
			engine.Registers[int(x)+i] = engine.Heap[int(engine.MemoryPointer)+i]
		}
	}

	return nil
}
