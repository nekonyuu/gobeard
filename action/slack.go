package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"

	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
)

type SlackMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

type Slack struct{}

func (Slack) Trigger(sub source.EpisodeSubscription) {
	var s source.Subscription
	source.GetPersistence("series").Find(bson.M{"source": sub.Source, "series.id": sub.SeriesId}).One(&s)

	msg := SlackMessage{
		Channel:  util.GetConfig().Slack.Channel,
		Username: "GoBeard Torrent Pwner",
		Text:     fmt.Sprintf("Subscription update: *%s* - S%02.0fE%02.0f - %s.", s.Series.Title, sub.Info.Season, sub.Info.Number, sub.Info.Title),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("cannot marshal Slack message: %s", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", util.GetConfig().Slack.WebhookUrl, bytes.NewBuffer(body))
	_, err = client.Do(req)
	if err != nil {
		logrus.Errorf("cannot post Slack message: %s", err)
	}
}
