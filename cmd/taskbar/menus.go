package main

import (
	"context"
	"fmt"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/getlantern/systray"
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
	menu := systray.AddMenuItem(fmt.Sprintf(`Slack status set to "%s %s"`, status.DefaultStatusText, status.DefaultStatusEmoji), "")
	menu.Hide()
	menu.Disable()
	return menu
}

func setupErrorMenu(ctx context.Context) *systray.MenuItem {
	menu := systray.AddMenuItem("", "")
	menu.Hide()
	menu.Disable()
	return menu
}