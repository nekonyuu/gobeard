package action

import (
	"github.com/apognu/gobeard/action/torrent"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
)

type Downloader interface {
	Download(e source.EpisodeSubscription, hash string, url string) error
}

type Action interface {
	Trigger(source.EpisodeSubscription)
}

var actions = map[string]Action{
	"log":    Log{},
	"slack":  Slack{},
	"strike": Strike{},
}

var downloaders = map[string]Downloader{
	"watchdir":     torrent.WatchDir{},
	"transmission": torrent.Transmission{},
}

func GetActions() []Action {
	ac := make([]Action, 0)
	for _, a := range util.GetConfig().Actions {
		ac = append(ac, actions[a])
	}

	return ac
}

func GetDownloaders() []Downloader {
	do := make([]Downloader, 0)
	for _, d := range util.GetConfig().Downloaders {
		do = append(do, downloaders[d])
	}

	return do
}
