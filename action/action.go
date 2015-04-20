package action

import (
	"github.com/apognu/gobeard/action/torrent"
	"github.com/apognu/gobeard/source"
)

type Downloader interface {
	Download(e source.EpisodeSubscription, hash string, url string)
}

type Action interface {
	Trigger(source.EpisodeSubscription)
}

func GetActions() []Action {
	return []Action{
		Log{},
		// Slack{},
		Strike{torrent.Transmission{}},
		// Strike{torrent.WatchDir{}},
	}
}
