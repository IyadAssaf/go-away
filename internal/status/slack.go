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
	DefaultStatusText      = "On webcam"
	DefaultStatusEmoji     = "🎥"
	DefaultWaitTimeSeconds = 5
)

type SlackStatus struct {
	client        *slack.Client
	statusText    string
	statusEmoji   string
	refreshRate   int64
	slackToken    string
	statusIsSet   bool
	statusIsUnset bool
}

func NewSlackStatus() *SlackStatus {
	return &SlackStatus{
		client:        slack.New(os.Getenv("SLACK_API_TOKEN")),
		statusText:    DefaultStatusText,
		statusEmoji:   DefaultStatusEmoji,
		refreshRate:   DefaultWaitTimeSeconds,
		statusIsSet:   false,
		statusIsUnset: false,
	}
}

func (s *SlackStatus) WithSlackToken(token string) *SlackStatus {
	s.slackToken = token

	log.Debugf("setting token %s", s.slackToken)
	s.client = slack.New(s.slackToken)
	return s
}

func (s *SlackStatus) WithStatusText(text string) *SlackStatus {
	s.statusText = text
	return s
}

func (s *SlackStatus) WithStatusEmoji(emoji string) *SlackStatus {
	s.statusEmoji = emoji
	return s
}

func (s *SlackStatus) WithRefreshRate(rate int64) *SlackStatus {
	s.refreshRate = rate
	return s
}

func (s *SlackStatus) SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func (s *SlackStatus) DoNotDistrub(ctx context.Context) error {
	s.statusIsUnset = false
	if !s.statusIsSet {
		log.Debugf("Setting status on slack using token %s", s.slackToken)
		//TODO rate limit how often we send this
		if err := s.client.SetUserCustomStatusContext(ctx, s.statusText, s.statusEmoji, 0); err != nil {
			log.Debugf("error setting status %s", err)
			return err
		}
		log.Debugf("successfully set status")
		s.statusIsSet = true
	}

	return nil
}

func (s *SlackStatus) Clear(ctx context.Context) error {
	//TODO rate limit how often we send this
	s.statusIsSet = false

	if !s.statusIsUnset {
		log.Debugf("Unsetting status on slack using token %s", s.slackToken)
		//TODO rate limit how often we send this
		if err := s.client.UnsetUserCustomStatusContext(ctx); err != nil {
			log.Debugf("error unsetting status %s", err)
			return err
		}
		log.Debugf("successfylly unset status")
		s.statusIsUnset = true
	}

	return nil
}

func (s *SlackStatus) SetStatusWhenWebcamIsBusy(ctx context.Context, isOnNotif chan bool) error {
	defer s.Clear(ctx)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		for {
			log.Debugf("Looping")
			isOn, err := webcamchecker.IsWebcamOn(ctx)
			if err != nil {
				errCh <- err
				return
			}

			if isOnNotif != nil {
				isOnNotif <- isOn
			}

			log.Debugf("webcam is on %+v", isOn)

			if isOn {
				if err := s.DoNotDistrub(ctx); err != nil {
					errCh <- err
					return
				}
			} else {
				if err := s.Clear(ctx); err != nil {
					errCh <- err
					return
				}
			}
			time.Sleep(time.Second * time.Duration(s.refreshRate))
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