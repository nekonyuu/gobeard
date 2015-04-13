package action

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type GetStrike struct{}

func (GetStrike) Trigger(e source.EpisodeSubscription) {
	const ApiSearchEndpoint = "https://getstrike.net/api/v2/torrents/search/?phrase=%s+s%02.0fe%02.0f+%s"
	const ApiSearchNoQualityEndpoint = "https://getstrike.net/api/v2/torrents/search/?phrase=%s+s%02.0fe%02.0f"
	const ApiDownloadEndpoint = "https://getstrike.net/torrents/api/download/%s.torrent"

	var series source.Subscription
	source.GetPersistence("series").Find(bson.M{"source": e.Source, "series.id": e.SeriesId}).One(&series)

	var resp *http.Response
	var err error

	// Iterate over the desired qualities for the first match
	for _, q := range util.GetConfig().Torrents.Quality {
		u := fmt.Sprintf(ApiSearchEndpoint, url.QueryEscape(series.Series.Title), e.Info.Season, e.Info.Number, q)
		resp, err = http.Get(u)
		if err != nil {
			logrus.Errorf("error getting torrents listing: %s", err)
			return
		}
		if resp.StatusCode != 200 {
			logrus.Infof("no torrent found for quality %s, dropping it: %s", q, u)
			resp = nil
			continue
		}

		break
	}

	// Drop the quality requirement if none matched
	if resp == nil {
		u := fmt.Sprintf(ApiSearchNoQualityEndpoint, url.QueryEscape(series.Series.Title), e.Info.Season, e.Info.Number)
		resp, err = http.Get(u)
		if err != nil {
			logrus.Errorf("error getting torrents listing: %s", err)
			return
		}
		if resp.StatusCode != 200 {
			logrus.Errorf("no torrents were found for the request")
			return
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error reading torrents listing body: %s", err)
		return
	}

	var tor map[string]interface{}
	err = json.Unmarshal(body, &tor)
	if err != nil {
		logrus.Errorf("unable to unmarshal JSON: %s", err)
		return
	}

	to := tor["torrents"].([]interface{})[0]
	t := to.(map[string]interface{})
	torrent_hash := t["torrent_hash"].(string)
	torrent_url := fmt.Sprintf(ApiDownloadEndpoint, torrent_hash)

	resp, err = http.Get(torrent_url)
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

	err = os.Rename(tmp.Name(), util.GetConfig().Torrents.WatchDir+"/"+torrent_hash+".torrent")
	if err != nil {
		logrus.Errorf("error creating output torrent file: %s", err)
		return
	}

	source.GetPersistence("subscriptions").UpdateId(e.Id, bson.M{"$set": bson.M{"state": source.StateSeen}})
}
