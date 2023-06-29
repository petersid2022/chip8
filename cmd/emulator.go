package chip8

import "fmt"

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type CPU struct {
	// The Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM,
	// from location 0x000 (0) to 0xFFF
	// Most Chip-8 programs start at location 0x200 (512)
	memory [4096]uint8

	// The Chip 8 has 35 opcodes which are all two bytes long.
	opcode uint16

	// V is 16 general purpose 8-bit registers, usually referred to as Vx, where x is a hexadecimal digit (0 through F)
	V [16]uint8

	// I is 16-bit registers. This register is generally used to store memory addresses
	I uint16

	//  ProgramCounter (PC) should be 16-bit, and is used to store the currently executing address.
	pc uint16

	// When these registers (delay_timer (DT) and sound_timer (ST)) are non-zero,
	// they are automatically decremented at a rate of 60Hz.
	// The system’s buzzer sounds whenever the sound timer reaches zero.
	dt uint8
	st uint8

	// Graphics:
	// The graphics of the Chip 8 are black and white and the screen has a total of 2048 pixels (64 x 32).
	display [64][32]uint8

	// The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
	sp uint8

	// The stack is an array of 16 16-bit values,
	// used to store the address that the interpreter shoud return to when finished with a subroutine.
	stack [16]uint16

	// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad
	keypad [16]uint8
}

func (cpu *CPU) Init() {
	// Program counter starts at 0x200 (512)
	pc := 0x200
	cpu.pc = uint16(pc)

	// Reset the current opcode
	cpu.opcode = 0

	// Reset the stack pointer
	cpu.sp = 0

	// Reset the index register
	cpu.I = 0

	// Clear Display
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			cpu.display[i][j] = 0
		}
	}

	// Clear stack
	for i := 0; i < len(cpu.stack); i++ {
		cpu.stack[i] = 0
	}

	// Clear register V0-VF
	for i := 0; i < 16; i++ {
		cpu.V[i] = 0
	}

	// Clear Memory
	for i := 0; i < len(cpu.memory); i++ {
		cpu.memory[i] = 0
	}

	// Load fontSet
	for i := 0; i < 80; i++ {
		cpu.memory[i] = fontSet[i]
	}

	// Reset the delay_timer and the sound_timer registers
	cpu.dt = 0
	cpu.st = 0
}

func (cpu *CPU) emulateCycle() {
	// Fetch opcode
	// 0xF0 is 1111 0000 in binary, so by doing a bitwise AND operation
	// it preserves the leftmost 4 bits and sets the rightmost 4 bits to 0.
	// 0x0F is 0000 1111 in binary, so by doing a bitwise AND operation
	// it preserves the rightmost 4 bits and sets the leftmost 4 bits to 0.
	// Finally, the two results are combined using bitwise OR to form the 16-bit opcode value.
	cpu.opcode = uint16(cpu.memory[cpu.pc]&0xF0) | uint16(cpu.memory[cpu.pc+1]&0x0F)

	// Decode opcode
	switch cpu.opcode & 0xF000 { // 0xF000 is 1111 0000 0000 0000 in binary
	case 0x0000:
		switch cpu.opcode & 0x000F {
		case 0x0000: // 0x00E0: Clears the screen
			for i := 0; i < 64; i++ {
				for j := 0; j < 32; j++ {
					cpu.display[i][j] = 0
				}
			}
		case 0x000E: // 0x00EE: Returns from subroutine
            cpu.pc = cpu.stack[cpu.sp]
            cpu.sp = cpu.sp - 1
		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", cpu.opcode)
		}
	case 0xA000: // ANNN: Sets I to the address NNN
		cpu.I = cpu.opcode & 0x0FFF
		cpu.pc += 2
    case 0x2000:
        cpu.stack[cpu.pc] = cpu.pc
        cpu.sp = cpu.sp + 1
        cpu.pc = cpu.opcode & 0x0FFF
    case 0x0004:
        if(cpu.V[(cpu.opcode & 0x00F0) >> 4] > (0xFF - cpu.V[(cpu.opcode & 0x0F00) >> 8])) {
            cpu.V[0xF] = 1 //carry
        } else {
            cpu.V[0xF] = 0
        }
        cpu.V[(cpu.opcode & 0x0F00) >> 8] += cpu.V[(cpu.opcode & 0x00F0) >> 4]
        cpu.pc = cpu.pc + 2
    case 0x0033:
        cpu.memory[cpu.I] = cpu.V[(cpu.opcode & 0x0F00) >> 8] / 100
        cpu.memory[cpu.I + 1] = (cpu.V[(cpu.opcode & 0x0F00) >> 8] / 10) % 10
        cpu.memory[cpu.I + 2] = (cpu.V[(cpu.opcode & 0x0F00) >> 8] % 100) % 10
        cpu.pc = cpu.pc + 2
	default:
		fmt.Printf("Unknown opcode: 0x%x\n", cpu.opcode)
	}

	// Update timers
	if cpu.dt > 0 {
		cpu.dt = cpu.dt - 1
	}

	if cpu.st > 0 {
		if cpu.st == 1 {
			fmt.Println("BEEP!")
			cpu.st = cpu.st - 1
		}
	}
}
