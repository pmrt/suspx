package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05.999 MST"
	// Separators for easy reading
	AppSection = "=============================================================="
	Section    = "-----------"
)

var (
	width          int
	height         int
	threshold      int
	margin         int
	cd             int
	out            string
	download       bool
	nonConsecutive bool
)

// TODO - Make `out` to default to params e.g.: cd5_m14010_t12.png
// TODO - Maybe we could experiment with detecting users with cooldown of 20
// minutes, running a simulation to detect users with consistent 20 minutes
// between pixels. Then run another simulation with only those users, this could
// improve the results of cd > 5
func main() {
	fmt.Println(AppSection)
	if download {
		downloadAll()
		fmt.Println("[+] All parts downloaded. decompress and re-run with the desired parameters")
		os.Exit(0)
	}

	fmt.Println("Ready to run simulation with following parameters:")
	fmt.Printf(
		"width:%d; height:%d; SusThreshold:%d; TimeMargin:%d; TimeCooldown:%d NonConsecutive:%t\n",
		width, height, threshold, margin, cd, nonConsecutive,
	)

	var filestats []*Filestats
	datasets := flag.Args()
	if len(datasets) == 0 {
		datasets, filestats = orderCSV()
	}
	fmt.Println(Section)
	fmt.Printf("Datasets (%d): %v\n", len(datasets), datasets)
	fmt.Println(Section)
	fmt.Printf(
		"First rows of datasets range from: (%v) to (%v).\n",
		filestats[0].FirstRowTimestamp, filestats[len(filestats)-1].FirstRowTimestamp,
	)

	fmt.Println(Section)
	for {
		var YesOrNo string
		fmt.Print("Execute simulation with these parameters? (y/n): ")
		fmt.Scanln(&YesOrNo)
		resp := strings.ToLower(YesOrNo)
		if resp == "y" || resp == "yes" {
			break
		} else if resp == "n" || resp == "no" {
			os.Exit(0)
		}
		fmt.Printf("\n")
	}

	fmt.Println(Section)
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
		Datasets:     datasets,
	}).
		Setup().
		Run().
		ExportPNG(out)
	fmt.Printf("canvas exported to %s\n", out)
}

func init() {
	flag.IntVar(&width, "w", 2000, "Pixel width of the canvas")
	flag.IntVar(&height, "h", 2000, "Pixel height of the canvas")
	flag.IntVar(&threshold, "threshold", 12, "Suspicious threshold. Above this threshold of consecutive pixels, the following consecutive suspicious pixels will be drawn")
	flag.IntVar(&margin, "margin", 14010, "If the pixel is placed in the time span of the cooldown and this time margin (in milliseconds), the pixel will be considered suspicious")
	flag.IntVar(&margin, "m", 14010, "Same as -margin")
	flag.IntVar(&cd, "cd", 5, "Time (in minutes) considered for the cooldown. A cooldown of 20 minutes with threshold 1, for example, will also make suspicious users with cooldown of 5 minutes with 4 suspicious pixels in the bucket")
	flag.BoolVar(&download, "d", false, "If a -d argument is provided, the download tool will be executed")
	flag.BoolVar(&nonConsecutive, "nc", false, "If a -nc argument is provided, all suspicious pixels defined by other parameters are drawn on the canvas, including the non-consecutive ones")
	flag.StringVar(&out, "o", "res.png", "Resulting PNG filename")
	flag.Parse()
}
