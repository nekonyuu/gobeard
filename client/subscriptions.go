package client

import (
	"fmt"

	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
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
	for _, s := range series {
		fmt.Printf("  %s: %s\n", s.GetId(), s.Series.Title)
	}

	return nil
}
