package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type textureImpl struct {
	tex *C.SDL_Texture
}

func (t *textureImpl) Release()

func (t *textureImpl) Size() image.Point

func (t *textureImpl) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: t.Size(),
	}
}

func (t *textureImpl) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle)

func (t *textureImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)
