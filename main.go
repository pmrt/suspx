package main

import (
	"fmt"

	"github.com/pmrt/suspx/simulation"
	"github.com/pmrt/suspx/utils"
)

const (
	// Separators for easy reading
	AppSection = "=============================================================="
	Section    = "-----------"
)

// TODO - Make `out` to default to params e.g.: cd5_m14010_t12.png
// TODO - Maybe we could experiment with detecting users with cooldown of 20
// minutes, running a simulation to detect users with consistent 20 minutes
// between pixels. Then run another simulation with only those users, this could
// improve the results of cd > 5
func main() {
	fmt.Println(AppSection)
	datasets, filestats := utils.OrderCSV()
	fmt.Printf("Datasets (%d): %v\n", len(datasets), datasets)
	fmt.Printf(
		"First rows of datasets range from: (%v) to (%v)\n",
		filestats[0].FirstRowTimestamp, filestats[len(filestats)-1].FirstRowTimestamp,
	)

	fmt.Println(Section)
	simulation.New(simulation.SimulationOptions{
		Datasets: datasets,
	}).
		Setup().
		Run()
}
