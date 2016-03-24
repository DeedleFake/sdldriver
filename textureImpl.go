package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"unsafe"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type textureImpl struct {
	tex *C.SDL_Texture
}

func (t *textureImpl) Release() {
	C.SDL_DestroyTexture(t.tex)
	t.tex = nil
}

func (t *textureImpl) Size() image.Point {
	var format C.Uint32
	var a, w, h C.int
	C.SDL_QueryTexture(t.tex,
		&format,
		&a,
		&w,
		&h,
	)

	return image.Pt(int(w), int(h))
}

func (t *textureImpl) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: t.Size(),
	}
}

func (t *textureImpl) lock(r image.Rectangle) ([]uint8, int) {
	var pix *uint8
	var pitch C.int
	C.SDL_LockTexture(t.tex,
		toSDLRect(r),
		&pix,
		&pitch,
	)

	return *(*[]uint8)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(pix)),
		Len:  r.Dx() * r.Dy(),
		Cap:  r.Dx() * r.Dy(),
	})), int(pitch)
}

func (t *textureImpl) unlock() {
	C.SDL_UnlockTexture(t.tex)
}

func (t *textureImpl) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	pix, pitch := t.lock(image.Rectangle{
		Min: dp,
		Max: sr.Size(),
	}.Intersect(t.Bounds()))
	defer t.unlock()

	rgba := src.RGBA()

	for row := 0; row < rgba.Bounds().Dy(); row++ {
		copy(
			pix[row*pitch:row*pitch+pitch],
			rgba.Pix[(row-rgba.Bounds().Min.Y)*rgba.Stride+(sr.Min.X*4):(row-rgba.Bounds().Min.Y)*rgba.Stride+(sr.Max.X*4)],
		)
	}
}

func (t *textureImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	dr = dr.Intersect(t.Bounds())

	pix, pitch := t.lock(dr)
	defer t.unlock()

	r, g, b, a := src.RGBA()
	rgba := [...]uint32{r, g, b, a}

	var c int
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			pix[(y*pitch)+x] = uint8(rgba[c] * 255 / 0xFFFF)

			c++
			if c >= len(rgba) {
				c = 0
			}
		}
	}
}
