package action

import "github.com/apognu/gobeard/source"

type Action interface {
	Trigger(source.EpisodeSubscription)
}

func GetActions() []Action {
	return []Action{
		Log{},
		// Slack{},
		GetStrike{},
	}
}
