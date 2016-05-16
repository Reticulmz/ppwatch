// +build linux

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-ps"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var multipleSpaces = regexp.MustCompile(" +")

type unixWindow struct {
	Pid         int
	WindowTitle string
}

func (tc *WindowTitleChecker) check() (bool, error) {
	var osuPid int

	// First, get the pid of the process that matches our ProcessName
	processes, err := ps.Processes()
	if err != nil {
		return false, err
	}

	for _, process := range processes {
		if strings.Contains(process.Executable(), tc.ProcessName) {
			osuPid = process.Pid()
			log.Debugf("found osu! process with pid %d", osuPid)
			break
		}
	}

	if osuPid == 0 {
		return false, fmt.Errorf("unable to find osu! process")
	}

	// Get all open windows.
	// format: [windowid] [desktopid] [pid] [clienthostname] [title]
	cmd := exec.Command("wmctrl", "-lp")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	var window *unixWindow

	for _, s := range strings.Split(strings.Trim(string(out), " "), "\n") {
		if s != "" {
			// Turn multiple spaces into single spaces and split on spaces up
			// to 5 entries (last being window title)
			str := strings.SplitN(multipleSpaces.ReplaceAllString(s, " "), " ", 5)

			// PID is the 3rd entry
			pid, err := strconv.Atoi(str[2])
			if err != nil {
				return false, err
			}

			if pid == osuPid {
				window = &unixWindow{pid, str[len(str)-1]}
				break
			}
		}
	}

	if window == nil {
		return false, fmt.Errorf("can't find osu! window")
	}

	if window.WindowTitle != tc.LastTitle {
		tc.LastTitle = window.WindowTitle
		return true, nil
	}

	return false, nil
}
