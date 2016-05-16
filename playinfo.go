package main

import (
	"fmt"
	"strings"
	"time"
)

type PlayInfo struct {
	Time time.Time `json:"time"`

	GameMode   string `json:"gamemode"`
	Beatmap    string `json:"beatmap"`
	BeatmapID  int    `json:"beatmap_id"`
	Difficulty string `json:"difficulty"`

	Rank     string `json:"rank"`
	Score    int    `json:"score"`
	MaxCombo int    `json:"maxcombo"`
	Perfect  bool   `json:"perfect"`

	TotalPP  float32 `json:"totalpp"`
	GainedPP float32 `json:"gainedpp"`
}

func (play *PlayInfo) String() string {
	totalpp := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", play.TotalPP), "0"), ".")

	// Display + on positive and - on negative
	gainedpp := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", play.GainedPP), "0"), ".")
	if play.GainedPP >= 0 {
		gainedpp = fmt.Sprintf("+%s", gainedpp)
	}

	perfect := ""
	if play.Perfect {
		perfect = " PERFECT"
	}

	var gamemode string
	switch play.GameMode {
	case "osu":
		gamemode = "osu!"
	case "taiko":
		gamemode = "osu!taiko"
	case "ctb":
		gamemode = "CatchTheBeat"
	case "mania":
		gamemode = "osu!mania"
	default:
		gamemode = "unknown gamemode"
	}

	return fmt.Sprintf(
		"%s [%s] (%s) | %s %dx%s (%d) | %s PP (%s)",
		play.Beatmap,
		play.Difficulty,
		gamemode,
		play.Rank,
		play.MaxCombo,
		perfect,
		play.Score,
		totalpp,
		gainedpp,
	)
}
