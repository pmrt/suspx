package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type SimulationOptions struct {
	Width        int
	Height       int
	SusThreshold int
	Datasets     []string
	TimeMargin   time.Duration
	TimeCooldown time.Duration
}

type Simulation struct {
	SimulationOptions
	canvas    *Canvas
	ht        Hashtable
	maxPlaced int
}

func (s *Simulation) Setup() *Simulation {
	s.canvas = NewCanvas(s.Width, s.Height)
	s.canvas.FillEmpty()
	return s
}

// Suspicious determines if the given pixel is worth storing in the
// Suspicious bucket.
func (s *Simulation) Suspicious(rawpx *RawPixel) bool {
	bucket := s.ht[rawpx.Name]
	last := bucket.LastPx

	if last == nil {
		// first pixel of the user
		return false
	}

	delta := rawpx.At.Sub(last.At).Milliseconds()
	return delta < (s.TimeCooldown.Milliseconds() + s.TimeMargin.Milliseconds())
}

func (s *Simulation) Run() *Simulation {
	l := len(s.Datasets)
	fmt.Printf("Simulating pixels... [0/%d]", l)
	for i, dataset := range s.Datasets {
		fmt.Printf("\rSimulating pixels... [%d/%d]", i+1, l)
		f, err := os.Open(dataset)
		if err != nil {
			panic(err)
		}

		// wrap in a func so we can defer Close() and avoid future memory leaks
		// because we forgot a Close()
		func() {
			defer f.Close()
			r := csv.NewReader(f)
			// Skip header
			if _, err := r.Read(); err != nil {
				panic(err)
			}
			for {
				row, err := r.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					panic(err)
				}

				ts, name, hex, pos := row[0], row[1], row[2], row[3]
				x, y, err := parseCoord(pos)
				if err != nil {
					if errors.Is(err, ErrModerationTool) {
						// skip pixels by moderation rect tool
						continue
					}
					panic(err)
				}

				t, err := time.Parse(TimeLayout, ts)
				if err != nil {
					panic(err)
				}
				if s.ht[name] == nil {
					// ensure the bucket is always initialized
					s.ht[name] = NewPixelBucket(s.SusThreshold)
				}

				raw := &RawPixel{
					Name: name,
					At:   t,
					Hex:  hex,
				}
				bucket := s.ht[name]
				// If the sus threshold is surpassed, the following consecutive
				// suspicious pixels are drawn on the canvas...
				if s.Suspicious(raw) {
					if bucket.isFull() {
						s.canvas.Set(x, y, raw)
					} else {
						bucket.AddSus(1)
					}
				} else if !nonConsecutive {
					// ...reset the bucket otherwise so we only draw consecutive pixels
					bucket.ResetSus()
				}
				bucket.LastPx = raw

				bucket.AddAll(1)
				nplaced := bucket.All()
				if nplaced > s.maxPlaced {
					s.maxPlaced = nplaced
				}
			}
		}()
	}
	fmt.Printf("\n")
	fmt.Printf("* A total of %d users placed a pixel\n", len(s.ht))
	fmt.Printf("* The user with most pixels placed have placed %d pixels\n", s.maxPlaced)
	return s
}

func (s *Simulation) ExportPNG(name string) *Simulation {
	s.canvas.PNG(name)
	return s
}

func NewSimulation(opts SimulationOptions) *Simulation {
	return &Simulation{
		SimulationOptions: opts,
		ht:                make(map[string]*PixelBucket),
	}
}
