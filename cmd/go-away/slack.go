package main

import (
	"context"
	"github.com/IyadAssaf/webcamchecker"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"time"
)

var log = logrus.New()

const (
	defaultWaitTime    = time.Second * 10
	defaultStatusText  = "On webcam"
	defaultStatusEmoji = "🎥"
)

type slackStatus struct {
	client      *slack.Client
	statusText  string
	statusEmoji string
}

func (s *slackStatus) DoNotDistrub(ctx context.Context) error {
	log.Debugf("Setting status on slack")
	return s.client.SetUserCustomStatusContext(ctx, defaultStatusText, defaultStatusEmoji, 0)
}

func (s *slackStatus) Clear(ctx context.Context) error {
	log.Debugf("Unsetting status on slack")
	return s.client.UnsetUserCustomStatusContext(ctx)
}

func (s *slackStatus) SetStatusWhenWebcamIsBusy(ctx context.Context) error {
	defer s.Clear(ctx)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		for {
			isOn, err := webcamchecker.IsWebcamOn(ctx)
			if err != nil {
				errCh<-err
				return
			}

			log.Debugf("webcam is on %+v", isOn)

			if isOn {
				if err := s.DoNotDistrub(ctx); err != nil {
					errCh<-err
					return
				}
			} else {
				if err := s.Clear(ctx); err != nil {
					errCh<-err
					return
				}
			}
			time.Sleep(defaultWaitTime)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				//cancelling is A-okay
				errCh <- nil
			default:
				errCh <- ctx.Err()
			}
		}
	}()

	err := <-errCh
	log.Debugf("context is finished")

	return err
}
