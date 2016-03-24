package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type screenImpl struct {
	r *C.SDL_Renderer
}

func (s *screenImpl) NewBuffer(size image.Point) (screen.Buffer, error) {
	sur := C.SDL_CreateRGBSurface(0,
		C.int(size.X),
		C.int(size.Y),
		32,
		0xFF000000,
		0x00FF0000,
		0x0000FF00,
		0x000000FF,
	)

	return &bufferImpl{
		sur: sur,
	}, nil
}

func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error)

func (s *screenImpl) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error)
