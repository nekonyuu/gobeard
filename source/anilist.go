package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/apognu/gobeard/util"
)

type Anilist struct {
	token string
}

func NewAnilist() Source {
	a := Anilist{}
	a.Authenticate()
	return a
}

func (Anilist) Name() string {
	return "anilist"
}

func (a Anilist) Authenticate() {
	const apiAuthEndpoint = "https://anilist.co/api/auth/access_token"

	payload := url.Values{}
	payload.Add("client_id", util.GetConfig().Anilist.ClientID)
	payload.Add("client_secret", util.GetConfig().Anilist.ClientSecret)
	payload.Add("grant_type", "client_credentials")

	resp, err := http.PostForm(apiAuthEndpoint, payload)
	if err != nil {
		logrus.Errorf("failed to send authentication request on anilist api")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Errorf("failed to authenticate on anilist api, http code %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to authenticate on anilist api : no body")
		return
	}

	var raw map[string]interface{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		logrus.Errorf("failed to parse upstream data: %s", err)
		return
	}

	a.token = raw["access_token"].(string)

	return
}

func (a Anilist) SearchSeries(title string) []Series {
	const apiEndpoint = "https://anilist.co/api/anime/search/%s?access_token=%s"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, title, a.token))
	if err != nil {
		logrus.Errorf("failed to search for %s", title)
		return []Series{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to search for %s", title)
		return []Series{}
	}

	var series []Series

	if resp.StatusCode == 200 {
		if len(strings.TrimSpace(string(body))) == 0 {
			return series
		}

		var raw []map[string]interface{}

		err = json.Unmarshal(body, &raw)
		if err != nil {
			logrus.Errorf("failed to parse upstream data: %s", err)
		}

		for _, show := range raw {
			series = append(series, Series{
				Id:    show["id"].(float64),
				Title: show["title_romaji"].(string),
			})
		}
	} else {
		logrus.Warnf("could not find %s : http code %d", title, resp.StatusCode)
	}

	return series
}

func (a Anilist) GetSeries(id int) Series {
	const apiEndpoint = "https://anilist.co/api/anime/%d?access_token=%s"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, id, a.token))
	if err != nil {
		logrus.Errorf("failed to search for %d: %s", id, err)
		return Series{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to search for %d", id)
		return Series{}
	}

	if resp.StatusCode == 200 {
		var raw map[string]interface{}
		err = json.Unmarshal(body, &raw)
		if err != nil {
			logrus.Errorf("failed to parse upstream data: %s", err)
			return Series{}
		}

		return Series{
			Id:           raw["id"].(float64),
			Title:        raw["title_romaji"].(string),
			Summary:      raw["description"].(string),
			EpisodeCount: raw["total_episodes"].(float64),
		}
	}
	return Series{}
}

func (a Anilist) ListEpisodes(id int) []Episode {
	const apiEndpoint = "https://anilist.co/api/anime/%d?access_token=%s"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, id, a.token))
	if err != nil {
		logrus.Errorf("failed to search for %d: %s", id, err)
		return []Episode{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to search for %d", id)
		return []Episode{}
	}

	var raw map[string]interface{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		logrus.Errorf("failed to parse upstream data: %s", err)
		return []Episode{}
	}

	episodes := []Episode{}

	airing, ok := raw["airing"].(map[string]interface{})
	if ok {
		airstamp, err := time.Parse("2006-01-02T15:04:05-07:00", airing["time"].(string))
		if err != nil {
			return episodes
		}

		episodes = append(episodes, Episode{
			Id:       raw["id"].(float64),
			Season:   1,
			Number:   airing["next_episode"].(float64),
			Title:    raw["title_romaji"].(string),
			Airstamp: airstamp,
		})
	}

	return episodes
}

func (s Anilist) GetPoller(id int, quit <-chan int) func(chan []Episode) {
	return func(c chan []Episode) {
	Poller:
		for {
			select {
			case <-quit:
				break Poller
			case <-time.After(util.GetConfig().CheckInterval):
				c <- s.ListEpisodes(id)
			}
		}
	}
}
