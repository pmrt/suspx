package main

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func findCSV(path string) []string {
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

func orderCSV() []string {
	datasets := findCSV(".")

	// hashtable of first rows to files
	// rowtable := make(map[string]string, len(datasets))
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
		t, err := time.Parse("2006-01-02 15:04:05.999 MST", peek[0])
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
	return all
}
