package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
)

func main() {
	app := &cli.App{
		Name:        "goaway",
		Usage:       "./goaway",
		Description: "update slack with a status when you're on webcam",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug logging",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "status-text",
				Usage: "text to use for slack status",
				Value: defaultStatusText,
			},
			&cli.StringFlag{
				Name:  "status-emoji",
				Usage: "emoji to use for slack status",
				Value: defaultStatusEmoji,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(logrus.DebugLevel)
			}

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)

			ctx, cancel := context.WithCancel(c.Context)
			go func() {
				<-ch
				cancel()
			}()

			s := &slackStatus{
				client:      slack.New(os.Getenv("SLACK_API_TOKEN")),
				statusText:  c.String("status-text"),
				statusEmoji: c.String("status-emoji"),
			}
			return s.SetStatusWhenWebcamIsBusy(ctx)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
