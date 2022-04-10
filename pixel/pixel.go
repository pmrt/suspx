package pixel

import (
	"image/color"
	"time"
)

type Pixel struct {
	Name  string
	Color *color.RGBA
	At    time.Time
}

type RawPixel struct {
	Name string
	At   time.Time
	Hex  string
}
