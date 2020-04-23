package main

import (
	"context"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/getlantern/systray"
	"os"
)

var (
	settingMessage = status.DefaultStatusText
	settingEmoji   = status.DefaultStatusEmoji
	settingToken   = os.Getenv("SLACK_API_TOKEN")
)

func setupQuitMenu(ctx context.Context, cancel context.CancelFunc) *systray.MenuItem {
	menu := systray.AddMenuItem("Quit", "")
	go func() {
		<-menu.ClickedCh
		cancel()
		systray.Quit()
	}()
	return menu
}

func setupStatusMenu(ctx context.Context) *systray.MenuItem {
	menu := systray.AddMenuItem("", "")
	menu.Hide()
	menu.Disable()
	return menu
}

func setupErrorMenu(ctx context.Context) *systray.MenuItem {
	menu := systray.AddMenuItem("Error", "")
	menu.Hide()
	menu.Disable()
	return menu
}
