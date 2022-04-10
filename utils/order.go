package utils

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
		if askDownload("No CSV datasets found, do you want to initialize the download tool?") {
			DownloadAll()
			os.Exit(0)
		}
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

func askDownload(msg string) bool {
	var dl string
	for {
		fmt.Printf("%s (y/n): ", msg)
		fmt.Scanln(&dl)
		switch strings.ToLower(dl) {
		case "y", "yes":
			return true
		case "n", "no":
			os.Exit(0)
		}
	}
}
