package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"strconv"
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

func toSDLColor(format *C.SDL_PixelFormat, c color.Color, op draw.Op) C.Uint32 {
	r, g, b, a := c.RGBA()
	r *= 255 / 0xFFFF
	g *= 255 / 0xFFFF
	b *= 255 / 0xFFFF
	a *= 255 / 0xFFFF

	switch op {
	case draw.Src:
		return C.SDL_MapRGB(format, C.Uint8(r), C.Uint8(g), C.Uint8(b))
	case draw.Over:
		return C.SDL_MapRGBA(format, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a))
	}

	panic("Unknown op: " + strconv.FormatInt(int64(op), 10))
}

func Main(main func(screen.Screen)) {
	C.SDL_Init(C.SDL_INIT_EVERYTHING)
	defer C.SDL_Quit()

	main(new(screenImpl))
}
