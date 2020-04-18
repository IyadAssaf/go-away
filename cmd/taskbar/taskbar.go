package main

import (
	"context"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/sirupsen/logrus"

	"github.com/getlantern/systray"
	"log"
)

var errMenu, statusMenu, quitMenu, prefMenu *systray.MenuItem

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
	errMenu = setupErrorMenu(ctx)

	prefMenu = setupPreferencesMenu(ctx, tokenCh)

	quitMenu = setupQuitMenu(ctx, cancel)
	go func() {
		for {
			t := <-tokenCh
			log.Println("Setting token", t)
			slackStatus = slackStatus.WithSlackToken(t)
			trigger <- struct{}{}
			clearError()
		}
	}()

	go switchOnIcon(ctx, isOnCh)

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
		clearError()
	}
}

func switchOnIcon(ctx context.Context, isOnCh chan bool) error {
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
