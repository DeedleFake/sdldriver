package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type sdlError string

func (err sdlError) Error() string {
	return string(err)
}

func Main(main func(screen.Screen)) {
	C.SDL_Init(C.SDL_INIT_EVERYTHING)
	defer C.SDL_Quit()

	main(new(screenImpl))
}
