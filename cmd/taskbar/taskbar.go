package main

import (
	"context"
	"fmt"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/getlantern/systray"
	"log"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTooltip("Go Away")
	systray.SetIcon(cameraOffIconData)

	ctx, cancel := context.WithCancel(context.Background())

	statusCameraOnText := systray.AddMenuItem(fmt.Sprintf(`Slack status set to "%s %s"`, status.DefaultStatusText, status.DefaultStatusEmoji), "")
	statusCameraOnText.Hide()
	statusCameraOnText.Disable()

	mQuitOrig := systray.AddMenuItem("Quit", "")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		cancel()
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	s := status.NewSlackStatus(status.DefaultStatusText, status.DefaultStatusEmoji, status.DefaultWaitTimeSeconds)


	isOnCh := make(chan bool)
	go func() {
		for {
			isOn := <-isOnCh
			switch isOn {
			case true:
				systray.SetIcon(cameraOnIconData)
				statusCameraOnText.Show()
			case false:
				systray.SetIcon(cameraOffIconData)
				statusCameraOnText.Hide()
			}
		}
	}()

	err := s.SetStatusWhenWebcamIsBusy(ctx, isOnCh)
	if err != nil {
		log.Fatal(err)
		// TODO ... show a message something?
	}
}

func onExit() {
}


