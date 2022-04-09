package instruments

import (
	"fmt"

	"github.com/pmrt/suspx/pixel"
)

const StatsInstrumentName = "stats"

type StatsBucket struct {
	AllCount int
}

func (b *StatsBucket) String() string {
	return StatsInstrumentName
}

type StatsInstrument struct {
	maxPlaced int
}

func (s *StatsInstrument) Run(bucket InstrumentBucket, rawpx *pixel.RawPixel, _ *Hashtable) (shouldDraw bool) {
	bkt, ok := bucket.(*StatsBucket)
	if !ok {
		panic("expected type of InstrumentBucket to be *StatsBucket")
	}

	bkt.AllCount++
	if nplaced := bkt.AllCount; nplaced > s.maxPlaced {
		s.maxPlaced = nplaced
	}
	return
}

func (s *StatsInstrument) Bucket() InstrumentBucket {
	return new(StatsBucket)
}

// No setup needed
func (s *StatsInstrument) Setup() {}

func (b *StatsInstrument) ShouldExport() bool {
	return false
}

func (s *StatsInstrument) Report(ht *Hashtable) {
	fmt.Printf("\n* The user with most pixels placed have placed %d pixels\n", s.maxPlaced)
	fmt.Printf("* A total of %d users placed a pixel\n", len(*ht))
}
