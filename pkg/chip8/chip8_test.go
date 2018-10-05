package chip8

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c.Mem[0] != 0xf0 {
		t.FailNow()
	}
	if c.Mem[77] != 0xf0 {
		t.FailNow()
	}
}
