package main

import (
	"fmt"
	"os"

	"github.com/petersid2022/chip8/cmd"
	sdl "github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "SDL2 GFX"
var winWidth, winHeight int32 = 800, 600

func mapKey(sdlKey sdl.Keycode) int {
	switch sdlKey {
	case sdl.K_1:
		return 0x1
	case sdl.K_2:
		return 0x2
	case sdl.K_3:
		return 0x3
	case sdl.K_4:
		return 0xC
	case sdl.K_q:
		return 0x4
	case sdl.K_w:
		return 0x5
	case sdl.K_e:
		return 0x6
	case sdl.K_r:
		return 0xD
	case sdl.K_a:
		return 0x7
	case sdl.K_s:
		return 0x8
	case sdl.K_d:
		return 0x9
	case sdl.K_f:
		return 0xE
	case sdl.K_z:
		return 0xA
	case sdl.K_x:
		return 0x0
	case sdl.K_c:
		return 0xB
	case sdl.K_v:
		return 0xF
	default:
		return -1 // Invalid key
	}
}

func run() int {
	// Initialize the Chip8 system and load the game into the memory
	chip8 := chip8.CPU{}
	chip8.Init()
    arg := os.Args[1]
	chip8.LoadRom("/home/petrside/github/chip8/" + arg)

	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	// Setting up graphics and creating a window
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		return 1
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 2
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 3 // don't use os.Exit(3); otherwise, previous deferred calls will never run
	}
	renderer.Clear()
	defer renderer.Destroy()

	// Initialize the key states array
	keyStates := &[16]bool{}
	// Emulation loop
	for {
		// Handle keyboard events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				// Handle key down event
				if t.Type == sdl.KEYDOWN {
					// Map the keyboard key to the corresponding Chip8 keypad key
					chip8Key := mapKey(t.Keysym.Sym)

					// Set the corresponding key state in the keyStates array
					if chip8Key != -1 {
						(*keyStates)[chip8Key] = true
					}
				} else if t.Type == sdl.KEYUP {
					// Map the keyboard key to the corresponding Chip8 keypad key
					chip8Key := mapKey(t.Keysym.Sym)

					// Clear the corresponding key state in the keyStates array
					if chip8Key != -1 {
						(*keyStates)[chip8Key] = false
					}
				}
			}
		}

		// Emulate one cycle
		chip8.EmulateCycle()

		// If the draw flag is set, update the screen
		if chip8.DrawFlag {
			// Draw graphics
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()

			windowWidth, windowHeight := window.GetSize()
			pixelWidth := windowWidth / 64
			pixelHeight := windowHeight / 32

			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					if chip8.Display[i][j] == 1 {
						renderer.SetDrawColor(255, 255, 255, 255)
					} else {
						renderer.SetDrawColor(0, 0, 0, 255)
					}
					renderer.FillRect(&sdl.Rect{
						X: int32(j) * pixelWidth,
						Y: int32(i) * pixelHeight,
						W: pixelWidth,
						H: pixelHeight,
					})
				}
			}

			renderer.Present()

			// Reset the draw flag
			chip8.DrawFlag = false
		}

		// Store key press state (Press and Release)
		chip8.SetKeys(*keyStates)

		// Delay to control the emulation speed
		sdl.Delay(100 / 60)
	}
}

func main() {
	os.Exit(run())
}
