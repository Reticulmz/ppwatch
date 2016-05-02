package main

import (
	log "github.com/Sirupsen/logrus"
)

type WindowTitleChecker struct {
	ProcessName        string
	PartialWindowTitle string

	LastTitle string
}

func NewWindowTitleChecker(processname, partialtitle string) *WindowTitleChecker {
	tc := &WindowTitleChecker{
		ProcessName:        processname,
		PartialWindowTitle: partialtitle,
		LastTitle:          "",
	}

	return tc
}

// Check return
func (tc *WindowTitleChecker) Check() (bool, error) {
	lasttitle := tc.LastTitle
	haschanged, err := tc.check()

	if haschanged {
		log.Debugf("title changed from '%s' to '%s'", lasttitle, tc.LastTitle)
	} else {
		log.Debugf("no title change")
	}

	return haschanged, err
}
