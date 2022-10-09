package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tpiperatgod/fmw/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
