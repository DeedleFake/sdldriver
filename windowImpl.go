package sdldriver

import (
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
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

	stage lifecycle.Stage
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
	r := lifecycle.Event{
		From: w.stage,
	}

	switch ev.event {
	case C.SDL_WINDOWEVENT_SHOWN:
		r.To = lifecycle.StageVisible

	case C.SDL_WINDOWEVENT_HIDDEN, C.SDL_WINDOWEVENT_MINIMIZED:
		r.To = lifecycle.StageAlive

	case C.SDL_WINDOWEVENT_EXPOSED:
		return paint.Event{}

	case C.SDL_WINDOWEVENT_FOCUS_GAINED:
		r.To = lifecycle.StageFocused

	case C.SDL_WINDOWEVENT_FOCUS_LOST:
		r.To = lifecycle.StageVisible

	default:
		return nil
	}

	if r.From == r.To {
		return nil
	}

	w.stage = r.To
	return r
}

func (w *windowImpl) keyEvent(ev *C.SDL_KeyboardEvent, dir key.Direction) interface{} {
	return key.Event{
		Rune:      rune(ev.keysym.sym),
		Code:      keyMap[ev.keysym.sym],
		Modifiers: modMap(ev.keysym.mod),
		Direction: dir,
	}
}

func (w *windowImpl) mouseMotionEvent(ev *C.SDL_MouseMotionEvent) interface{} {
	return mouse.Event{
		X: float32(ev.x),
		Y: float32(ev.y),
	}
}

func (w *windowImpl) mouseButtonEvent(ev *C.SDL_MouseButtonEvent, dir mouse.Direction) interface{} {
	return mouse.Event{
		X:         float32(ev.x),
		Y:         float32(ev.y),
		Button:    mouseButtonMap[ev.button],
		Direction: dir,
	}
}

func (w *windowImpl) mouseWheelEvent(ev *C.SDL_MouseWheelEvent) interface{} {
	var button mouse.Button
	switch {
	case ev.y == 0:
		return nil
	case ev.y < 0:
		button = mouse.ButtonWheelUp
	case ev.y > 0:
		button = mouse.ButtonWheelDown
	}

	return mouse.Event{
		Button:    button,
		Direction: mouse.DirPress,
	}
}

func (w *windowImpl) NextEvent() interface{} {
	for {
		var ev C.SDL_Event
		C.SDL_WaitEvent(&ev)
		switch (*C.SDL_CommonEvent)(unsafe.Pointer(&ev))._type {
		case C.SDL_QUIT:
			r := lifecycle.Event{
				From: w.stage,
				To:   lifecycle.StageDead,
			}
			w.stage = r.To
			return r

		case C.SDL_WINDOWEVENT:
			r := w.windowEvent((*C.SDL_WindowEvent)(unsafe.Pointer(&ev)))
			if r == nil {
				continue
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
		case C.SDL_MOUSEWHEEL:
			r := w.mouseWheelEvent((*C.SDL_MouseWheelEvent)(unsafe.Pointer(&ev)))
			if r == nil {
				continue
			}
			return r
		}
	}
}

func (w *windowImpl) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	C.SDL_BlitSurface(
		src.(*bufferImpl).sur,
		toSDLRect(sr),
		C.SDL_GetWindowSurface(w.win),
		toSDLRect(image.Rectangle{Min: dp, Max: dp}),
	)
}

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	sur := C.SDL_GetWindowSurface(w.win)
	C.SDL_FillRect(
		sur,
		toSDLRect(dr),
		toSDLColor(sur.format, src, op),
	)
}

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
	panic("Not implemented.")
}

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
	C.SDL_RenderCopy(
		w.ren,
		src.(*textureImpl).tex,
		toSDLRect(sr),
		toSDLRect(image.Rectangle{
			Min: dp,
			Max: dp.Add(sr.Size()),
		}),
	)
}

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
	panic("Not implemented.")
}

func (w *windowImpl) Publish() screen.PublishResult {
	C.SDL_UpdateWindowSurface(w.win)
	return screen.PublishResult{}
}

func modMap(sdl C.Uint16) (mod key.Modifiers) {
	if sdl&C.KMOD_SHIFT != 0 {
		mod |= key.ModShift
	}
	if sdl&C.KMOD_CTRL != 0 {
		mod |= key.ModControl
	}
	if sdl&C.KMOD_ALT != 0 {
		mod |= key.ModAlt
	}
	if sdl&C.KMOD_GUI != 0 {
		mod |= key.ModMeta
	}

	return mod
}

var (
	keyMap = map[C.SDL_Keycode]key.Code{
		C.SDLK_RETURN:    key.CodeReturnEnter,
		C.SDLK_ESCAPE:    key.CodeEscape,
		C.SDLK_BACKSPACE: key.CodeDeleteBackspace,
		C.SDLK_TAB:       key.CodeTab,
		C.SDLK_SPACE:     key.CodeSpacebar,
		//C.SDLK_EXCLAIM,
		//C.SDLK_QUOTEDBL,
		//C.SDLK_HASH,
		//C.SDLK_PERCENT,
		//C.SDLK_DOLLAR,
		//C.SDLK_AMPERSAND,
		C.SDLK_QUOTE: key.CodeApostrophe,
		//C.SDLK_LEFTPAREN,
		//C.SDLK_RIGHTPAREN,
		//C.SDLK_ASTERISK,
		//C.SDLK_PLUS,
		C.SDLK_COMMA:  key.CodeComma,
		C.SDLK_MINUS:  key.CodeHyphenMinus,
		C.SDLK_PERIOD: key.CodeFullStop,
		C.SDLK_SLASH:  key.CodeSlash,
		C.SDLK_0:      key.Code0,
		C.SDLK_1:      key.Code1,
		C.SDLK_2:      key.Code2,
		C.SDLK_3:      key.Code3,
		C.SDLK_4:      key.Code4,
		C.SDLK_5:      key.Code5,
		C.SDLK_6:      key.Code6,
		C.SDLK_7:      key.Code7,
		C.SDLK_8:      key.Code8,
		C.SDLK_9:      key.Code9,
		//C.SDLK_COLON,
		C.SDLK_SEMICOLON: key.CodeSemicolon,
		//C.SDLK_LESS,
		C.SDLK_EQUALS: key.CodeEqualSign,
		//C.SDLK_GREATER,
		//C.SDLK_QUESTION,
		//C.SDLK_AT,
		C.SDLK_LEFTBRACKET:  key.CodeLeftSquareBracket,
		C.SDLK_BACKSLASH:    key.CodeBackslash,
		C.SDLK_RIGHTBRACKET: key.CodeRightSquareBracket,
		//C.SDLK_CARET,
		//C.SDLK_UNDERSCORE,
		C.SDLK_BACKQUOTE: key.CodeGraveAccent,
		C.SDLK_a:         key.CodeA,
		C.SDLK_b:         key.CodeB,
		C.SDLK_c:         key.CodeC,
		C.SDLK_d:         key.CodeD,
		C.SDLK_e:         key.CodeE,
		C.SDLK_f:         key.CodeF,
		C.SDLK_g:         key.CodeG,
		C.SDLK_h:         key.CodeH,
		C.SDLK_i:         key.CodeI,
		C.SDLK_j:         key.CodeJ,
		C.SDLK_k:         key.CodeK,
		C.SDLK_l:         key.CodeL,
		C.SDLK_m:         key.CodeM,
		C.SDLK_n:         key.CodeN,
		C.SDLK_o:         key.CodeO,
		C.SDLK_p:         key.CodeP,
		C.SDLK_q:         key.CodeQ,
		C.SDLK_r:         key.CodeR,
		C.SDLK_s:         key.CodeS,
		C.SDLK_t:         key.CodeT,
		C.SDLK_u:         key.CodeU,
		C.SDLK_v:         key.CodeV,
		C.SDLK_w:         key.CodeW,
		C.SDLK_x:         key.CodeX,
		C.SDLK_y:         key.CodeY,
		C.SDLK_z:         key.CodeZ,
		C.SDLK_CAPSLOCK:  key.CodeCapsLock,
		C.SDLK_F1:        key.CodeF1,
		C.SDLK_F2:        key.CodeF2,
		C.SDLK_F3:        key.CodeF3,
		C.SDLK_F4:        key.CodeF4,
		C.SDLK_F5:        key.CodeF5,
		C.SDLK_F6:        key.CodeF6,
		C.SDLK_F7:        key.CodeF7,
		C.SDLK_F8:        key.CodeF8,
		C.SDLK_F9:        key.CodeF9,
		C.SDLK_F10:       key.CodeF10,
		C.SDLK_F11:       key.CodeF11,
		C.SDLK_F12:       key.CodeF12,
		//C.SDLK_PRINTSCREEN,
		//C.SDLK_SCROLLLOCK,
		C.SDLK_PAUSE:        key.CodePause,
		C.SDLK_INSERT:       key.CodeInsert,
		C.SDLK_HOME:         key.CodeHome,
		C.SDLK_PAGEUP:       key.CodePageUp,
		C.SDLK_DELETE:       key.CodeDeleteForward,
		C.SDLK_END:          key.CodeEnd,
		C.SDLK_PAGEDOWN:     key.CodePageDown,
		C.SDLK_RIGHT:        key.CodeRightArrow,
		C.SDLK_LEFT:         key.CodeLeftArrow,
		C.SDLK_DOWN:         key.CodeDownArrow,
		C.SDLK_UP:           key.CodeUpArrow,
		C.SDLK_NUMLOCKCLEAR: key.CodeKeypadNumLock,
		C.SDLK_KP_DIVIDE:    key.CodeKeypadSlash,
		C.SDLK_KP_MULTIPLY:  key.CodeKeypadAsterisk,
		C.SDLK_KP_MINUS:     key.CodeKeypadHyphenMinus,
		C.SDLK_KP_PLUS:      key.CodeKeypadPlusSign,
		C.SDLK_KP_ENTER:     key.CodeKeypadEnter,
		C.SDLK_KP_1:         key.CodeKeypad1,
		C.SDLK_KP_2:         key.CodeKeypad2,
		C.SDLK_KP_3:         key.CodeKeypad3,
		C.SDLK_KP_4:         key.CodeKeypad4,
		C.SDLK_KP_5:         key.CodeKeypad5,
		C.SDLK_KP_6:         key.CodeKeypad6,
		C.SDLK_KP_7:         key.CodeKeypad7,
		C.SDLK_KP_8:         key.CodeKeypad8,
		C.SDLK_KP_9:         key.CodeKeypad9,
		C.SDLK_KP_0:         key.CodeKeypad0,
		C.SDLK_KP_PERIOD:    key.CodeFullStop,
		//C.SDLK_APPLICATION,
		//C.SDLK_POWER,
		C.SDLK_KP_EQUALS: key.CodeKeypadEqualSign,
	}

	mouseButtonMap = map[C.Uint8]mouse.Button{
		C.SDL_BUTTON_LEFT:   mouse.ButtonLeft,
		C.SDL_BUTTON_MIDDLE: mouse.ButtonMiddle,
		C.SDL_BUTTON_RIGHT:  mouse.ButtonRight,
	}
)
