package source

import "time"

type Series struct {
	Id           float64
	Title        string
	Summary      string
	EpisodeCount float64
}

type Episode struct {
	Id       float64
	Season   float64
	Number   float64
	Title    string
	Airstamp time.Time
}
