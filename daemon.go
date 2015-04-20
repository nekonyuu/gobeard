package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gobeard/action"
	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
)

func NewDaemon() {
	event := make(chan source.EpisodeSubscription, util.GetConfig().EventsQueueSize)
	go checkForSubscriptions(event)

	for {
		select {
		case e := <-event:
			for _, a := range action.GetActions() {
				go a.Trigger(e)
			}
		}
	}
}

func checkForSubscriptions(e chan source.EpisodeSubscription) {
	subs := make([]source.Subscription, 0)
	subl := 0
	quits := make([]chan int, 0)
	for {
		for _, s := range subs {
			quit := make(chan int)
			go s.Monitor(e, quit)
			quits = append(quits, quit)
		}

		for {
			err := source.GetPersistence("series").Find(nil).All(&subs)
			if err != nil {
				logrus.Fatalf("cannot retrieve subscriptions: %s", err)
			}
			if len(subs) != subl {
				subl = len(subs)
				for _, c := range quits {
					close(c)
				}
				quits = make([]chan int, 0)
				break
			}
			time.Sleep(util.GetConfig().CheckInterval)
		}
	}
}
