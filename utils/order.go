package utils

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pmrt/suspx/global"
)

func FindCSV(path string) []string {
	all := make([]string, 0, 10)
	filepath.WalkDir(path, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == ".csv" {
			all = append(all, s)
		}
		return nil
	})
	return all
}

type Filestats struct {
	Name              string
	FirstRowTimestamp time.Time
}

func OrderCSV() ([]string, []*Filestats) {
	datasets := FindCSV(".")

	if len(datasets) == 0 {
		panic("No CSV datasets found in current path")
	}

	// slice with file details for later sorting
	filestats := make([]*Filestats, 0, len(datasets))

	for _, dataset := range datasets {
		f, err := os.Open(dataset)
		if err != nil {
			panic(err)
		}

		r := csv.NewReader(f)
		// skip header
		if _, err := r.Read(); err != nil {
			panic(err)
		}
		peek, err := r.Read()
		if err != nil {
			panic(err)
		}
		t, err := time.Parse(global.TimeLayout, peek[0])
		if err != nil {
			panic(err)
		}
		filestats = append(filestats, &Filestats{
			Name:              dataset,
			FirstRowTimestamp: t,
		})
	}

	sort.Slice(filestats, func(a, b int) bool {
		return filestats[a].FirstRowTimestamp.Before(filestats[b].FirstRowTimestamp)
	})

	all := make([]string, 0, len(filestats))
	for _, stat := range filestats {
		all = append(all, stat.Name)
	}
	return all, filestats
}
