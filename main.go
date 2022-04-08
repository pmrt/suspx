package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"
)

// 2022-04-04 21:13:48.993 UTC,0cSPtex/zsj8oFNjQM0V+sKmD+DB/aK3pJb4gRbRodHVPPT6wiF+d5/ayYxBQdbYDVUb5YazL8oyEgj11Ev4Ng==,#FFFFFF,"134,1634"
// 2022-04-04 21:13:48.993 UTC,o/9SjLe/n0CbJ7LedDysW5VyQElHmtUQfBNOmfWOkWbtTJkAnPiShPBvDmMjRRJUpoFHKV4f258n0YBNAm086Q==,#D4D7D9,"135,1600"
// 2022-04-04 21:13:48.993 UTC,e4d82+uuPlqWH71c4e/qOKumk2TB53ue1H57QUUbaLLHY4TGFBLuQ1rv4KdLTrT7r/41cuOXN2bfnddHQTdmrQ==,#FFFFFF,"110,1558"
func main() {
	c := NewCanvas(2000, 2000)
	c.FillEmpty()

	f, err := os.Open("2022_place_canvas_history-000000000077.csv")
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
		c.Set(x, y, RawPixel{
			Name: name,
			At:   t,
			Hex:  hex,
		})
	}

	c.PNG("res.png")
}
