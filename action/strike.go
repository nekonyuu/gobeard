package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
	"gopkg.in/mgo.v2/bson"
)

type Strike struct {
	Subaction Downloader
}

func (a Strike) Trigger(e source.EpisodeSubscription) {
	const ApiSearchEndpoint = "https://getstrike.net/api/v2/torrents/search/?phrase=%s+s%02.0fe%02.0f+%s"
	const ApiSearchNoQualityEndpoint = "https://getstrike.net/api/v2/torrents/search/?phrase=%s+s%02.0fe%02.0f"
	const ApiDownloadEndpoint = "https://getstrike.net/torrents/api/download/%s.torrent"

	var series source.Subscription
	source.GetPersistence("series").Find(bson.M{"source": e.Source, "series.id": e.SeriesId}).One(&series)

	var resp *http.Response
	var err error

	// Iterate over the desired qualities for the first match
QualityLoop:
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

		for _, d := range GetDownloaders() {
			err = d.Download(e, torrent_hash, torrent_url)
			if err != nil {
				continue QualityLoop
			}
		}
	}
}
