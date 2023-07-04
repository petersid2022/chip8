package chip8

import (
	"fmt"
	"math/rand"
	"os"
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
	// MoSound_timer Chip-8 programs start at location 0x200 (512)
	Memory [4096]uint8

	// The Chip 8 has 35 Opcodes which are all two bytes long.
	Opcode uint16

	// V is 16 general purpose 8-bit registers, usually referred to as Vx, where x is a hexadecimal digit (0 through F)
	V [16]uint8

	// I is 16-bit registers. This register is generally used to store Memory addresses
	I uint16

	//  ProgramCounter (PC) should be 16-bit, and is used to store the currently executing address.
	Pc uint16

	// When these registers (delay_timer (DT) and sound_timer (ST)) are non-zero,
	// they are automatically decremented at a rate of 60Hz.
	// The system’s buzzer sounds whenever the sound timer reaches zero.
	Delay_timer uint8
	Sound_timer uint8

	// Graphics:
	// The graphics of the Chip 8 are black and white and the screen has a total of 2048 pixels (64 x 32).
	Display [32][64]uint8

	// The Stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
	Stack_pointer uint8

	// The Stack is an array of 16 16-bit values,
	// used to store the address that the interpreter shoud return to when finished with a subroutine.
	Stack [16]uint16

	// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad
	Keypad [16]uint8

	DrawFlag bool
}

func (cpu *CPU) Init() {
	// Program counter starts at 0x200 (512)
	Pc := 0x200
	cpu.Pc = uint16(Pc)

	// Reset the current Opcode
	cpu.Opcode = 0

	// Reset the Stack pointer
	cpu.Stack_pointer = 0

	// Reset the index register
	cpu.I = 0

	// Clear Display
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			cpu.Display[i][j] = 0
		}
	}

	// Clear stack
	for i := 0; i < len(cpu.Stack); i++ {
		cpu.Stack[i] = 0
	}

	// Clear register V0-VF
	for i := 0; i < 16; i++ {
		cpu.V[i] = 0
	}

	// Clear Memory
	for i := 0; i < len(cpu.Memory); i++ {
		cpu.Memory[i] = 0
	}

	// Load fontSet
	for i := 0; i < 80; i++ {
		cpu.Memory[i] = fontSet[i]
	}

	// Reset the delay_timer and the sound_timer registers
	cpu.Delay_timer = 0
	cpu.Sound_timer = 0
}

func (cpu *CPU) EmulateCycle() {
	// Emulation cycle: Fetch -> Decode -> Execute
	// Every cycle, the method EmulateCycle is called which emulates one cycle of the Chip 8 CPU.
	// During this cycle, the emulator will Fetch, Decode and Execute one Opcode.

	// Fetch Opcode
	// One way of doing that is this:
	// 0xF0 is 1111 0000 in binary, so by doing a bitwise AND operation
	// it preserves the leftmoSound_timer 4 bits and sets the rightmoSound_timer 4 bits to 0.
	// 0x0F is 0000 1111 in binary, so by doing a bitwise AND operation
	// it preserves the rightmoSound_timer 4 bits and sets the leftmoSound_timer 4 bits to 0.
	// Finally, the two results are combined using bitwise OR to form the 16-bit Opcode value.
	// cpu.Opcode = uint16(cpu.Memory[cpu.pc]&0xF0) | uint16(cpu.Memory[cpu.pc+1]&0x0F)
	// Or, you can simply shift left the cpu.Memory address and then perform an OR operation with the new addr.
	cpu.Opcode = (uint16(cpu.Memory[cpu.Pc]) << 8) | uint16(cpu.Memory[cpu.Pc +1])

	// Decode Opcode
	// As we have stored our current Opcode, we need to decode the Opcode and
	// check the Opcode table to see what it means.
	switch cpu.Opcode & 0xF000 { // 0xF000 is 1111 0000 0000 0000 in binary
	case 0x0000:
		switch cpu.Opcode & 0x000F { // 0x000F is 0000 0000 0000 1111
		case 0x0000: // 0x00E0: Clears the screen
			for i := 0; i < 64; i++ {
				for j := 0; j < 32; j++ {
					cpu.Display[i][j] = 0x0
				}
			}
			cpu.Pc = cpu.Pc + 2
		case 0x000E: // 0x00EE: Returns from subroutine
			cpu.Stack_pointer = cpu.Stack_pointer - 1
			cpu.Pc = cpu.Stack[cpu.Stack_pointer]
			cpu.Pc = cpu.Pc + 2
		default:
			fmt.Printf("Unknown Opcode [0x0000]: 0x%X\n", cpu.Opcode)
		}

	case 0x1000: // 1NNN: Jumps to address NNN
		cpu.Pc = cpu.Opcode & 0x0FFF

	case 0x2000: // 2NNN: Calls subroutine at NNN.
		cpu.Stack[cpu.Pc] = cpu.Pc
		cpu.Stack_pointer = cpu.Stack_pointer + 1
		cpu.Pc = cpu.Opcode & 0x0FFF

	case 0x3000: // 3XNN: Skips the next instruction if VX equals NN
		if uint16(cpu.V[(cpu.Opcode&0x0F00)>>8]) == (cpu.Opcode & 0x00FF) {
			cpu.Pc = cpu.Pc + 4 // Skip next instruction
		} else {
			cpu.Pc = cpu.Pc + 2
		}

	case 0x4000: // 4XNN: Skips the next instruction if VX does not equal NN
		if uint16(cpu.V[(cpu.Opcode&0x0F00)>>8]) != (cpu.Opcode & 0x00FF) {
			cpu.Pc = cpu.Pc + 4 // Skip next instruction
		} else {
			cpu.Pc = cpu.Pc + 2
		}

	case 0x5000: // 5XY0: Skips the next instruction if VX equals VY
		if uint16(cpu.V[(cpu.Opcode&0x0F00)>>8]) == uint16(cpu.V[(cpu.Opcode&0x00F0)>>4]) {
			cpu.Pc = cpu.Pc + 4 // Skip next instruction
		} else {
			cpu.Pc = cpu.Pc + 2
		}

	case 0x6000: // 6XNN: Sets VX to NN.
		cpu.V[(cpu.Opcode&0x0F00)>>8] = uint8(cpu.Opcode & 0x00FF)
		cpu.Pc = cpu.Pc + 2

	case 0x7000: // 7XNN: Adds NN to VX (carry flag is not changed).
		cpu.V[(cpu.Opcode&0x0F00)>>8] += uint8(cpu.Opcode & 0x00FF)
		cpu.Pc = cpu.Pc + 2

    // Chip-8 ALU (arithmetic logic unit)
    // Performs arithmetic and bitwise operations.
	case 0x8000:
		switch cpu.Opcode & 0x000F { // 0x000F is 0000 0000 0000 1111
		case 0x0000: // 8XY0: Sets Vx to the value of Vy
			cpu.V[(cpu.Opcode&0x0F00)>>4] = cpu.V[(cpu.Opcode & 0x00F0)]
			cpu.Pc = cpu.Pc + 2

		case 0x0001: // 8XY1: Sets VX to VX or VY. (bitwise OR operation)
			cpu.V[(cpu.Opcode&0x0F00)>>4] = cpu.V[(cpu.Opcode&0x0F00)>>4] | cpu.V[(cpu.Opcode&0x00F0)]
			cpu.Pc = cpu.Pc + 2

		case 0x0002: // 8XY2: Sets VX to VX and VY. (bitwise AND operation)
			cpu.V[(cpu.Opcode&0x0F00)>>4] = cpu.V[(cpu.Opcode&0x0F00)>>4] & cpu.V[(cpu.Opcode&0x00F0)]
			cpu.Pc = cpu.Pc + 2

		case 0x0003: // 8XY3: Sets VX to VX xor VY. (bitwise XOR operation)
			cpu.V[(cpu.Opcode&0x0F00)>>4] = cpu.V[(cpu.Opcode&0x0F00)>>4] ^ cpu.V[(cpu.Opcode&0x00F0)]
			cpu.Pc = cpu.Pc + 2

		case 0x0004: // 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.
			if cpu.V[(cpu.Opcode&0x00F0)>>4] > (0xFF - cpu.V[(cpu.Opcode&0x0F00)>>8]) {
				cpu.V[0xF] = 1 //carry
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[(cpu.Opcode&0x0F00)>>8] += cpu.V[(cpu.Opcode&0x00F0)>>4]
			cpu.Pc = cpu.Pc + 2 // Because every instruction is 2 bytes long

		case 0x0005: // 8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.
			if cpu.V[(cpu.Opcode&0x00F0)>>4] > cpu.V[(cpu.Opcode&0x0F00)>>8] {
				cpu.V[0xF] = 0 // There is a borrow
			} else {
				cpu.V[0xF] = 1 // No borrow
			}
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.V[(cpu.Opcode&0x0F00)>>8] - cpu.V[(cpu.Opcode&0x00F0)>>4]
			cpu.Pc = cpu.Pc + 2 // Because every instruction is 2 bytes long

		case 0x0006: // 0x8XY6 Shifts VY right by one and stores the result to VX (VY remains unchanged). VF is set to the value of the leaSound_timer significant bit of VY before the shift
			cpu.V[0xF] = cpu.V[(cpu.Opcode&0x0F00)>>8] & 0x1
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.V[(cpu.Opcode&0x0F00)>>8] >> 1
			cpu.Pc = cpu.Pc + 2
		case 0x0007: // 0x8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if cpu.V[(cpu.Opcode&0x0F00)>>8] > cpu.V[(cpu.Opcode&0x00F0)>>4] {
				cpu.V[0xF] = 0
			} else {
				cpu.V[0xF] = 1
			}
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.V[(cpu.Opcode&0x00F0)>>4] - cpu.V[(cpu.Opcode&0x0F00)>>8]
			cpu.Pc = cpu.Pc + 2
		case 0x000E: // 0x8XYE Shifts VY left by one and copies the result to VX. VF is set to the value of the moSound_timer significant bit of VY before the shift
			cpu.V[0xF] = cpu.V[(cpu.Opcode&0x0F00)>>8] >> 7
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.V[(cpu.Opcode&0x0F00)>>8] << 1
			cpu.Pc = cpu.Pc + 2
		default:
			fmt.Printf("Invalid Opcode %X\n", cpu.Opcode)
		}

	case 0x9000: // 9XY0: Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);
		if cpu.V[(cpu.Opcode&0x0F00)>>8] != cpu.V[(cpu.Opcode&0x00F0)>>4] {
			cpu.Pc = cpu.Pc + 4 // Skip the next instruction
		} else {
			cpu.Pc = cpu.Pc + 2 // Go to the rightmoSound_timer instruction
		}

	case 0xA000: // ANNN: Sets I to the address NNN
		cpu.I = cpu.Opcode & 0x0FFF
		cpu.Pc = cpu.Pc + 2 // Because every instruction is 2 bytes long

	case 0xB000: // BNNN: Jumps to the address NNN plus V0.
		cpu.Pc = (cpu.Opcode & 0x0FFF) + uint16(cpu.V[0x0])

	case 0xC000: // CXNN: Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		cpu.V[(cpu.Opcode&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8((cpu.Opcode & 0x00FF))
		cpu.Pc = cpu.Pc + 2

	case 0xD000: // DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
		x := cpu.V[(cpu.Opcode&0x0F00)>>8]
		y := cpu.V[(cpu.Opcode&0x00F0)>>4]
		height := cpu.V[(cpu.Opcode & 0x000F)]
		var pixel uint8
		var xline uint16
		var yline uint16

		cpu.V[0xF] = 0
		for yline = 0; yline < uint16(height); yline++ {
			pixel = cpu.Memory[cpu.I+yline]
			for xline = 0; xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if (pixel & (0x80 >> xline)) != 0 {
						if cpu.Display[(y + uint8(yline))][x+uint8(xline)] == 1 {
							cpu.V[0xF] = 1
						}
						cpu.Display[(y + uint8(yline))][x+uint8(xline)] ^= 1
					}
				}
			}
		}
		cpu.DrawFlag = true
		cpu.Pc = cpu.Pc + 2

	case 0xE000:
		switch cpu.Opcode & 0x00FF {
		case 0x009E: // 0xEX9E Skips the next instruction if the key stored in VX is pressed
			if cpu.Keypad[cpu.V[(cpu.Opcode&0x0F00)>>8]] != 0 {
				cpu.Pc = cpu.Pc + 4
			} else {
				cpu.Pc = cpu.Pc + 2
			}
		case 0x00A1: // 0xEXA1 Skips the next instruction if the key stored in VX isn't pressed
			if cpu.Keypad[cpu.V[(cpu.Opcode&0x0F00)>>8]] == 0 {
				cpu.Pc = cpu.Pc + 4
			} else {
				cpu.Pc = cpu.Pc + 2
			}
		default:
			fmt.Printf("Invalid Opcode %X\n", cpu.Opcode)
		}

	case 0xF000:
		switch cpu.Opcode & 0x00FF { // 0x000F is 0000 0000 0000 1111

		case 0x0007: // FX07: Sets VX to the value of the delay timer.
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.Delay_timer
			cpu.Pc = cpu.Pc + 2

		case 0x000A: // FX0A: A key press is awaited, and then stored in VX (blocking operation, all instruction halted until next key event).
			pressed := false
			for i := 0; i < len(cpu.Keypad); i++ {
				if cpu.Keypad[i] != 0 {
					cpu.V[(cpu.Opcode&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			cpu.Pc = cpu.Pc + 2

		case 0x0015: // FX15: Sets the delay timer to VX.
			cpu.Delay_timer = cpu.V[(cpu.Opcode&0x0F00)>>8]
			cpu.Pc = cpu.Pc + 2

		case 0x0018: // FX18: Sets the sound timer to VX.
			cpu.Sound_timer = cpu.V[(cpu.Opcode&0x0F00)>>8]
			cpu.Pc = cpu.Pc + 2

		case 0x001E: // FX1E: Adds VX to I. VF is not affected.
			cpu.I = cpu.I + uint16(cpu.V[(cpu.Opcode&0x0F00)>>8])
			cpu.Pc = cpu.Pc + 2

		case 0x0029: // FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			cpu.I = uint16(cpu.V[(cpu.Opcode&0x0F00)>>8]) * 0x5
			cpu.Pc = cpu.Pc + 2

		case 0x0033: // FX33: Stores the binary-coded decimal representation of VX, with the hundreds digit in Memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
			cpu.Memory[cpu.I] = cpu.V[(cpu.Opcode&0x0F00)>>8] / 100
			cpu.Memory[cpu.I+1] = (cpu.V[(cpu.Opcode&0x0F00)>>8] / 10) % 10
			cpu.Memory[cpu.I+2] = (cpu.V[(cpu.Opcode&0x0F00)>>8] % 100) % 10
			cpu.Pc = cpu.Pc + 2 // Because every instruction is 2 bytes long

		case 0x0055: // FX55: Stores from V0 to VX (including VX) in Memory, starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.
			cpu.Memory[cpu.I] = uint8(cpu.V[(cpu.Opcode&0x0F00)>>8])
			for i := uint16(0); i <= ((cpu.Opcode & 0x0F00) >> 8); i++ {
				cpu.Memory[cpu.I+i] = cpu.V[i]
			}
			cpu.Pc = cpu.Pc + 2

		case 0x0065: // FX65: Fills from V0 to VX (including VX) with values from Memory, starting at address I. The offset from I is increased by 1 for each value read, but I itself is left unmodified.
			cpu.V[(cpu.Opcode&0x0F00)>>8] = cpu.Memory[cpu.I]
			for i := uint16(0); i <= ((cpu.Opcode & 0x0F00) >> 8); i++ {
				cpu.V[i] = cpu.Memory[cpu.I+i]
			}
			cpu.Pc = cpu.Pc + 2

		default:
			fmt.Printf("Unknown Opcode [0x8000]: 0x%X\n", cpu.Opcode)
		}

	default:
		fmt.Printf("Unknown Opcode: 0x%x\n", cpu.Opcode)
	}

	// Update timers
	if cpu.Delay_timer > 0 {
		cpu.Delay_timer = cpu.Delay_timer - 1
	}

	if cpu.Sound_timer > 0 {
		if cpu.Sound_timer == 1 {
			fmt.Println("BEEP!")
			cpu.Sound_timer = cpu.Sound_timer - 1
		}
	}
}

func (cpu *CPU) LoadRom(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Load the ROM into Memory
	for i := 0; i < len(data); i++ {
		cpu.Memory[i+512] = data[i]
	}

	fmt.Println("ROM loaded successfully")
}
