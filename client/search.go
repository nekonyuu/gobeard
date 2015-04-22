package client

import (
	"fmt"
	"strings"

	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
	"github.com/ryanuber/columnize"
)

func CmdSearch() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdSearch,
		UsageLine: "search <search term>...",
		Short:     "search TV series from all built-in sources",
	}

	return
}

func runCmdSearch(cmd *commander.Command, args []string) error {
	if len(args) < 1 {
		argumentError(cmd, "you must provide at least one search term")
	}

	t := strings.Join(args, " ")
	if strings.TrimSpace(t) == "" {
		argumentError(cmd, "search term cannot be empty")
	}

	for name, src := range source.GetSources() {
		series := src.SearchSeries(t)

		fmt.Printf("Search results for `%s` and source `%s`:\n", t, name)
		o := make([]string, 0)
		for _, item := range series {
			o = append(o, fmt.Sprintf("%s|%s\n", Bold(fmt.Sprintf("%s/%0.0f", name, item.Id)), item.Title))
		}
		fmt.Println(columnize.Format(o, &columnize.Config{Prefix: "  "}))
	}
	return nil
}
