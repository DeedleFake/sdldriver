package main

import (
	driver "github.com/DeedleFake/sdldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"image"
	"image/color"
	"log"
)

func main() {
	driver.Main(func(s screen.Screen) {
		win, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  640,
			Height: 480,
		})
		if err != nil {
			log.Fatalf("Failed to create window: %v", err)
		}
		defer win.Release()

		win.Fill(image.ZR, color.Black, screen.Src)
		win.Fill(image.Rect(10, 10, 110, 60), color.NRGBA{255, 0, 255, 255}, screen.Src)
		win.Publish()

		for {
			ev := win.NextEvent()
			log.Printf("Event: %#v", win.NextEvent())

			switch ev := ev.(type) {
			case lifecycle.Event:
				if ev.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}
