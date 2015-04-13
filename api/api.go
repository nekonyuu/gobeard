package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/apognu/gobeard/source"
	"github.com/apognu/gobeard/util"
)

func httpWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NewWebUi() {
	if util.GetConfig().Api.Addr == "" {
		logrus.Warn("api listening address not set, disabling API")
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", serveApp)

	r.HandleFunc("/api/series/search", searchSeries).Methods("GET")
	r.HandleFunc("/api/series/{source:[a-z]+}/{id:[0-9]+}/episodes", listSeriesEpisodes).Methods("GET")

	r.HandleFunc("/api/subscriptions", listSubscriptions).Methods("GET")
	r.HandleFunc("/api/subscriptions/{source:[a-z]+}/{id:[0-9]+}", getSubscription).Methods("GET")
	r.HandleFunc("/api/subscriptions", addSubscription).Methods("PUT")
	r.HandleFunc("/api/subscriptions/{source:[a-z]+}/{id:[0-9]+}", deleteSubscription).Methods("DELETE")
	r.HandleFunc("/api/subscriptions/{source:[a-z]+}/{id:[0-9]+}", markSubscription).Methods("POST")

	http.Handle("/", httpWrapper(r))
	http.ListenAndServe(util.GetConfig().Api.Addr, nil)
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Serve app"))
}

type SearchResults struct {
	Source  string
	Results []source.Series
}

func searchSeries(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("term")
	if t == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var results []SearchResults
	for name, src := range source.GetSources() {
		var result []source.Series
		series := src.SearchSeries(t)
		for _, s := range series {
			result = append(result, s)
		}

		results = append(results, SearchResults{
			Source:  name,
			Results: result,
		})
	}

	json, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func listSeriesEpisodes(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	src, sid := p["source"], p["id"]
	id, _ := strconv.Atoi(sid)

	sources := source.GetSources()
	if _, ok := sources[src]; !ok {
		w.WriteHeader(http.StatusNotFound)
	}
	episodes := sources[src].ListEpisodes(id)

	json, err := json.Marshal(episodes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func listSubscriptions(w http.ResponseWriter, r *http.Request) {
	var subs []source.Subscription
	source.GetPersistence("series").Find(nil).All(&subs)

	json, err := json.Marshal(subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func getSubscription(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	src, sid := p["source"], p["id"]
	id, _ := strconv.Atoi(sid)

	var subs []source.EpisodeSubscription
	err := source.GetPersistence("subscriptions").Find(bson.M{"source": src, "series": id}).All(&subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(subs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json, err := json.Marshal(subs)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(json)
}

func addSubscription(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form["source"]) < 1 || len(r.Form["id"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	src := r.Form["source"][0]
	sid := r.Form["id"][0]

	id, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sources := source.GetSources()
	if _, ok := sources[src]; !ok {
		w.WriteHeader(http.StatusNotFound)
	}

	source.NewSubscription(sources[src], id)

	w.WriteHeader(http.StatusCreated)
}

func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	src, sid := p["source"], p["seriesid"]
	id, _ := strconv.Atoi(sid)

	info, err := source.GetPersistence("series").RemoveAll(bson.M{"source": src, "series.id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if info.Removed == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func markSubscription(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	src, sid := p["source"], p["id"]

	r.ParseForm()
	if len(r.Form["type"]) < 1 || len(r.Form["state"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := r.Form["type"][0]
	s := r.Form["state"][0]

	id := src + "/" + sid

	if s != source.StateIgnored && s != source.StateUnseen && s != source.StateSeen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	var info *mgo.ChangeInfo
	if t == "series" {
		info, err = source.GetPersistence("subscriptions").UpdateAll(bson.M{"_id": bson.M{"$regex": `^` + id + `/\d+`}}, bson.M{"$set": bson.M{"state": s}})
	} else if t == "episode" {
		if len(r.Form["episode"]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		e := r.Form["episode"][0]
		id = id + "/" + e
		info, err = source.GetPersistence("subscriptions").UpdateAll(bson.M{"_id": id}, bson.M{"$set": bson.M{"state": s}})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if info.Updated == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
