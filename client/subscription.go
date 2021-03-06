package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/gonuts/commander"
	"github.com/ryanuber/columnize"
	"gopkg.in/mgo.v2/bson"
)

func CmdSubscription() (cmd *commander.Command) {
	cmd = &commander.Command{
		Run:       runCmdSubscription,
		UsageLine: "subscription <source_id/series_id>",
		Short:     "print the state of all episodes in a subscription",
	}

	return
}

func runCmdSubscription(cmd *commander.Command, args []string) error {
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

	var series source.Subscription
	source.GetPersistence("series").Find(bson.M{"source": sourceId, "series.id": seriesId}).One(&series)

	if series.Id == "" {
		logrus.Fatalf("given series does not exist or was marked for deletion")
		return nil
	}

	var subs []source.EpisodeSubscription
	err = source.GetPersistence("subscriptions").Find(bson.M{"source": sourceId, "series": seriesId}).All(&subs)
	if err != nil {
		logrus.Fatalf("cannot retrieve subscription: %s", err)
	}

	fmt.Printf("Subscription status for %s (%s)\n", Bold(os.Args[2]), Blue(series.Series.Title))

	if len(subs) == 0 {
		fmt.Println("  No episode for this subscription yet, is the daemon running?")
		return nil
	}

	o := []string{fmt.Sprintf("%s|%s|%s|%s", Bold("ID"), "Episode", "Title", "State")}
	// o := make([]string, 0)
	for _, s := range subs {
		var st string
		switch s.State {
		case source.StateUnseen:
			st = Yellow(s.State)
		case source.StateIgnored:
			st = Red(s.State)
		case source.StateSeen:
			st = Green(s.State)
		}

		o = append(o, fmt.Sprintf("%s|S%02.0fE%02.0f|%s|[%-16s]\n", Bold(fmt.Sprintf("%s/%d/%0.0f", s.Source, s.SeriesId, s.Info.Id)), s.Info.Season, s.Info.Number, s.Info.Title, st))
	}
	// fmt.Println(o)
	fmt.Println(columnize.Format(o, &columnize.Config{Prefix: "  "}))

	return nil
}
