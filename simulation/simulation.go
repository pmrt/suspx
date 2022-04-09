package simulation

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pmrt/suspx/canvas"
	"github.com/pmrt/suspx/global"
	"github.com/pmrt/suspx/instruments"
	"github.com/pmrt/suspx/pixel"
	"github.com/pmrt/suspx/utils"
)

var (
	width    int
	height   int
	out      string
	download bool
)

type SimulationOptions struct {
	Datasets []string
}

type Simulation struct {
	SimulationOptions
	inst   instruments.Instrument
	canvas *canvas.Canvas
	ht     instruments.Hashtable
}

func (s *Simulation) Setup() *Simulation {
	insts := instruments.Setup()
	s.inst = askInstrument(insts)
	s.inst.Setup()
	flag.Parse()

	if download {
		utils.DownloadAll()
		fmt.Println("[+] All parts downloaded. decompress and re-run with the desired parameters")
		os.Exit(0)
	}

	s.canvas = canvas.New(width, height)
	s.canvas.FillEmpty()
	return s
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
				x, y, err := canvas.ParseCoord(pos)
				if err != nil {
					if errors.Is(err, canvas.ErrModerationTool) {
						// skip pixels by moderation rect tool
						continue
					}
					panic(err)
				}

				t, err := time.Parse(global.TimeLayout, ts)
				if err != nil {
					panic(err)
				}
				if s.ht[name] == nil {
					// ensure the bucket is always initialized
					s.ht[name] = s.inst.Bucket()
				}
				rawpx := &pixel.RawPixel{
					Name: name,
					At:   t,
					Hex:  hex,
				}

				if s.inst.Run(s.ht[name], rawpx, &s.ht) {
					s.canvas.Set(x, y, rawpx)
				}
			}
		}()
	}
	fmt.Printf("\n")
	s.inst.Report(&s.ht)
	if s.inst.ShouldExport() {
		s.ExportPNG()
	}
	return s
}

func (s *Simulation) ExportPNG() *Simulation {
	s.canvas.PNG(out)
	fmt.Printf("[+] canvas exported to %s\n", out)
	return s
}

func New(opts SimulationOptions) *Simulation {
	return &Simulation{
		SimulationOptions: opts,
		ht:                make(map[string]instruments.InstrumentBucket),
	}
}

func askInstrument(insts map[string]instruments.Instrument) instruments.Instrument {
	strs := make([]string, 0, len(insts))
	for key := range insts {
		strs = append(strs, key)
	}

	var resp string
	fmt.Printf("Instruments available: %s\n", strings.Join(strs, " | "))
	for {
		fmt.Print("Select an instrument to use for this simulation: ")
		fmt.Scanln(&resp)
		if inst := insts[resp]; inst != nil {
			return inst
		}
	}
}

func init() {
	flag.IntVar(&width, "w", 2000, "Pixel width of the canvas")
	flag.IntVar(&height, "h", 2000, "Pixel height of the canvas")
	flag.BoolVar(&download, "d", false, "If a -d argument is provided, the download tool will be executed")
	flag.StringVar(&out, "o", "res.png", "Resulting PNG filename")
}
