package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
)

func CmdList() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdList,
		UsageLine: "list <source_id/series_id>",
		Short:     "list episodes of one TV series",
	}

	return
}

func runCmdList(cmd *commander.Command, args []string) error {
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
		argumentError(cmd, "cannot find given source")
	}
	episodes := sources[sourceId].ListEpisodes(seriesId)

	fmt.Printf("Episode listing for series `%s`:\n", os.Args[2])
	for _, item := range episodes {
		fmt.Printf("  S%02.0fE%02.0f: %s (%s)\n", item.Season, item.Number, item.Title, item.Airdate.Format("2006-01-02"))
	}

	return nil
}
