package sdldriver

import (
	"errors"
	"golang.org/x/exp/shiny/screen"
	"image"
	"runtime"
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
	if sur == nil {
		return nil, sdlError(C.GoString(C.SDL_GetError()))
	}

	b := &bufferImpl{
		sur: sur,
	}
	runtime.SetFinalizer(b, (*bufferImpl).Release)

	return b, nil
}

func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	if s.r == nil {
		return nil, errors.New("No renderer has been created yet; create a window first")
	}

	tex := C.SDL_CreateTexture(
		s.r,
		C.SDL_PIXELFORMAT_RGBA8888,
		C.SDL_TEXTUREACCESS_STREAMING,
		C.int(size.X),
		C.int(size.Y),
	)
	if tex == nil {
		return nil, sdlError(C.GoString(C.SDL_GetError()))
	}

	t := &textureImpl{
		tex: tex,
	}
	runtime.SetFinalizer(t, (*textureImpl).Release)

	return t, nil
}

func (s *screenImpl) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	var win *C.SDL_Window
	var ren *C.SDL_Renderer
	ok := C.SDL_CreateWindowAndRenderer(
		C.int(opts.Width),
		C.int(opts.Height),
		C.SDL_WINDOW_ALLOW_HIGHDPI,
		&win,
		&ren,
	)
	if ok != 0 {
		return nil, sdlError(C.GoString(C.SDL_GetError()))
	}

	s.r = ren

	w := &windowImpl{
		win: win,
		ren: ren,
	}
	runtime.SetFinalizer(w, (*windowImpl).Release)

	return w, nil
}
