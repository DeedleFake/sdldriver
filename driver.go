package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type sdlError string

func (err sdlError) Error() string {
	return string(err)
}

func toSDLRect(r image.Rectangle) *C.SDL_Rect {
	if r == image.ZR {
		return nil
	}

	return &C.SDL_Rect{
		x: C.int(r.Min.X),
		y: C.int(r.Min.Y),
		w: C.int(r.Dx()),
		h: C.int(r.Dy()),
	}
}

func Main(main func(screen.Screen)) {
	C.SDL_Init(C.SDL_INIT_EVERYTHING)
	defer C.SDL_Quit()

	main(new(screenImpl))
}
