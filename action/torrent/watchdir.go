package torrent

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
	"gopkg.in/mgo.v2/bson"
)

type WatchDir struct{}

func (WatchDir) Download(e source.EpisodeSubscription, hash string, url string) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("error retrieving torrent file: %s", err)
		return
	}
	defer resp.Body.Close()

	tmp, err := ioutil.TempFile(os.TempDir(), "gobeard")
	if err != nil {
		logrus.Errorf("cannot create temporary file: %s", err)
		return
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		logrus.Errorf("error writing torrent file: %s", err)
		return
	}

	err = os.Rename(tmp.Name(), util.GetConfig().Torrents.WatchDir+"/"+hash+".torrent")
	if err != nil {
		logrus.Errorf("error creating output torrent file: %s", err)
		return
	}

	source.GetPersistence("subscriptions").UpdateId(e.Id, bson.M{"$set": bson.M{"state": source.StateSeen}})
}
