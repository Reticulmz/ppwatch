package main

import (
	"fmt"
	"strings"
	"time"
)

type PlayInfo struct {
	Time time.Time `json:"time"`

	Beatmap    string `json:"beatmap"`
	BeatmapID  int    `json:"beatmap_id"`
	Difficulty string `json:"difficulty"`

	Rank     string `json:"rank"`
	Score    int    `json:"score"`
	MaxCombo int    `json:"maxcombo"`
	Perfect  bool   `json:"perfect"`

	GainedPP float32 `json:"gainedpp"`
}

func (play *PlayInfo) String() string {
	// Display + on positive and - on negative
	gainedpp := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", play.GainedPP), "0"), ".")
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
