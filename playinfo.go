package main

import (
	"time"
	"fmt"
	"strings"
)

type PlayInfo struct {
	Time time.Time

	Beatmap string
	BeatmapID int
	Difficulty string

	Rank string
	Score int
	MaxCombo int
	Perfect bool

	GainedPP float32
}

func (play *PlayInfo) String() string {
	// Display + on positive and - on negative
	gainedpp := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", play.GainedPP), "0"), ".")
	if play.GainedPP > 0 {
		gainedpp = fmt.Sprintf("+%s", gainedpp)
	}

	gainedpp = fmt.Sprintf("%s PP", gainedpp)

	perfect := ""
	if play.Perfect {
		perfect = " PERFECT"
	}

	return fmt.Sprintf(
		"%s [%s] | %s %dx%s (%d) | %s",
		play.Beatmap,
		play.Difficulty,
		play.Rank,
		play.MaxCombo,
		perfect,
		play.Score,
		gainedpp,
	)
}