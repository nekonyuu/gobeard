package action

import (
	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
	"github.com/longnguyen11288/go-transmission/transmission"
	"gopkg.in/mgo.v2/bson"
)

var client *transmission.TransmissionClient = nil

type Transmission struct{}

func (Transmission) Download(e source.EpisodeSubscription, hash string, url string) error {
	c := util.GetConfig().Torrents.Transmission
	if client == nil {
		cl := transmission.New(c.Endpoint, c.Username, c.Password)
		client = &cl
	}

	info, err := client.AddTorrentByURL(url, c.DownloadDir)
	if err != nil {
		logrus.Errorf("error starting download: %s", err)
		return err
	}
	if info.HashString == "" {
		logrus.Errorf("error starting downloaded file, hash was empty")
		return err
	}

	if util.GetConfig().Slack.WebhookUrl != "" {
		Slack{}.Trigger(e)
	}

	logrus.Infof("transmission: torrent added: %s (%s)", info.HashString, info.Name)
	source.GetPersistence("subscriptions").UpdateId(e.Id, bson.M{"$set": bson.M{"state": source.StateSeen}})

	return nil
}
