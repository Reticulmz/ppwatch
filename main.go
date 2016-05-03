package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

const (
	VERSION = "v0.2.0"
)

const (
	UserNamePrompt = "What is your osu! username?"

	APIKeyPrompt = `What is your osu! API key?
	If you don't have an API key, go to https://osu.ppy.sh/p/api`

	GameModePrompt = `What game modes would you like to get ranks for?
	Available modes: osu, taiko, ctb, mania
	Separate the modes by a space (example: 'osu ctb')`
)

var config Config

func main() {
	var err error

	var debug = flag.Bool("debug", false, "enable debug logging")
	var nocolor = flag.Bool("nocolor", false, "disable color in output")
	var configpath = flag.String("config", "", "path to configuration")
	var jsonoutput = flag.Bool("json", false, "output as json")

	flag.Parse()

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(NewLogFormatter(*nocolor))

	if *jsonoutput {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("ppwatch %s", VERSION)

	user, err := user.Current()
	if err != nil {
		log.Panicf("failed to get user information: %s", err)
	}

	configPath := *configpath
	if configPath == "" {
		configPath = path.Join(user.HomeDir, ".ppwatch.yml")
	}

	// check config file exists, prompt for necessary information if not
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = yaml.Unmarshal([]byte(DefaultConfig), &config)
		if err != nil {
			log.Fatalf("failed to unmarshal default config: %s", err)
		}

		config.UserName, err = PromptForInput(UserNamePrompt, "")
		if err != nil {
			log.Fatalf("failed to prompt for input: %s", err)
		}

		config.APIKey, err = PromptForInput(APIKeyPrompt, "")
		if err != nil {
			log.Fatalf("failed to prompt for input: %s", err)
		}

		modestring, err := PromptForInput(GameModePrompt, "osu taiko ctb mania")
		if err != nil {
			log.Fatalf("failed to prompt for input: %s", err)
		}

		modes := strings.Split(strings.Trim(modestring, " \r\n"), " ")
		var c []string
		for _, str := range modes {
			c = append(c, strings.Trim(str, " "))
		}

		config.GameModes = c

		cfgout, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("error saving config: %s", err)
		}

		err = ioutil.WriteFile(configPath, cfgout, 0600)
		if err != nil {
			log.Fatalf("error saving config: %s", err)
		}
	} else {
		// if we're here, config file exists at start, so load it up

		configyml, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("failed to read configuration: %s", err)
		}

		err = yaml.Unmarshal(configyml, &config)
		if err != nil {
			log.Fatalf("failed to unmarshal configuration: %s", err)
		}
	}

	api := NewAPIChecker(config.UserName, config.APIKey)
	api.LastTime, _ = time.Parse("2006-01-02 15:04:05", config.LastTime)
	api.LastPP = config.LastPP

	log.Debugf("process: \"%s\"", config.ProcessName)
	log.Debugf("partial window title: \"%s\"", config.PartialWindowTitle)

	tc := NewWindowTitleChecker(config.ProcessName, config.PartialWindowTitle)

	waitforwindow := true
	log.Info("waiting for osu! window...")

	for waitforwindow {
		_, err := tc.Check()

		if err != nil {
			if err.Error() != "window does not exist" {
				log.Fatalf("error checking for window: %s", err)
			}
		} else {
			waitforwindow = false
		}

		time.Sleep(1)
	}

	log.Infof("found osu! window: %s", tc.LastTitle)

	for {
		haschanged, err := tc.Check()

		if err != nil {
			if err.Error() == "window does not exist" {
				if !waitforwindow {
					log.Infof("osu! window has closed, waiting for it to return...")
					waitforwindow = true
				}
			} else {
				log.Fatalf("error checking for window change: %s", err)
			}
		} else {
			if waitforwindow {
				log.Infof("osu! window has returned: %s", tc.LastTitle)
				waitforwindow = false
			}
		}

		if haschanged {
			duration, err := time.ParseDuration(config.WaitTime)
			if err != nil {
				log.Errorf("error parsing wait_time from config, defaulting to 2500ms")
				duration, _ = time.ParseDuration("2500ms")
			}

			time.Sleep(duration)

			newplay, data, err := api.CheckForPlay(config.GameModes)
			if err != nil {
				log.Errorf("error getting play data: %s", err)
				continue
			}

			if newplay {
				if *jsonoutput {
					out, err := json.Marshal(data)
					if err != nil {
						log.Errorf("error marshalling json: %s", err)
						continue
					}

					fmt.Println(string(out))
				} else {
					log.Infof("%s", data)
				}

				// Write the new PP and date values to the config, and save
				config.LastTime = api.LastTime.Format("2006-01-02 15:04:05")
				config.LastPP = api.LastPP

				cfgout, err := yaml.Marshal(config)
				if err != nil {
					log.Errorf("error saving config: %s", err)
					continue
				}

				err = ioutil.WriteFile(configPath, cfgout, 0600)
				if err != nil {
					log.Errorf("error saving config: %s", err)
					continue
				}
			}
		}

		time.Sleep(1)
	}
}
