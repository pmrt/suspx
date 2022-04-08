package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"
)

type SimulationOptions struct {
	Width      int
	Height     int
	BucketSize string
	Datasets   []string
}

type Simulation struct {
	SimulationOptions
	canvas *Canvas
	ht     Hashtable
}

func (s *Simulation) Setup() {
	s.canvas = NewCanvas(s.Width, s.Height)
	s.canvas.FillEmpty()
}

func (s *Simulation) IsEligible(rawpx RawPixel) bool {
	return true
}

func (s *Simulation) Run() {
	for _, dataset := range s.Datasets {
		f, err := os.Open(dataset)
		if err != nil {
			panic(err)
		}
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
				panic(err)
			}
			t, err := time.Parse("2006-01-02 15:04:05.999999999 MST", ts)
			if err != nil {
				panic(err)
			}

			raw := RawPixel{
				Name: name,
				At:   t,
				Hex:  hex,
			}
			if s.IsEligible(raw) {
				s.canvas.Set(x, y, raw)
			}
		}
	}
}

func NewSimulation(opts SimulationOptions) *Simulation {
	return &Simulation{
		SimulationOptions: opts,
		ht:                make(map[string]*PixelBucket),
	}
}
