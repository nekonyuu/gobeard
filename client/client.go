package client

import (
	"os"

	"github.com/fatih/color"
	"github.com/gonuts/commander"
	"github.com/sirupsen/logrus"
)

var Blue func(...interface{}) string = color.New(color.FgBlue).SprintFunc()
var Yellow func(...interface{}) string = color.New(color.FgYellow).SprintFunc()
var Green func(...interface{}) string = color.New(color.FgGreen).SprintFunc()
var Red func(...interface{}) string = color.New(color.FgRed).SprintFunc()

func NewClient() {
	var CmdLine = &commander.Command{
		UsageLine: os.Args[0],
		Short:     "Keep track of your TV shows subscriptions",
	}

	CmdLine.Subcommands = []*commander.Command{
		CmdSearch(),
		CmdMark(),
		CmdList(),
		CmdSubscribe(),
		CmdUnsubscribe(),
		CmdSubscriptions(),
		CmdSubscription(),
	}

	CmdLine.Dispatch(os.Args[1:])
}

func argumentError(cmd *commander.Command, msg string) {
	logrus.Error(msg)
	cmd.Usage()
	os.Exit(1)
}
