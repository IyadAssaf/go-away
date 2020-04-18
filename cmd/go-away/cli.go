package main

import (
	"context"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	app := &cli.App{
		Name:        "go-away",
		Usage:       "./go-away",
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
				Value: status.DefaultStatusText,
			},
			&cli.StringFlag{
				Name:  "status-emoji",
				Usage: "emoji to use for slack status",
				Value: status.DefaultStatusEmoji,
			},
			&cli.Int64Flag{
				Name:  "refresh-rate",
				Usage: "number of seconds to refresh webcam status",
				Value: int64(status.DefaultWaitTimeSeconds / time.Second),
			},
		},
		Action: func(c *cli.Context) error {

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)

			ctx, cancel := context.WithCancel(c.Context)
			go func() {
				<-ch
				cancel()
			}()

			s := status.NewSlackStatus(c.String("status-text"), c.String("status-emoji"), c.Int64("refresh-rate"))
			if c.Bool("debug") {
				s.SetLogLevel(logrus.DebugLevel)
			}
			return s.SetStatusWhenWebcamIsBusy(ctx, nil)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
