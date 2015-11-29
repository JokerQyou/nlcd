package pcd8544

import (
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

const (
	WHITE = 0
	BLACK = 1

	LCDWIDTH  = 84
	LCDHEIGHT = 48

	DISPLAY_BLANK    = 0x0
	DISPLAY_NORMAL   = 0x4
	DISPLAY_ALLON    = 0x1
	DISPLAY_INVERTED = 0x5

	POWERDOWN           = 0x04
	ENTRYMODE           = 0x02
	EXTENDEDINSTRUCTION = 0x01

	FUNCTIONSET     = 0x20
	DISPLAY_CONTROL = 0x08
	SETYADDR        = 0x40
	SETXADDR        = 0x80

	SETTEMP = 0x04
	SETBIAS = 0x10
	SETVOP  = 0x80

	CLKCONST_1 = 8000
	CLKCONST_2 = 400

	LSBFIRST = 0
	MSBFIRST = 1

	LCD_CMD  = uint8(embd.Low)
	LCD_DATA = uint8(embd.High)
)

type PCD8544 struct {
	PIN_SCLK uint8
	PIN_DIN  uint8
	PIN_DC   uint8
	PIN_CS   uint8
	PIN_RST  uint8
	CONTRAST uint8

	line   uint8
	column uint8
	custom map[uint8][5]uint8
}

type LCD interface {
}

func pre_shift(v uint8) int {
	var (
		v2   int
		zero uint8 = 0
	)
	if v == zero {
		v2 = 1
	} else {
		v2 = 0
	}
	return v2
}

func shiftOut(dataPin uint8, clockPin uint8, bitOrder uint8, val uint8) {
	var (
		one   uint8 = 1
		seven uint8 = 7
	)

	for i := uint8(0); i < 8; i++ {
		if bitOrder == LSBFIRST {
			embd.DigitalWrite(dataPin, pre_shift(val&(one<<i)))
		} else {
			embd.DigitalWrite(dataPin, pre_shift(val&(one<<(seven-i))))
		}
		embd.DigitalWrite(clockPin, embd.High)
		embd.DigitalWrite(clockPin, embd.Low)
	}
}

func (p *PCD8544) send(type_ uint8, data uint8) {
	embd.DigitalWrite(p.PIN_DC, int(type_))
	embd.DigitalWrite(p.PIN_CS, embd.Low)
	shiftOut(p.PIN_DIN, p.PIN_SCLK, MSBFIRST, data)
	embd.DigitalWrite(p.PIN_CS, embd.High)
}

func (p *PCD8544) Close() {
	p.Clear()
	p.setPower(false)
}

func (p *PCD8544) Clear() {
	p.SetCursor(0, 0)
	for i := 0; i < LCDWIDTH*(LCDHEIGHT/8); i++ {
		p.send(LCD_DATA, 0x00)
	}
	p.SetCursor(0, 0)
}

func (p *PCD8544) setPower(on bool) {
	if on {
		p.send(LCD_CMD, 0x20)
	} else {
		p.send(LCD_CMD, 0x24)
	}
}

func (p *PCD8544) Display() {
	p.setPower(true)
}

func (p *PCD8544) NoDisplay() {
	p.setPower(false)
}

func (p *PCD8544) SetInverse(inverse bool) {
	if inverse {
		p.send(LCD_CMD, 0x0d)
	} else {
		p.send(LCD_CMD, 0x0c)
	}
}

func (p *PCD8544) SetContrast(level uint8) {
	if level > 90 {
		level = 90
	}
	p.send(LCD_CMD, 0x21) // extended instruction set control (H=1)
	p.send(LCD_CMD, 0x80|(level&0x7f))
	p.send(LCD_CMD, 0x20) // extended instruction set control (H=0)
}

func (p *PCD8544) SetCursor(column uint8, line uint8) {
	p.column = column % LCDWIDTH
	p.line = line % (LCDHEIGHT/9 + 1)
	p.send(LCD_CMD, 0x80|p.column)
	p.send(LCD_CMD, 0x40|p.line)
}

func (p *PCD8544) CreateChar() {
	//
}

func (p *PCD8544) init() {
	// All pin direction are output
	embd.SetDirection(p.PIN_SCLK, embd.Out)
	embd.SetDirection(p.PIN_DIN, embd.Out)
	embd.SetDirection(p.PIN_DC, embd.Out)
	embd.SetDirection(p.PIN_CS, embd.Out)
	embd.SetDirection(p.PIN_RST, embd.Out)

	// Reset controller state
	embd.DigitalWrite(p.PIN_RST, embd.High)
	embd.DigitalWrite(p.PIN_CS, embd.High)
	embd.DigitalWrite(p.PIN_RST, embd.Low)
	// delay(100)
	time.Sleep(100 * time.Millisecond)
	embd.DigitalWrite(p.PIN_RST, embd.High)

	// LCD parameters
	p.send(LCD_CMD, 0x21) // extended instruction set control (H=1)
	p.send(LCD_CMD, 0x13) // bias system (1:48)
	p.send(LCD_CMD, 0xc2) // default Vop (3.06 + 66 * 0.06 = 7v)

	p.Clear()

	// Activate LCD
	p.send(LCD_CMD, 0x08) // display blank
	p.send(LCD_CMD, 0x0c) // normal mode (0x0d = inverse mode)
	// delay(100)
	time.Sleep(100 * time.Millisecond)

	// Place the cursor at the origin
	p.SetCursor(0, 0)
}

func (p *PCD8544) Home() {
	p.SetCursor(0, p.line)
}

func (p *PCD8544) Write(char uint8) {
	// ASCII 7-bit only
	if char >= 0x80 {
		return
	}
	var (
		glyph [5]uint8
	)
	if char >= ' ' {
		// Regular ASCII characters are kept in flash to save RAM...
		// memcpy_P(pgm_buffer, &charset[chr - ' '], sizeof(pgm_buffer));
		glyph = CHARSET[char]
	} else {
		// Custom glyphs, on the other hand, are stored in RAM...
		if _, ok := p.custom[char]; ok {
			glyph = p.custom[char]
		} else {
			// Default to a space character if unset...
			p.custom[char] = p.custom[' ']
			// memcpy_P(pgm_buffer, &charset[0], sizeof(pgm_buffer));
			glyph = p.custom[char]
		}
	}

	for _, i := range glyph {
		p.send(LCD_DATA, i)
	}
	p.send(LCD_DATA, 0x00)
	p.column = (p.column + 6) % LCDWIDTH
	if p.column == 0 {
		p.line = (p.line + 1) % (LCDHEIGHT/9 + 1)
	}
}

func New(clk uint8, din uint8, dc uint8, cs uint8, rst uint8, contrast uint8) *PCD8544 {
	p := PCD8544{
		PIN_SCLK: clk,
		PIN_DIN:  din,
		PIN_DC:   dc,
		PIN_CS:   cs,
		PIN_RST:  rst,
		CONTRAST: contrast,
		line:     0,
		column:   0,
		custom:   map[uint8][5]uint8{},
	}
	p.init()
	return &p
}
