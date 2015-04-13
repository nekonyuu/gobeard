package client

import (
	"strconv"
	"strings"

	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
	"gopkg.in/mgo.v2/bson"
)

func CmdUnsubscribe() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdUnsubscribe,
		UsageLine: "unsubscribe <source_id/series_id>",
		Short:     "unsubscribe from a given TV series",
	}

	return
}

func runCmdUnsubscribe(cmd *commander.Command, args []string) error {
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

	source.GetPersistence("series").RemoveAll(bson.M{"source": sourceId, "series.id": seriesId})

	return nil
}
