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

	trigger := make(chan struct{})
	//defer close(trigger)
	isOnCh := make(chan bool)
	//defer close(isOnCh)

	slackStatus := status.NewSlackStatus()
	slackStatus.SetLogLevel(logrus.DebugLevel)
	slackStatus.WithSlackToken(settingToken)

	ctx, cancel := context.WithCancel(context.Background())

	statusMenu = setupStatusMenu(ctx)
	errMenu = setupErrorMenu(ctx)
	prefMenu = setupPreferencesMenu(ctx, slackStatus, trigger)
	quitMenu = setupQuitMenu(ctx, cancel)

	go switchOnIcon(ctx, isOnCh)

	go loop(ctx, slackStatus, trigger, isOnCh)

	trigger <- struct{}{}
}

func loop(ctx context.Context, s *status.SlackStatus, trigger chan struct{}, isOnCh chan bool) {
	for {
		<-trigger
		log.Println("Trigger")
		clearError()
		err := s.SetStatusWhenWebcamIsBusy(ctx, isOnCh)
		if err != nil {
			logError(err)
			continue
		}
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
