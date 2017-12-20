package mkv

import "time"

type Flag byte

const (
	FKeyFrame    Flag = 0x80
	FNoLace      Flag = 0x06
	FNoDuration  Flag = 0x08
	FDiscardable Flag = 0x01
)

type Block struct {
	Track    Vint
	TimeCode time.Duration
	Flag

	// wow, now we're getting to the good stuff
	// ... maybe later
}
