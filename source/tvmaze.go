package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/apognu/gobeard/util"
)

type TVMaze struct {
}

func NewTVMaze() Source {
	return TVMaze{}
}

func (TVMaze) Name() string {
	return "tvmaze"
}

func (TVMaze) SearchSeries(title string) []Series {
	const apiEndpoint = "http://api.tvmaze.com/search/shows?q=%s"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, title))
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

	var raw []map[string]interface{}
	var series []Series
	err = json.Unmarshal(body, &raw)
	if err != nil {
		logrus.Errorf("failed to parse upstream data: %s", err)
	}

	for _, item := range raw {
		show := item["show"].(map[string]interface{})

		series = append(series, Series{
			Id:    show["id"].(float64),
			Title: show["name"].(string),
		})
	}

	return series
}

func (TVMaze) GetSeries(id int) Series {
	const apiEndpoint = "http://api.tvmaze.com/shows/%d"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, id))
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

	var raw map[string]interface{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		logrus.Errorf("failed to parse upstream data: %s", err)
		return Series{}
	}

	return Series{
		Id:    raw["id"].(float64),
		Title: raw["name"].(string),
	}
}

func (TVMaze) ListEpisodes(id int) []Episode {
	const apiEndpoint = "http://api.tvmaze.com/shows/%d/episodes"

	resp, err := http.Get(fmt.Sprintf(apiEndpoint, id))
	if err != nil {
		logrus.Errorf("failed to list episodes for %d: %s", id, err)
		return []Episode{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to list episodes for %d: %s", id, err)
		return []Episode{}
	}

	var raw []map[string]interface{}
	var episodes []Episode
	err = json.Unmarshal(body, &raw)
	if err != nil {
		logrus.Errorf("failed to parse upstream data: %s", err)
		return []Episode{}
	}

	for _, item := range raw {
		airdate, err := time.Parse("2006-01-02", item["airdate"].(string))
		if err != nil {
			continue
		}

		episodes = append(episodes, Episode{
			Id:      item["id"].(float64),
			Season:  item["season"].(float64),
			Number:  item["number"].(float64),
			Title:   item["name"].(string),
			Airdate: airdate,
		})
	}

	return episodes
}

func (s TVMaze) GetPoller(id int, quit <-chan int) func(chan []Episode) {
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
