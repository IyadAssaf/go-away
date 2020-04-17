package status

import (
	"context"
	"github.com/IyadAssaf/webcamchecker"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"os"
	"time"
)

// TODO make property of SlackStatus
var log = logrus.New()

const (
	DefaultStatusText  = "On webcam"
	DefaultStatusEmoji = "🎥"
	DefaultWaitTime    = time.Second * 10
)

type SlackStatus struct {
	client      *slack.Client
	statusText  string
	statusEmoji string
}

func NewSlackStatus(statusText, statusEmoji string) *SlackStatus {
	return &SlackStatus{
		client:      slack.New(os.Getenv("SLACK_API_TOKEN")),
		statusText:  stringOrDefault(statusText, DefaultStatusText),
		statusEmoji: stringOrDefault(statusEmoji, DefaultStatusEmoji),
	}
}

func (s *SlackStatus) SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func (s *SlackStatus) DoNotDistrub(ctx context.Context) error {
	log.Debugf("Setting status on slack")
	return s.client.SetUserCustomStatusContext(ctx, s.statusText, s.statusEmoji, 0)
}

func (s *SlackStatus) Clear(ctx context.Context) error {
	log.Debugf("Unsetting status on slack")
	return s.client.UnsetUserCustomStatusContext(ctx)
}

func (s *SlackStatus) SetStatusWhenWebcamIsBusy(ctx context.Context) error {
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
			time.Sleep(DefaultWaitTime)
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

func stringOrDefault(s, def string) string {
	if s != "" {
		return s
	}
	return def
}