package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	width     int
	height    int
	threshold int
	margin    int
	cd        int
	out       string
)

func main() {
	fmt.Println("Running simulation with following parameters:")
	fmt.Printf(
		"width:%d; height:%d; SusThreshold:%d; TimeMargin:%d; TimeCooldown:%d\n",
		width, height, threshold, margin, cd,
	)
	fmt.Printf("Datasets: %v\n", flag.Args())

	NewSimulation(SimulationOptions{
		// Canvas size
		Width:  width,
		Height: height,
		// Above the sus threshold of consecutive sus pixels, the pixels will be
		// drawn on the canvas. If a non-suspicious pixel comes in the bucket will
		// be reset, but the old pixels will remain on the canvas until they are
		// overwritten. fg
		SusThreshold: threshold,
		TimeMargin:   time.Millisecond * time.Duration(margin),
		TimeCooldown: time.Minute * time.Duration(cd),
		Datasets:     flag.Args(),
	}).
		Setup().
		Run().
		ExportPNG(out)
	fmt.Printf("canvas exported to %s\n", out)
}

func init() {
	flag.IntVar(&width, "w", 2000, "Pixel width of the canvas")
	flag.IntVar(&height, "h", 2000, "Pixel height of the canvas")
	flag.IntVar(&threshold, "threshold", 3, "Suspicious threshold. Above this threshold of consecutive pixels, the following consecutive pixels will be drawn.")
	flag.IntVar(&margin, "margin", 500, "Time margin (ms) from the cooldown to consider a pixel suspicious")
	flag.IntVar(&cd, "cd", 5, "Time cooldown (min) considered for the cooldown")
	flag.StringVar(&out, "o", "res.png", "Resulting PNG filename (default: res.png)")
	flag.Parse()
}
