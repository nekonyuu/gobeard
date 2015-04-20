package action

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/apognu/gobeard/source"
)

type Log struct{}

func (Log) Trigger(s source.EpisodeSubscription) {
	logrus.WithFields(logrus.Fields{
		"source":      s.Source,
		"series":      s.SeriesId,
		"episodeid":   fmt.Sprintf("S%02.0fE%02.0f", s.Info.Season, s.Info.Number),
		"episodename": s.Info.Title,
	}).Info("action taken")
}
