package chip8

import (
	"math/rand"
)

func (c *Chip8) emulateCycle() {
	op := uint16(c.Mem[c.PC])<<8 | uint16(c.Mem[c.PC+1])
	switch op & 0xf000 {
	case 0x0000:
		switch op & 0x0fff {
		case 0x00e0:
			c.Gfx = [2048]uint8{}
			c.Draw = true
			c.PC += 2
		case 0x00ee:
			c.SP--
			c.PC = c.Stack[c.SP]
			c.PC += 2
		default:
			panic("Not implemented")
		}
	case 0x1000:
		c.PC = op & 0x0fff
	case 0x2000:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = op & 0x0fff
	case 0x3000:
		x := (op & 0x0F00) >> 8
		c.PC += 2
		if c.V[x] == byte(op) {
			c.PC += 2
		}
	case 0x4000:
		x := (op & 0x0F00) >> 8
		c.PC += 2
		if c.V[x] != byte(op) {
			c.PC += 2
		}
	case 0x5000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		c.PC += 2
		if c.V[x] == c.V[y] {
			c.PC += 2
		}
	case 0x6000:
		x := (op & 0x0F00) >> 8
		kk := byte(op)
		c.V[x] = kk
		c.PC += 2
	case 0x7000:
		x := (op & 0x0F00) >> 8
		kk := byte(op)
		c.V[x] = c.V[x] + kk
		c.PC += 2
	case 0x8000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		switch op & 0x000f {
		case 0x0000:
			c.V[x] = c.V[y]
			c.PC += 2
		case 0x0001:
			c.V[x] |= c.V[y]
			c.PC += 2
		case 0x0002:
			c.V[x] &= c.V[y]
			c.PC += 2
		case 0x0003:
			c.V[x] ^= c.V[y]
			c.PC += 2
		case 0x0004:
			r := uint16(c.V[x]) + uint16(c.V[y])
			var cf byte
			if r > 0xFF {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = byte(r)
			c.PC += 2
		case 0x0005:
			var cf byte
			if c.V[x] > c.V[y] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] - c.V[y]
			c.PC += 2
		case 0x0006:
			var cf byte
			if (c.V[x] & 0x01) == 0x01 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] / 2
			c.PC += 2
		case 0x0007:
			var cf byte
			if c.V[y] > c.V[x] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[y] - c.V[x]
			c.PC += 2
		case 0x000e:
			var cf byte
			if (c.V[x] & 0x80) == 0x80 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] * 2
			c.PC += 2
		default:
			panic("Not implemented")
		}
	case 0x9000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		c.PC += 2
		if c.V[x] != c.V[y] {
			c.PC += 2
		}
	case 0xa000:
		c.I = op & 0x0FFF
		c.PC += 2
	case 0xb000:
		c.PC = op&0x0FFF + uint16(c.V[0])
	case 0xc000:
		x := (op & 0x0F00) >> 8
		kk := byte(op)
		c.V[x] = kk + uint8(uint16(rand.Int()%0xff))
		c.PC += 2
	case 0xd000:
		height := uint8(op & 0x000f)
		x := c.V[op&0x0f00>>8]
		y := c.V[op&0x00f0>>4]
		c.V[0xf] = 0

		for yline := uint16(0); yline < uint16(height); yline++ {
			pixel := c.Mem[c.I+uint16(yline)]
			for xline := uint16(0); xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.Gfx[uint16(x)+xline+((uint16(y)+yline)*ScreenWidth)] == 1 {
						c.V[0xf] = 1
					}
					c.Gfx[uint16(x)+xline+((uint16(y)+yline)*ScreenWidth)] ^= 1
				}
			}
		}
		c.Draw = true
		c.PC += 2
	case 0xe000:
		x := (op & 0x0F00) >> 8
		switch op & 0x00ff {
		case 0x009e:
			c.PC += 2
			if c.V[x] == c.Key {
				c.PC += 2
			}
		case 0x00a1:
			c.PC += 2
			if c.V[x] != c.Key {
				c.PC += 2
			}
		default:
			panic("Not implemented")
		}
		c.Key = 0
	case 0xf000:
		x := (op & 0x0F00) >> 8
		switch op & 0x00ff {
		case 0x0007:
			c.V[x] = c.DT
			c.PC += 2
		case 0x000a:
			c.V[x] = c.Key
			c.Key = 0
			c.PC += 2
		case 0x0015:
			c.DT = c.V[x]
			c.PC += 2
		case 0x0018:
			c.ST = c.V[x]
			c.PC += 2
		case 0x001e:
			c.I = c.I + uint16(c.V[x])
			c.PC += 2
		case 0x0029:
			c.I = uint16(c.V[x]) * uint16(0x05)
			c.PC += 2
		case 0x0033:
			c.Mem[c.I] = c.V[x] / 100
			c.Mem[c.I+1] = (c.V[x] / 10) % 10
			c.Mem[c.I+2] = (c.V[x] % 100) % 10
			c.PC += 2
		case 0x0055:
			for i := 0; uint16(i) <= x; i++ {
				c.Mem[c.I+uint16(i)] = c.V[i]
			}
			c.PC += 2
		case 0x0065:
			for i := 0; byte(i) <= byte(x); i++ {
				c.V[uint16(i)] = c.Mem[c.I+uint16(i)]
			}
			c.PC += 2
		default:
			panic("Not implemented")
		}
	default:
		panic("Not implemented")
	}
}
