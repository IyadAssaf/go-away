package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "goaway",
		Usage:       "./goaway",
		Description: "update slack with a status when you're on webcam",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(logrus.DebugLevel)
			}

			s := &slackStatus{
				client: slack.New(os.Getenv("SLACK_API_TOKEN")),
			}
			return s.SetStatusWhenWebcamIsBusy(context.Background())
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
