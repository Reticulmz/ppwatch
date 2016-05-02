package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlayInfoString(t *testing.T) {
	playinfo := &PlayInfo{
		Time:       time.Now(),
		Beatmap:    "Test - Test",
		BeatmapID:  0,
		Difficulty: "Test",
		Rank:       "SH",
		Score:      1000000,
		MaxCombo:   1000,
		Perfect:    true,
		GainedPP:   1.01,
	}

	expected := "Test - Test [Test] | SH 1000x PERFECT (1000000) | +1.01 PP"
	playstr := playinfo.String()

	assert.Equal(t, expected, playstr)
}
