package source

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/apognu/gobeard/util"
)

const (
	StateUnseen  = "unseen"
	StateIgnored = "ignored"
	StateSeen    = "seen"
)

var (
	db *mgo.Session
)

func GetPersistence(c string) *mgo.Collection {
	if db == nil {
		var err error
		db, err = mgo.DialWithTimeout(util.GetConfig().MongoDb.Host, 2*time.Second)
		if err != nil {
			logrus.Fatalf("error connecting to MondoDB: %s", err)
		}
	}

	return db.DB("gobeard").C(c)
}

type EpisodeSubscription struct {
	Id       string  `bson:"_id"`
	Source   string  `bson:"source"`
	SeriesId int     `bson:"series"`
	Info     Episode `bson:"info"`
	State    string  `bson:"state"`
}

type Subscription struct {
	Id     string `bson:"_id"`
	Source string `bson:"source"`
	Series Series `bson:"series"`
}

func (s Subscription) GetId() string {
	return s.Source + "/" + fmt.Sprintf("%0.0f", s.Series.Id)
}

func (s Subscription) GetSeriesId(id float64) string {
	return s.GetId() + "/" + fmt.Sprintf("%0.0f", id)
}

func NewSubscription(src Source, id int) Subscription {
	series := src.GetSeries(id)

	s := Subscription{
		Source: src.Name(),
		Series: series,
	}
	s.Id = s.GetId()

	err := GetPersistence("series").Insert(s)
	if err != nil {
		logrus.Errorf("error persisting series: %s", err)
	}

	return s
}

func (sub Subscription) Monitor(e chan EpisodeSubscription, quit <-chan int) {
	c := make(chan []Episode)
	q := make(chan int)
	go GetSources()[sub.Source].GetPoller(int(sub.Series.Id), q)(c)

Main:
	for {
		select {
		case <-quit:
			// If the subscription was deleted, report this to the database
			n, err := GetPersistence("series").Find(bson.M{"series.id": sub.Series.Id}).Count()
			if err == nil && n == 0 {
				GetPersistence("subscriptions").RemoveAll(bson.M{"source": sub.Source, "series": sub.Series.Id})
			}
			close(q)
			break Main

		case update := <-c:
			for _, item := range update {
				ep := EpisodeSubscription{
					Id:       sub.GetSeriesId(item.Id),
					Source:   sub.Source,
					SeriesId: int(sub.Series.Id),
					Info:     item,
					State:    StateIgnored,
				}
				GetPersistence("subscriptions").Insert(ep)
			}

			newEpisodes := make([]EpisodeSubscription, 0)
			err := GetPersistence("subscriptions").Find(bson.M{"series": sub.Series.Id, "state": StateUnseen}).All(&newEpisodes)
			if err != nil {
				logrus.Errorf("filed to retrieve subscription: %s", err)
				continue
			}

			for _, ep := range newEpisodes {
				if ep.State == StateUnseen && ep.Info.Airdate.Before(time.Now()) {
					e <- ep
				}
			}
		}
	}
}
