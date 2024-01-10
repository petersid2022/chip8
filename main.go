package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/petersid2022/chip8/cmd"
	sdl "github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	winTitle            string = "CHIP8 emulator"
	winWidth, winHeight int32  = 800, 600
	delay               uint32 = 100
	target_fps          uint32 = 60
)

//go:embed font.ttf
var contentfont embed.FS
var fontSize = 24

type MenuItem struct {
	Text   string
	Bounds sdl.Rect
}

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

//go:embed roms
var content embed.FS

func showMenu(renderer *sdl.Renderer, font *ttf.Font) string {
	files, err := content.ReadDir("roms")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read ROM directory: %s\n", err)
		return ""
	}

	menuItems := make([]MenuItem, len(files))
	lineHeight := fontSize + 10

	const itemsPerColumn = 10
	columnWidth := winWidth / 4
	numColumns := (len(files) + itemsPerColumn - 1) / itemsPerColumn
	columnSpacing := (winWidth - int32(numColumns)*columnWidth) / (int32(numColumns) + 1)

	for i, file := range files {
		columnIndex := i / itemsPerColumn
		itemIndex := i % itemsPerColumn
		itemText := fmt.Sprintf("%d) %s", i+1, file.Name())
		itemRect := sdl.Rect{
			X: (int32(columnIndex) * (columnWidth + columnSpacing)) + columnSpacing,
			Y: 96 + (int32(lineHeight) * int32(itemIndex)),
			W: columnWidth,
			H: int32(lineHeight),
		}
		menuItems[i] = MenuItem{Text: itemText, Bounds: itemRect}
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return ""
			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN {
					for i, item := range menuItems {
						if t.X >= item.Bounds.X && t.X < item.Bounds.X+item.Bounds.W &&
							t.Y >= item.Bounds.Y && t.Y < item.Bounds.Y+item.Bounds.H {
							return files[i].Name()
						}
					}
				}
			case *sdl.KeyboardEvent:
				// Handle key down event
				if t.Type == sdl.KEYDOWN {
					// Exit the game if the "Backspace" key is pressed
					if t.Keysym.Sym == sdl.K_ESCAPE {
						// restart the game
						fmt.Println("Exiting")
						return ""
					}
					if t.Keysym.Sym == sdl.K_i {
						// increment by -5 target_fps
						// until target_fps is 0
						if target_fps > 0 {
							target_fps -= 5
						}
					}
					if t.Keysym.Sym == sdl.K_p {
						// increment by +5 target_fps
						if target_fps < 100 {
							target_fps += 5
						}
					}
					if t.Keysym.Sym == sdl.K_j {
						// increment by -100 delay
						// if delay is 100 do nothing
						if delay > 100 {
							delay -= 100
						}
					}
					if t.Keysym.Sym == sdl.K_l {
						// increment by +100 delay
						// if delay is 1000 do nothing
						if delay < 1000 {
							delay += 100
						}
					}
				}

			}
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// Render "Select ROM to play" text
		// -----------------------------
		// -----------------------------
		// -----------------------------

		textSurface, err := font.RenderUTF8Solid("Click on a ROM to play", sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			return ""
		}
		defer textSurface.Free()
		textTexture, err := renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return ""
		}
		defer textTexture.Destroy()
		_, _, textWidth, textHeight, err := textTexture.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			return ""
		}
		textX := (winWidth - 2*textWidth) / 2
		textY := (96 - 2*textHeight) / 2
		renderer.Copy(textTexture, nil, &sdl.Rect{X: textX, Y: textY, W: textWidth * 2, H: textHeight * 2})

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// DELAY TEXT
		// -----------------------------
		// -----------------------------
		// -----------------------------

		delaySurface, err := font.RenderUTF8Solid(fmt.Sprintf("delay: %d (j: -100, l: +100)", delay), sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			return ""
		}
		defer delaySurface.Free()
		delayTexture, err := renderer.CreateTextureFromSurface(delaySurface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return ""
		}
		defer delayTexture.Destroy()
		_, _, delayWidth, delayHeight, err := delayTexture.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			return ""
		}
		// delayX := (winWidth - delayWidth) / 8
		delayY := int32(winHeight - delayHeight - 8)
		renderer.Copy(delayTexture, nil, &sdl.Rect{X: columnSpacing, Y: delayY, W: delayWidth, H: delayHeight})

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// TARGET_FPS TEXT
		// -----------------------------
		// -----------------------------
		// -----------------------------

		target_fpsSurface, err := font.RenderUTF8Solid(fmt.Sprintf("target_fps: %d (i: -5, p: +5)", target_fps), sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			return ""
		}
		defer target_fpsSurface.Free()
		target_fpsTexture, err := renderer.CreateTextureFromSurface(target_fpsSurface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return ""
		}
		defer target_fpsTexture.Destroy()
		_, _, target_fpsWidth, target_fpsHeight, err := target_fpsTexture.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			return ""
		}
		// target_fpsX := (winWidth - target_fpsWidth) / 8
		target_fpsY := int32(winHeight - target_fpsHeight - 8 - delayHeight - 8)
		renderer.Copy(target_fpsTexture, nil, &sdl.Rect{X: columnSpacing, Y: target_fpsY, W: target_fpsWidth, H: target_fpsHeight})

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// Render "Credits" text
		// -----------------------------
		// -----------------------------
		// -----------------------------

		creditsSurface, err := font.RenderUTF8Solid("(c) Peter Sideris 2023", sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			return ""
		}
		defer creditsSurface.Free()
		creditsTexture, err := renderer.CreateTextureFromSurface(creditsSurface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return ""
		}
		defer creditsTexture.Destroy()
		_, _, creditsWidth, creditsHeight, err := creditsTexture.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			return ""
		}
		// creditsX := (winWidth - creditsWidth) / 2
		creditsX := (winWidth - columnSpacing - creditsWidth)
		creditsY := int32(winHeight - delayHeight - 8)
		renderer.Copy(creditsTexture, nil, &sdl.Rect{X: creditsX, Y: creditsY, W: creditsWidth, H: creditsHeight})

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// Render "Escape" text
		// -----------------------------
		// -----------------------------
		// -----------------------------

		exitSurface, err := font.RenderUTF8Solid("Press <Escape> to exit.", sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			return ""
		}
		defer exitSurface.Free()
		exitTexture, err := renderer.CreateTextureFromSurface(exitSurface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return ""
		}
		defer exitTexture.Destroy()
		_, _, exitWidth, exitHeight, err := exitTexture.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			return ""
		}
		exitX := (winWidth - columnSpacing - exitWidth)
		exitY := int32(winHeight - target_fpsHeight - 8 - delayHeight - 8)
		renderer.Copy(exitTexture, nil, &sdl.Rect{X: exitX, Y: exitY, W: exitWidth, H: exitHeight})

		// -----------------------------
		// -----------------------------
		// -----------------------------
		// Render the menu items
		// -----------------------------
		// -----------------------------
		// -----------------------------

		for _, item := range menuItems {
			itemSurface, err := font.RenderUTF8Solid(item.Text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
				return ""
			}
			defer itemSurface.Free()

			itemTexture, err := renderer.CreateTextureFromSurface(itemSurface)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
				return ""
			}
			defer itemTexture.Destroy()

			renderer.Copy(itemTexture, nil, &item.Bounds)
		}

		renderer.Present()
		sdl.Delay(16)
	}
}

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	// Setting up graphics and creating a window
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		return 2
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 3
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 4
	}
	renderer.Clear()
	defer renderer.Destroy()

	if err = ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return 4
	}
	defer ttf.Quit()

	fontData, err := contentfont.ReadFile("font.ttf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read font file: %s\n", err)
		return 5
	}

	rwops, err := sdl.RWFromMem(fontData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed on rwops: %s\n", err)
		return 5
	}

	font, err := ttf.OpenFontRW(rwops, 1, fontSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font from RWops: %s\n", err)
		return 5
	}

	defer font.Close()

	romName := showMenu(renderer, font)
	if romName == "" {
		return 0
	}

	// Initialize the Chip8 system and load the game into memory
	chip8 := chip8.CPU{}
	chip8.Init()
	chip8.LoadRom("./roms/" + romName)

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
					// Restart the game if the "ESC" key is pressed
					if t.Keysym.Sym == sdl.K_BACKSPACE {
						// restart the game
						fmt.Println("Restarting")
						return 1
					}

					// Exit the game if the "Backspace" key is pressed
					if t.Keysym.Sym == sdl.K_ESCAPE {
						// restart the game
						fmt.Println("Exiting")
						return 0
					}

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
			case *sdl.QuitEvent:
				return 0
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
			footerSurface, err := font.RenderUTF8Solid("<Escape> to exit, <Backspace> to restart", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render text: %s\n", err)
			}
			defer footerSurface.Free()

			footerTexture, err := renderer.CreateTextureFromSurface(footerSurface)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			}
			defer footerTexture.Destroy()

			// Get the dimensions of the text texture
			_, _, footerWidth, footerHeight, err := footerTexture.Query()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to query texture: %s\n", err)
			}

			// Position the text at the center of the window
			footerX := (winWidth - footerWidth) / 2
			footerY := int32(winHeight - footerHeight - 4)

			// Render the text
			renderer.Copy(footerTexture, nil, &sdl.Rect{X: footerX, Y: footerY, W: footerWidth, H: footerHeight})

			renderer.Present()

			// Reset the draw flag
			chip8.DrawFlag = false
		}

		// Store key press state (Press and Release)
		chip8.SetKeys(*keyStates)

		// Delay to control the emulation speed
		sdl.Delay(uint32(delay / target_fps))
	}
}

func main() {
	os.Stdout = nil
	for {
		returnValue := run()

		if returnValue == 0 {
			break
		}
	}
}
