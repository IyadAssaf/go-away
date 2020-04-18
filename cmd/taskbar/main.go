package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
}

func logError(err error) {
	logrus.Error(err)
	errMenu.SetTitle(fmt.Sprintf("Error: %s", err))
	errMenu.Show()
}

func clearError() {
	errMenu.Hide()
}