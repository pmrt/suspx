package canvas

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/pmrt/suspx/pixel"
)

type Canvas struct {
	Width, Height int
	Colors        [24]color.RGBA
	Grid          [][]*pixel.Pixel
}

func (c *Canvas) FillEmpty() {
	c.Grid = make([][]*pixel.Pixel, c.Height)
	for x := 0; x < c.Height; x++ {
		c.Grid[x] = make([]*pixel.Pixel, c.Height)
		for y := 0; y < c.Width; y++ {
			c.Grid[x][y] = &pixel.Pixel{
				Color: &color.RGBA{255, 255, 255, 0xff},
			}
		}
	}
}

func (c *Canvas) PNG(name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			img.Set(x, y, c.Grid[x][y].Color)
		}
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func (c *Canvas) Set(x, y int, rawpx *pixel.RawPixel) {
	px := c.Grid[x][y]
	px.Name, px.At = rawpx.Name, rawpx.At
	ParseHex(px, rawpx.Hex)
}

var ErrModerationTool = errors.New("coordinates from the moderation rect tool")

func ParseCoord(pos string) (x int, y int, err error) {
	pos = strings.ReplaceAll(pos, "\"", "")
	parts := strings.Split(pos, ",")

	if len(parts) == 4 {
		// moderation rect tool
		err = ErrModerationTool
		return

	}

	if x, err = strconv.Atoi(parts[0]); err != nil {
		return
	}
	if y, err = strconv.Atoi(parts[1]); err != nil {
		return
	}
	return
}

// Thanks to @icza for the OG hex2rgba function.
func ParseHex(dst *pixel.Pixel, hex string) {
	c := dst.Color
	c.A = 0xff

	if hex[0] != '#' {
		panic("sethex: invalid format: missing #")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		panic("parseHex: invalid format")
	}

	switch len(hex) {
	case 7:
		c.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		c.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		c.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		c.R = hexToByte(hex[1]) * 17
		c.G = hexToByte(hex[2]) * 17
		c.B = hexToByte(hex[3]) * 17
	default:
		panic("parseHex: invalid format")
	}

}

func New(w, h int) *Canvas {
	c := &Canvas{
		Width:  w,
		Height: h,
		Colors: [24]color.RGBA{
			{190, 0, 57, 0xff},
			{255, 69, 0, 0xff},
			{255, 168, 0, 0xff},
			{255, 213, 52, 0xff},
			{0, 159, 97, 0xff},
			{0, 204, 120, 0xff},
			{126, 237, 86, 0xff},
			{0, 117, 111, 0xff},
			{0, 158, 170, 0xff},
			{36, 80, 164, 0xff},
			{54, 144, 234, 0xff},
			{81, 233, 244, 0xff},
			{73, 58, 193, 0xff},
			{105, 91, 255, 0xff},
			{129, 30, 159, 0xff},
			{180, 74, 192, 0xff},
			{255, 56, 129, 0xff},
			{255, 153, 170, 0xff},
			{156, 105, 38, 0xff},
			{109, 72, 47, 0xff},
			{0, 0, 0, 0xff},
			{137, 141, 144, 0xff},
			{212, 215, 217, 0xff},
			{255, 255, 255, 0xff},
		},
	}
	c.FillEmpty()
	return c
}
