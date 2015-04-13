package client

import (
	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"gopkg.in/mgo.v2/bson"
)

func CmdMark() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdMark,
		UsageLine: "mark -type <series|episode> -id <id> -state <ignored|unseen|seen>",
		Short:     "sets the state of an item",
		Flag:      *flag.NewFlagSet("mark", flag.ExitOnError),
	}

	cmd.Flag.String("type", "", "type of item to be marked (`series` or `episode`)")
	cmd.Flag.String("id", "", "ID of the item to mark (eg. `tvmaze/42`)")
	cmd.Flag.String("state", "", "state to mark the item as (`ignored`, `unseen`, `seen`)")

	return
}

func runCmdMark(cmd *commander.Command, args []string) error {
	t := cmd.Flag.Lookup("type").Value.Get().(string)
	id := cmd.Flag.Lookup("id").Value.Get().(string)
	s := cmd.Flag.Lookup("state").Value.Get().(string)

	if t == "" {
		argumentError(cmd, "option `type` is mandatory")
	}
	if t != "series" && t != "episode" {
		argumentError(cmd, "unknown value for attribute `type`")
	}
	if id == "" {
		argumentError(cmd, "option `id` is mandatory")
	}
	if s == "" {
		argumentError(cmd, "option `state` is mandatory")
	}
	if s != "ignored" && s != "unseen" && s != "seen" {
		argumentError(cmd, "unknown value for attribute `state`")
	}

	if t == "series" {
		source.GetPersistence("subscriptions").UpdateAll(bson.M{"_id": bson.M{"$regex": `^` + id + `/\d+`}}, bson.M{"$set": bson.M{"state": s}})
	} else if t == "episode" {
		source.GetPersistence("subscriptions").UpdateAll(bson.M{"_id": id}, bson.M{"$set": bson.M{"state": s}})
	}

	return nil
}
