package torrent

import (
	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
	"github.com/longnguyen11288/go-transmission/transmission"
	"gopkg.in/mgo.v2/bson"
)

type Transmission struct{}

func (Transmission) Download(e source.EpisodeSubscription, hash string, url string) error {
	c := util.GetConfig().Torrents.Transmission
	tr := transmission.New(c.Endpoint, c.Username, c.Password)
	info, err := tr.AddTorrentByURL(url, c.DownloadDir)
	if err != nil {
		logrus.Errorf("error starting download: %s", err)
		return err
	}
	if info.HashString == "" {
		logrus.Errorf("error starting downloaded file")
		return err
	}

	logrus.Infof("transmission: torrent added: %s (%s)", info.HashString, info.Name)
	source.GetPersistence("subscriptions").UpdateId(e.Id, bson.M{"$set": bson.M{"state": source.StateSeen}})

	return nil
}
