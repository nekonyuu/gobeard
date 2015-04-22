package client

import (
	"fmt"

	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
	"github.com/ryanuber/columnize"
)

func CmdSubscriptions() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdSubscriptions,
		UsageLine: "subscriptions",
		Short:     "list current subscriptions",
	}

	return
}

func runCmdSubscriptions(cmd *commander.Command, args []string) error {
	var series []source.Subscription
	source.GetPersistence("series").Find(nil).All(&series)

	fmt.Println("Active subscriptions:")
	o := make([]string, 0)
	for _, s := range series {
		o = append(o, fmt.Sprintf("%s|%s\n", Bold(s.GetId()), s.Series.Title))
	}
	fmt.Println(columnize.Format(o, &columnize.Config{Prefix: "  "}))

	return nil
}
