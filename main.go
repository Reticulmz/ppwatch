package main

import (
	"flag"
	"time"
	"path"
	"os"
	"os/user"
	"io/ioutil"

	"github.com/ghodss/yaml"
	log "github.com/Sirupsen/logrus"
)

const (
	VERSION = "v0.1.0"
)

func main() {
	var err error

	var debug = flag.Bool("debug", false, "enable debug logging")
	var nocolor = flag.Bool("nocolor", false, "disable color in output")
	var configpath = flag.String("config", "", "path to configuration")

	flag.Parse()

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(NewLogFormatter(*nocolor))
	log.Infof("ppwatch %s", VERSION)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	user, err := user.Current()
	if err != nil {
		log.Panicf("failed to get user information: %s", err)
	}

	configPath := *configpath
	if configPath == "" {
		configPath = path.Join(user.HomeDir, ".ppwatch.yml")
	}

	// check config file exists, write it if not
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		ioutil.WriteFile(configPath, []byte(DefaultConfig), 0600)

		log.Fatalf("Please edit '%s' with your username and osu! API key.", configPath)
	}

	var config Config
	configyml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read configuration: %s", err)
	}

	err = yaml.Unmarshal(configyml, &config)
	if err != nil {
		log.Panicf("failed to unmarshal configuration: %s", err)
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

			newplay, data, err := api.CheckForPlay()
			if err != nil {
				log.Errorf("error getting play data: %s", err)
				continue
			}

			if newplay {
				log.Infof("%s", data)

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