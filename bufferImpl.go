package sdldriver

import (
	"image"
	"reflect"
	"unsafe"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type bufferImpl struct {
	sur *C.SDL_Surface
}

func (b *bufferImpl) Release() {
	C.SDL_FreeSurface(b.sur)
	b.sur = nil
}

func (b *bufferImpl) Size() image.Point {
	return image.Pt(
		int(b.sur.w),
		int(b.sur.h),
	)
}

func (b *bufferImpl) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: b.Size(),
	}
}

func (b *bufferImpl) RGBA() *image.RGBA {
	return &image.RGBA{
		Pix: *(*[]uint8)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(b.sur.pixels)),
			Len:  int(b.sur.w * b.sur.h),
			Cap:  int(b.sur.w * b.sur.h),
		})),
		Stride: 1,
		Rect:   b.Bounds(),
	}
}
