package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"strings"
)

const getWindowTitlePowershellExpression = `gps |
? {$_.mainwindowtitle -like "*%s*"} |
? {$_.processname -like "%s"} |
select mainwindowtitle`

func (tc *WindowTitleChecker) check() (bool, error) {
	expr := fmt.Sprintf(getWindowTitlePowershellExpression, tc.PartialWindowTitle, tc.ProcessName)
	log.Debugf("executing powershell expression: %s", strings.Replace(expr, "\n", " ", -1))
	cmd := exec.Command("powershell", "-command", expr)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// Strip all leading/trailing whitespace, separate into lines, delete empty lines
	var outlines []string
	for _, s := range strings.Split(strings.Trim(string(out), " "), "\r\n") {
		if s != "" {
			outlines = append(outlines, s)
		}
	}

	log.Debugf("got %d lines from powershell", len(outlines))

	// Check that we have content in the array, if not the window doesn't exist
	if len(outlines) < 1 {
		return false, fmt.Errorf("window does not exist")
	}

	// Last line in array should now be our window title, so do the compare
	title := strings.Trim(outlines[len(outlines)-1], " ")
	if title != tc.LastTitle {
		tc.LastTitle = title
		return true, nil
	}

	return false, nil
}
