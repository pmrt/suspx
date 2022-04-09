package instruments

import (
	"flag"
	"fmt"
	"time"

	"github.com/pmrt/suspx/pixel"
)

var (
	threshold      int
	cd             time.Duration
	margin         time.Duration
	nonConsecutive bool
)

const BotInstrumentName = "bot"

type BotBucket struct {
	SusCount int
	LastPx   *pixel.RawPixel
}

func (b *BotBucket) String() string {
	return BotInstrumentName
}

type BotInstrument struct{}

func (b *BotInstrument) Run(bucket InstrumentBucket, rawpx *pixel.RawPixel, _ *Hashtable) (shouldDraw bool) {
	bkt, ok := bucket.(*BotBucket)
	if !ok {
		panic("expected type of InstrumentBucket to be *BotBucket")
	}

	// If the sus threshold is surpassed, the following consecutive
	// suspicious pixels are drawn on the canvas...
	if isSus(bkt.LastPx, rawpx) {
		if bkt.SusCount >= threshold {
			shouldDraw = true
		} else {
			bkt.SusCount++
		}
	} else if !nonConsecutive {
		// ...reset the bucket otherwise so we only draw consecutive pixels
		bkt.SusCount = 0
	}
	bkt.LastPx = rawpx
	return
}

func isSus(lastpx *pixel.RawPixel, rawpx *pixel.RawPixel) bool {
	if lastpx == nil {
		return false
	}
	delta := rawpx.At.Sub(lastpx.At).Milliseconds()
	return delta < (cd.Milliseconds() + margin.Milliseconds())
}

func (b *BotInstrument) Bucket() InstrumentBucket {
	return new(BotBucket)
}

// No reports needed
func (b *BotInstrument) Report(ht *Hashtable) {}

func (b *BotInstrument) ShouldExport() bool {
	return true
}

func (b *BotInstrument) Setup() {
	flag.IntVar(&threshold, "threshold", 12, "Suspicious threshold. Above this threshold of consecutive pixels, the following consecutive suspicious pixels will be drawn")

	var m int
	flag.IntVar(&m, "margin", 14010, "If the pixel is placed in the time span of the cooldown and this time margin (in milliseconds), the pixel will be considered suspicious")
	flag.IntVar(&m, "m", 14010, "Same as -margin")
	flag.BoolVar(&nonConsecutive, "nc", false, "If a -nc argument is provided, all suspicious pixels defined by other parameters are drawn on the canvas, including the non-consecutive ones")
	margin = time.Millisecond * time.Duration(m)

	var c int
	flag.IntVar(&c, "cd", 5, "Time (in minutes) considered for the cooldown. A cooldown of 20 minutes with threshold 1, for example, will also make suspicious users with cooldown of 5 minutes with 4 suspicious pixels in the bucket")
	cd = time.Minute * time.Duration(c)

	fmt.Printf(
		"Selected parameters: [threshold:%d] [time_margin:%d] [time_cooldown:%d]\n\n",
		threshold, margin, cd,
	)
}
