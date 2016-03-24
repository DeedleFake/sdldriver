package sdldriver

import (
	"image"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type bufferImpl struct {
	sur *C.SDL_Surface
}

func (b *bufferImpl) Release()

func (b *bufferImpl) Size() image.Point

func (b *bufferImpl) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: b.Size(),
	}
}

func (b *bufferImpl) RGBA() *image.RGBA
