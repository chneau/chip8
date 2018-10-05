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
		}
	case 0x1000:
		c.PC = op & 0x0fff
	case 0x2000:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = op & 0x0fff
	case 0x3000:
		if uint16(c.V[op&0x0f00>>8]) == op&0x00ff {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x4000:
		if uint16(c.V[op&0x0f00>>8]) != op&0x00ff {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x5000:
		if c.V[op&0x0f00>>8] == c.V[op&0x00f0>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x6000:
		c.V[op&0x0f00>>8] = uint8(op & 0x00ff)
		c.PC += 2
	case 0x7000:
		c.V[op&0x0f00>>8] += uint8(op & 0x00ff)
		c.PC += 2
	case 0x8000:
		switch op & 0x000f {
		case 0x0000:
			c.V[op&0x0f00>>8] = c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0001:
			c.V[op&0x0f00>>8] |= c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0002:
			c.V[op&0x0f00>>8] &= c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0003:
			c.V[op&0x0f00>>8] ^= c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0004:
			if c.V[op&0x00f0>>4] > (0xff - c.V[op&0x0f00>>8]) {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.V[op&0x0f00>>8] += c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0005:
			if c.V[op&0x00f0>>4] > c.V[op&0x0f00>>8] {
				c.V[0xf] = 0
			} else {
				c.V[0xf] = 1
			}
			c.V[op&0x0f00>>8] -= c.V[op&0x00f0>>4]
			c.PC += 2
		case 0x0006:
			c.V[0xf] = c.V[op&0x0f00>>8] & 0x1
			c.V[op&0x0f00>>8] >>= 1
			c.PC += 2
		case 0x0007:
			if c.V[op&0x00f0>>4] < c.V[op&0x0f00>>8] {
				c.V[0xf] = 0
			} else {
				c.V[0xf] = 1
			}
			c.V[op&0x0f00>>8] = c.V[op&0x00f0>>4] - c.V[op&0x0f00>>8]
			c.PC += 2
		case 0x000e:
			c.V[0xf] = c.V[op&0x0f00>>8] >> 7 //& 0x80
			c.V[op&0x0f00>>8] <<= 1
			c.PC += 2
		}
	case 0x9000:
		if c.V[op&0x0f00>>8] != c.V[op&0x00f0>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0xa000:
		c.I = op & 0x0fff
		c.PC += 2
	case 0xb000:
		c.PC = op&0x0fff + uint16(c.V[0])
	case 0xc000:
		c.V[op&0x0f00>>8] = uint8(uint16(rand.Int()%0xff) & (op & 0x00ff))
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
		switch op & 0x00ff {
		case 0x009e:
			if c.V[op&0x0f00>>8] == c.Key {
				c.PC += 4
			} else {
				c.PC += 2
			}
		case 0x00a1:
			if c.Key == c.V[op&0x0f00>>8] {
				c.PC += 2
			} else {
				c.PC += 4
			}
		}
		c.Key = 0
	case 0xf000:
		switch op & 0x00ff {
		case 0x0007:
			c.V[op&0x0f00>>8] = c.DT
			c.PC += 2
		case 0x000a:
			c.V[op&0x0f00>>8] = c.Key // value of keyboard
			c.Key = 0
			c.PC += 2
		case 0x0015:
			c.DT = c.V[op&0x0f00>>8]
			c.PC += 2
		case 0x0018:
			c.ST = c.V[op&0x0f00>>8]
			c.PC += 2
		case 0x001e:
			if (c.I + uint16(c.V[op&0x0f00>>4])) > 0xfff {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.I += uint16(c.V[op&0x0f00>>8])
			c.PC += 2
		case 0x0029:
			c.I = 0x05 * (op & 0x0f00 >> 8)
			c.PC += 2
		case 0x0033:
			c.Mem[c.I] = c.V[op&0x0f00>>8] / 100
			c.Mem[c.I+1] = (c.V[op&0x0f00>>8] / 10) % 10
			c.Mem[c.I+2] = (c.V[op&0x0f00>>8] % 100) % 10
			c.PC += 2
		case 0x0055:
			for i := uint16(0); i <= (op & 0x0f00 >> 8); i++ {
				c.Mem[i+i] = c.V[i]
			}
			c.I += (op & 0x0f00 >> 8) + 1
			c.PC += 2
		case 0x0065:
			for i := uint16(0); i <= (op & 0x0f00 >> 8); i++ {
				c.V[i] = c.Mem[c.I+i]
			}
			c.I += (op & 0x0f00 >> 8) + 1
			c.PC += 2
		}
	default:
		panic("Not implemented")
	}
}
