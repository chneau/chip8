package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/chneau/chip8/pkg/chip8"
	"github.com/nsf/termbox-go"
)

func init() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()
}

var keyMap = map[rune]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
	'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := chip8.New()
	err := c.LoadFromPath("roms/TICTAC")
	if err != nil {
		panic(err)
	}
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Key == termbox.KeyEsc {
				cancel()
				return
			}
			if v, exist := keyMap[ev.Ch]; exist == true {
				c.Key = v
			}
		}
	}()

	c.DrawCallback = func() {
		for k := 0; k < chip8.ScreenHeight; k++ {
			for l := 0; l < chip8.ScreenWidth; l++ {
				r := rune(' ')
				if c.Gfx[(l+(k*chip8.ScreenWidth))] == 1 {
					r = '█'
				}
				termbox.SetCell(l, k, r, termbox.ColorGreen, termbox.ColorBlack)
			}
		}
		termbox.Flush()
	}
	c.Run(ctx)
}
