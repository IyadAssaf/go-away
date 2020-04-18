package main

import (
	"github.com/gen2brain/beeep"
	"github.com/sirupsen/logrus"
)

func logError(err error) {
	logrus.Error(err)
	_ = beeep.Notify("go-away", err.Error(), "")
}
