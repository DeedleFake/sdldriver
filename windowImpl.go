package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"image"
	"image/color"
	"image/draw"
	"time"
	"unsafe"
)

// #cgo pkg-config: sdl2
// #include <SDL.h>
import "C"

type windowImpl struct {
	win *C.SDL_Window
	ren *C.SDL_Renderer
}

func (w *windowImpl) Release() {
	C.SDL_DestroyRenderer(w.ren)
	w.ren = nil

	C.SDL_DestroyWindow(w.win)
	w.win = nil
}

func (w *windowImpl) Send(ev interface{}) {
	t := C.SDL_RegisterEvents(1)
	C.SDL_PushEvent((*C.SDL_Event)(unsafe.Pointer(&C.SDL_UserEvent{
		_type:     t,
		timestamp: C.Uint32(time.Now().Unix()),
		windowID:  C.SDL_GetWindowID(w.win),
		code:      0,
		data1:     unsafe.Pointer(&ev),
	})))
}

func (w *windowImpl) windowEvent(ev *C.SDL_WindowEvent) interface{} {
	switch ev.event {
	case C.SDL_WINDOWEVENT_SHOWN:
		return lifecycle.Event{
			From: lifecycle.StageAlive,
			To:   lifecycle.StageVisible,
		}

	case C.SDL_WINDOWEVENT_HIDDEN:
		return lifecycle.Event{
			From: lifecycle.StageVisible,
			To:   lifecycle.StageAlive,
		}
	}

	return nil
}

func (w *windowImpl) keyEvent(ev *C.SDL_KeyboardEvent, dir key.Direction) interface{}

func (w *windowImpl) mouseMotionEvent(ev *C.SDL_MouseMotionEvent) interface{}

func (w *windowImpl) mouseButtonEvent(ev *C.SDL_MouseButtonEvent, dir mouse.Direction) interface{}

func (w *windowImpl) NextEvent() interface{} {
top:
	var ev C.SDL_Event
	C.SDL_WaitEvent(&ev)
	switch (*C.SDL_CommonEvent)(unsafe.Pointer(&ev))._type {
	case C.SDL_QUIT:
		return lifecycle.Event{
			From: lifecycle.StageAlive,
			To:   lifecycle.StageDead,
		}

	case C.SDL_WINDOWEVENT:
		r := w.windowEvent((*C.SDL_WindowEvent)(unsafe.Pointer(&ev)))
		if r == nil {
			goto top
		}
		return r

	case C.SDL_KEYUP:
		return w.keyEvent((*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev)), key.DirRelease)
	case C.SDL_KEYDOWN:
		return w.keyEvent((*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev)), key.DirPress)

	case C.SDL_MOUSEMOTION:
		return w.mouseMotionEvent((*C.SDL_MouseMotionEvent)(unsafe.Pointer(&ev)))
	case C.SDL_MOUSEBUTTONUP:
		return w.mouseButtonEvent((*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev)), mouse.DirRelease)
	case C.SDL_MOUSEBUTTONDOWN:
		return w.mouseButtonEvent((*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev)), mouse.DirPress)
	}

	goto top
}

func (w *windowImpl) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle)

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions)

func (w *windowImpl) Publish() screen.PublishResult
