package main

const DefaultConfig = `---
username: EDITME
apikey: EDITME

wait_time: 2500ms
process: osu!
partialtitle: osu!

last_time: 1970-01-01 00:00:00
last_pp: 0`

type Config struct {
	UserName string `json:"username"`
	APIKey string `json:"apikey"`

	WaitTime string `json:"wait_time"`
	ProcessName string `json:"process"`
	PartialWindowTitle string `json:"partialtitle"`

	LastTime string `json:"last_time"`
	LastPP float32 `json:"last_pp"`
}