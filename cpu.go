package chip8 

// Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM,
// from location 0x000 (0) to 0xFFF
// Most Chip-8 programs start at location 0x200 (512)
// and 16 general-purpose 8-bit registers named V0 to VF
var Memory [4096]uint8

// The Chip 8 has 35 opcodes which are all two bytes long.
var Opcodes [35]uint16

// Chip-8 has 16 general purpose 8-bit registers, usually referred to as Vx,
// where x is a hexadecimal digit (0 through F).
var Registers [16]uint8

// The stack is an array of 16 16-bit values,
// used to store the address that the interpreter shoud return to when finished with a subroutine.
var Stack [16]uint16

//  ProgramCounter should be 16-bit, and is used to store the currently executing address.
var ProgramCounter uint16 = 0x200

// The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
var StackPointer uint8 = 0x000

// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad
var Keypad [16]uint8 

// When these registers are non-zero, 
// they are automatically decremented at a rate of 60Hz.
// The system’s buzzer sounds whenever the sound timer reaches zero.
var delay_timer uint16 = 0
var sound_timer uint16 = 0

// Graphics:
// The graphics of the Chip 8 are black and white and the screen has a total of 2048 pixels (64 x 32).
var Display [64][32]uint8
