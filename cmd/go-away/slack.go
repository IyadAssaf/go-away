package main

import (
	"context"
	"github.com/IyadAssaf/webcamchecker"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"time"
)

var log = logrus.New()

//TODO make customisable from envs, cli
const (
	defaultWaitTime      = time.Second * 10
	defaultOnWebcamText  = "On camera"
	defaultOnWebcamEmoji = "🎥"
)

type slackStatus struct {
	client *slack.Client
}

func (s *slackStatus) DoNotDistrub(ctx context.Context) error {
	log.Debugf("Setting do not disturb status on slack")
	return s.client.SetUserCustomStatusContext(ctx, defaultOnWebcamText, defaultOnWebcamEmoji, 0)
}

func (s *slackStatus) Clear(ctx context.Context) error {
	log.Debugf("Un-setting do not disturb status on slack")
	return s.client.UnsetUserCustomStatusContext(ctx)
}

func (s *slackStatus) SetStatusWhenWebcamIsBusy(ctx context.Context) error {
	defer s.Clear(ctx)
	for {
		isOn, err := webcamchecker.IsWebcamOn(ctx)
		if err != nil {
			return err
		}

		if isOn {
			if err := s.DoNotDistrub(ctx); err != nil {
				return err
			}
		} else {
			if err := s.Clear(ctx); err != nil {
				return err
			}
		}
		time.Sleep(defaultWaitTime)
	}
}

