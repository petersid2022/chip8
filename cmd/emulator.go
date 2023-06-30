package chip8

import (
	"fmt"
	"math/rand"
)

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
	// Emulation cycle: Fetch -> Decode -> Execute
	// Every cycle, the method emulateCycle is called which emulates one cycle of the Chip 8 CPU.
	// During this cycle, the emulator will Fetch, Decode and Execute one opcode.

	// Fetch opcode
	// One way of doing that is this:
	// 0xF0 is 1111 0000 in binary, so by doing a bitwise AND operation
	// it preserves the leftmost 4 bits and sets the rightmost 4 bits to 0.
	// 0x0F is 0000 1111 in binary, so by doing a bitwise AND operation
	// it preserves the rightmost 4 bits and sets the leftmost 4 bits to 0.
	// Finally, the two results are combined using bitwise OR to form the 16-bit opcode value.
	// cpu.opcode = uint16(cpu.memory[cpu.pc]&0xF0) | uint16(cpu.memory[cpu.pc+1]&0x0F)
	// Or, you can simply shift left the cpu.memory address and then perform an OR operation with the new addr.
	cpu.opcode = (uint16(cpu.memory[cpu.pc]) << 8) | uint16(cpu.memory[cpu.pc+1])

	// Decode opcode
	// As we have stored our current opcode, we need to decode the opcode and
	// check the opcode table to see what it means.
	switch cpu.opcode & 0xF000 { // 0xF000 is 1111 0000 0000 0000 in binary
	case 0x0000:
		switch cpu.opcode & 0x000F { // 0x000F is 0000 0000 0000 1111
		case 0x0000: // 0x00E0: Clears the screen
			for i := 0; i < 64; i++ {
				for j := 0; j < 32; j++ {
					cpu.display[i][j] = 0x0
				}
			}
			cpu.pc = cpu.pc + 2
		case 0x000E: // 0x00EE: Returns from subroutine
			cpu.sp = cpu.sp - 1
			cpu.pc = cpu.stack[cpu.sp]
			cpu.pc = cpu.pc + 2
		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", cpu.opcode)
		}

	case 0x1000: // 1NNN: Jumps to address NNN
		cpu.pc = cpu.opcode & 0x0FFF

	case 0x2000: // 2NNN: Calls subroutine at NNN.
		cpu.stack[cpu.pc] = cpu.pc
		cpu.sp = cpu.sp + 1
		cpu.pc = cpu.opcode & 0x0FFF

	case 0x3000: // 3XNN: Skips the next instruction if VX equals NN
		if uint16(cpu.V[(cpu.opcode&0x0F00)>>8]) == (cpu.opcode & 0x00FF) {
			cpu.pc = cpu.pc + 4 // Skip next instruction
		} else {
			cpu.pc = cpu.pc + 2
		}

	case 0x4000: // 4XNN: Skips the next instruction if VX does not equal NN
		if uint16(cpu.V[(cpu.opcode&0x0F00)>>8]) != (cpu.opcode & 0x00FF) {
			cpu.pc = cpu.pc + 4 // Skip next instruction
		} else {
			cpu.pc = cpu.pc + 2
		}

	case 0x5000: // 5XY0: Skips the next instruction if VX equals VY
		if uint16(cpu.V[(cpu.opcode&0x0F00)>>8]) == uint16(cpu.V[(cpu.opcode&0x00F0)>>4]) {
			cpu.pc = cpu.pc + 4 // Skip next instruction
		} else {
			cpu.pc = cpu.pc + 2
		}

	case 0x6000: // 6XNN: Sets VX to NN.
		cpu.V[(cpu.opcode&0x0F00)>>8] = uint8(cpu.opcode & 0x00FF)
		cpu.pc = cpu.pc + 2

	case 0x7000: // 7XNN: Adds NN to VX (carry flag is not changed).
		cpu.V[(cpu.opcode&0x0F00)>>8] += uint8(cpu.opcode & 0x00FF)
		cpu.pc = cpu.pc + 2

	case 0x8000:
		switch cpu.opcode & 0x000F { // 0x000F is 0000 0000 0000 1111
		case 0x0000: // 8XY0: Sets Vx to the value of Vy
			cpu.V[(cpu.opcode&0x0F00)>>4] = cpu.V[(cpu.opcode & 0x00F0)]
			cpu.pc = cpu.pc + 2

		case 0x0001: // 8XY1: Sets VX to VX or VY. (bitwise OR operation)
			cpu.V[(cpu.opcode&0x0F00)>>4] = cpu.V[(cpu.opcode&0x0F00)>>4] | cpu.V[(cpu.opcode&0x00F0)]
			cpu.pc = cpu.pc + 2

		case 0x0002: // 8XY2: Sets VX to VX and VY. (bitwise AND operation)
			cpu.V[(cpu.opcode&0x0F00)>>4] = cpu.V[(cpu.opcode&0x0F00)>>4] & cpu.V[(cpu.opcode&0x00F0)]
			cpu.pc = cpu.pc + 2

		case 0x0003: // 8XY3: Sets VX to VX xor VY. (bitwise XOR operation)
			cpu.V[(cpu.opcode&0x0F00)>>4] = cpu.V[(cpu.opcode&0x0F00)>>4] ^ cpu.V[(cpu.opcode&0x00F0)]
			cpu.pc = cpu.pc + 2

		case 0x0004: // 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.
			if cpu.V[(cpu.opcode&0x00F0)>>4] > (0xFF - cpu.V[(cpu.opcode&0x0F00)>>8]) {
				cpu.V[0xF] = 1 //carry
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[(cpu.opcode&0x0F00)>>8] += cpu.V[(cpu.opcode&0x00F0)>>4]
			cpu.pc = cpu.pc + 2 // Because every instruction is 2 bytes long

		case 0x0005: // 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.

		case 0x0006: // 8XY6: Sets Vx to the value of Vy

		case 0x0007: // 8XY7: Sets Vx to the value of Vy

		case 0x000E: // 8XYE: Sets Vx to the value of Vy

		default:
			fmt.Printf("Unknown opcode [0x8000]: 0x%X\n", cpu.opcode)
		}

	case 0xA000: // ANNN: Sets I to the address NNN
		cpu.I = cpu.opcode & 0x0FFF
		cpu.pc = cpu.pc + 2 // Because every instruction is 2 bytes long

	case 0xB000: // BNNN: Jumps to the address NNN plus V0.
		cpu.pc = (cpu.opcode & 0x0FFF) + uint16(cpu.V[0x0])

	case 0xC000: // CXNN: Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		cpu.V[(cpu.opcode&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8((cpu.opcode & 0x00FF))
		cpu.pc = cpu.pc + 2

	case 0xF000:
		switch cpu.opcode & 0x00FF { // 0x000F is 0000 0000 0000 1111

		case 0x0007: // FX07: Sets VX to the value of the delay timer.
			cpu.V[(cpu.opcode&0x0F00)>>8] = cpu.dt
			cpu.pc = cpu.pc + 2

		case 0x000A: // FX0A: A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event).

		case 0x0015: // FX15: Sets the delay timer to VX.
			cpu.dt = cpu.V[(cpu.opcode&0x0F00)>>8]
			cpu.pc = cpu.pc + 2

		case 0x0018: // FX18: Sets the sound timer to VX.
			cpu.st = cpu.V[(cpu.opcode&0x0F00)>>8]
			cpu.pc = cpu.pc + 2

		case 0x001E: // FX1E: Adds VX to I. VF is not affected.
			cpu.I = cpu.I + uint16(cpu.V[(cpu.opcode&0x0F00)>>8])
			cpu.pc = cpu.pc + 2

		case 0x0029: // FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.

		case 0x0033: // FX33: Stores the binary-coded decimal representation of VX, with the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
			cpu.memory[cpu.I] = cpu.V[(cpu.opcode&0x0F00)>>8] / 100
			cpu.memory[cpu.I+1] = (cpu.V[(cpu.opcode&0x0F00)>>8] / 10) % 10
			cpu.memory[cpu.I+2] = (cpu.V[(cpu.opcode&0x0F00)>>8] % 100) % 10
			cpu.pc = cpu.pc + 2 // Because every instruction is 2 bytes long

		case 0x0055: // FX55: Stores from V0 to VX (including VX) in memory, starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.
			cpu.memory[cpu.I] = uint8(cpu.V[(cpu.opcode&0x0F00)>>8])
			for i := uint16(0); i <= ((cpu.opcode & 0x0F00) >> 8); i++ {
				cpu.memory[cpu.I+i] = cpu.V[i]
			}
			cpu.pc = cpu.pc + 2

		case 0x0065: // FX65: Fills from V0 to VX (including VX) with values from memory, starting at address I. The offset from I is increased by 1 for each value read, but I itself is left unmodified.
			cpu.V[(cpu.opcode&0x0F00)>>8] = cpu.memory[cpu.I]
			for i := uint16(0); i <= ((cpu.opcode & 0x0F00) >> 8); i++ {
				cpu.V[i] = cpu.memory[cpu.I+i]
			}
			cpu.pc = cpu.pc + 2

		default:
			fmt.Printf("Unknown opcode [0x8000]: 0x%X\n", cpu.opcode)
		}

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
