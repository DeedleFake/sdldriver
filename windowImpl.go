package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"image"
	"image/color"
	"image/draw"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type windowImpl struct {
	win *C.SDL_Window
	ren *C.SDL_Renderer
}

func (w *windowImpl) Release()

func (w *windowImpl) Send(ev interface{})

func (w *windowImpl) NextEvent() interface{}

func (w *windowImpl) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle)

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Publish() screen.PublishResult
