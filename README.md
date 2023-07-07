# Chip-8 Emulator in Go

According to Wikipedia:
> CHIP-8 is an interpreted programming language, developed by Joseph Weisbecker made on his 1802 Microprocessor. It was initially used on the COSMAC VIP and Telmac 1800 8-bit microcomputers in the mid-1970s.

This is a simple Chip-8 emulator written in Go. The Chip-8 is an interpreted programming language used in many early computer systems and video game consoles.

## Installation

To run the emulator, make sure you have Go installed on your system. Then, follow these steps:

1. Clone this repository: ```git clone https://github.com/petersid2022/chip8.git```
2. Navigate to the project directory: ```cd chip8```
3. Build the project: ```go build```
4. Run the emulator: ```./chip8```

## Key Bindings

```
Chip8 keypad         Keyboard mapping
1 | 2 | 3 | C        1 | 2 | 3 | 4
4 | 5 | 6 | D   =>   Q | W | E | R
7 | 8 | 9 | E   =>   A | S | D | F
A | 0 | B | F        Z | X | C | V
```


## Resources

If you're interested in learning more about the Chip-8 system check out the following resources:

* [Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM): A comprehensive guide to the Chip-8 system.
* [Chip-8 Wikipedia page](https://en.wikipedia.org/wiki/CHIP-8): General information about the Chip-8 system.
* [How to write an emulator (CHIP-8 interpreter)](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)

## TODO

* Add the [XO-CHIP](https://johnearnest.github.io/Octo/docs/XO-ChipSpecification.html) extension, which includes:
    1. 7 new opcodes
    2. 16-bit addressing for a total of ~64kb RAM
    3. Second display buffer allowing for 4 colors instead of the typical 2
    4. Improved sound support
    5. Modified Fx75 and Fx85 instructions to allow for 16 user flags instead of typical 8

## License

(The MIT License)

Copyright (c) 2023 Peter Sideris petersid2022@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
