package main

const DefaultConfig = `---
username: EDITME
apikey: EDITME
gamemodes:
 - osu

wait_time: 2500ms
process: osu!
partialtitle: osu!

last_time: 1970-01-01 00:00:00
last_pp:
  '0': 0
  '1': 0
  '2': 0
  '3': 0
`

type Config struct {
	UserName  string   `json:"username"`
	APIKey    string   `json:"apikey"`
	GameModes []string `json:"gamemodes"`

	WaitTime           string `json:"wait_time"`
	ProcessName        string `json:"process"`
	PartialWindowTitle string `json:"partialtitle"`

	LastTime string             `json:"last_time"`
	LastPP   map[string]float32 `json:"last_pp"`
}
