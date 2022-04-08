package main

import (
	"time"
)

func main() {
	NewSimulation(SimulationOptions{
		// Canvas size
		Width:  2000,
		Height: 2000,
		// Above the sus threshold of consecutive sus pixels, the pixels will be
		// drawn on the canvas. If a non-suspicious pixel comes in the bucket will
		// be reset, but the old pixels will remain on the canvas until they are
		// overwritten.
		SusThreshold: 3,
		TimeMargin:   time.Millisecond * 1000,
		TimeCooldown: time.Minute * 5,
		Datasets: []string{
			"2022_place_canvas_history-000000000077.csv",
		},
	}).
		Setup().
		Run().
		ExportPNG("res.png")
}
