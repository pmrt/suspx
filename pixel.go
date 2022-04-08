package main

import (
	"image/color"
	"time"
)

type Hashtable map[string]*PixelBucket

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

type PixelBucket struct {
	threshold int
	// Number of suspicious pixels in the bucket
	sus int
	// Last pixel by the same user
	LastPx *RawPixel
}

func (b *PixelBucket) Add(n int) {
	b.sus = b.sus + n
}

func (b *PixelBucket) Reset() {
	b.sus = 0
}

func (b *PixelBucket) isFull() bool {
	return b.sus >= b.threshold
}

func NewPixelBucket(n int) *PixelBucket {
	b := &PixelBucket{
		threshold: n,
	}
	return b
}
