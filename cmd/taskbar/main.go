package main

import "github.com/getlantern/systray"

func main() {
	systray.RunWithAppWindow("Preferences", 1500, 1500, onReady, onExit)
}

