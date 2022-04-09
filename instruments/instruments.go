package instruments

import "github.com/pmrt/suspx/pixel"

const NInstruments = 2

type Hashtable map[string]InstrumentBucket

type Instrument interface {
	Run(b InstrumentBucket, rawpx *pixel.RawPixel, ht *Hashtable) bool
	Bucket() InstrumentBucket
	Setup()
	Report(ht *Hashtable)
	ShouldExport() bool
}

type InstrumentBucket interface {
	String() string
}

func Setup() map[string]Instrument {
	return map[string]Instrument{
		BotInstrumentName:   new(BotInstrument),
		StatsInstrumentName: new(StatsInstrument),
	}
}
