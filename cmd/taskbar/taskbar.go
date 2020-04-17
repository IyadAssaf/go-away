package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Go Away")

	ctx, cancel := context.WithCancel(context.Background())
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		cancel()
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	s := status.NewSlackStatus(status.DefaultStatusText, status.DefaultStatusEmoji)
	err := s.SetStatusWhenWebcamIsBusy(ctx)
	if err != nil {
		log.Fatal(err)
		// TODO ... show a message something?
	}
}

func onExit() {
}
