package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

func Main(main func(screen.Screen)) {
	C.SDL_Init(C.SDL_INIT_EVERYTHING)
	defer C.SDL_Quit()

	main(new(screenImpl))
}
