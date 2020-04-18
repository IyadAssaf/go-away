package main

import (
	"context"
	"fmt"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/sirupsen/logrus"

	"github.com/getlantern/systray"
	"log"
)

var statusMenu, quitMenu, prefMenu *systray.MenuItem

func onReady() {
	systray.SetTooltip("Go Away")
	systray.SetIcon(cameraOffIconData)

	tokenCh := make(chan string)
	//defer close(tokenCh)
	trigger := make(chan struct{})
	//defer close(trigger)
	isOnCh := make(chan bool)
	//defer close(isOnCh)

	slackStatus := status.NewSlackStatus()
	slackStatus.SetLogLevel(logrus.DebugLevel)

	ctx, cancel := context.WithCancel(context.Background())
	statusMenu = setupStatusMenu(ctx)

	prefMenu = setupPreferencesMenu(ctx, tokenCh)

	quitMenu = setupQuitMenu(ctx, cancel)
	go func() {
		for {
			t := <-tokenCh
			log.Println("Setting token", t)
			slackStatus = slackStatus.WithSlackToken(t)
			trigger <- struct{}{}
		}
	}()

	go switchIcon(ctx, isOnCh)

	go loop(ctx, slackStatus, trigger, isOnCh)

	trigger<-struct{}{}
}

func loop(ctx context.Context, s *status.SlackStatus, trigger chan struct{}, isOnCh chan bool) {
	for {
		<-trigger
		log.Println("Trigger")
		err := s.SetStatusWhenWebcamIsBusy(ctx, isOnCh)
		if err != nil {
			logError(err)
		}
	}
}

func onExit() {
}

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

func switchIcon(ctx context.Context, isOnCh chan bool) error {
	for {
		isOn := <-isOnCh
		switch isOn {
		case true:
			systray.SetIcon(cameraOnIconData)
			statusMenu.Show()
		case false:
			systray.SetIcon(cameraOffIconData)
			statusMenu.Hide()
		}
	}
}