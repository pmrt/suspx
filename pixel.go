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
	cap    int
	bucket []*Pixel
}

func (b *PixelBucket) Append(px *Pixel) {
	b.bucket = append(b.bucket, px)
}

func (b *PixelBucket) isFull() bool {
	return len(b.bucket) >= b.cap
}

func NewPixelBucket(n int) *PixelBucket {
	b := &PixelBucket{
		cap:    n,
		bucket: make([]*Pixel, 0, n),
	}
	return b
}
