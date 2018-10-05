package chip8

import (
	"context"
	"io/ioutil"
	"time"
)

// the const
const (
	TimerFreq    = 240
	ScreenWidth  = 64
	ScreenHeight = 32
)

// Chip8 ...
type Chip8 struct {
	Mem          [0xfff]uint8 // 4k mem
	V            [16]uint8    // registers
	Stack        [16]uint16
	Gfx          [2048]uint8 // 32 * 64 = screen dims
	Key          byte
	Draw         bool
	SP, DT, ST   uint8
	I, PC        uint16
	DrawCallback func()
	BipCallback  func()
}

// LoadFromPath ...
func (c *Chip8) LoadFromPath(rom string) error {
	f, err := ioutil.ReadFile(rom)
	if err != nil {
		return err
	}
	c.Load(f)
	return nil
}

// Run ...
func (c *Chip8) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / TimerFreq)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.emulateCycle()
			if c.Draw {
				c.DrawCallback()
				c.Draw = false
			}
			c.updateTimers()
		}
	}
}

// Run ...
func (c *Chip8) updateTimers() {
	if c.ST > 0 {
		c.ST--
	}
	if c.ST == 1 {
		c.BipCallback()
	}
	if c.DT > 0 {
		c.DT--
	}
}

// Load a rom ...
func (c *Chip8) Load(rom []byte) {
	copy(c.Mem[0x200:], rom)
}

// New ...
func New() *Chip8 {
	c := &Chip8{}
	c.PC = 0x200
	c.DrawCallback = func() {}
	c.BipCallback = func() {}
	copy(c.Mem[:], []uint8{240, 144, 144, 144, 240, 32, 96, 32, 32, 112, 240, 16, 240, 128, 240, 240, 16, 240, 16, 240, 144, 144, 240, 16, 16, 240, 128, 240, 16, 240, 240, 128, 240, 144, 240, 240, 16, 32, 64, 64, 240, 144, 240, 144, 240, 240, 144, 240, 16, 240, 240, 144, 240, 144, 144, 224, 144, 224, 144, 224, 240, 128, 128, 128, 240, 224, 144, 144, 144, 224, 240, 128, 240, 128, 240, 240, 128, 240, 128, 128}) // FONT
	return c
}
