package client

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
)

func CmdSubscribe() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdSubscribe,
		UsageLine: "subscribe <source_id/series_id>",
		Short:     "subscribe to a given TV series",
	}

	return
}

func runCmdSubscribe(cmd *commander.Command, args []string) error {
	if len(args) != 1 {
		argumentError(cmd, "you must provide the series ID")
	}

	id := strings.Split(args[0], "/")
	if len(id) != 2 {
		argumentError(cmd, "cannot parse the given ID to a series identifier")
	}
	sourceId := id[0]
	seriesId, err := strconv.Atoi(id[1])
	if err != nil {
		argumentError(cmd, "cannot parse the given ID to a series identifier")
	}

	sources := source.GetSources()
	if _, ok := sources[sourceId]; !ok {
		logrus.Fatalf("cannot find source `%s`\n", sourceId)
	}

	source.NewSubscription(sources[sourceId], seriesId)

	return nil
}
