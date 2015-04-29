package source

type Source interface {
	Name() string
	SearchSeries(title string) []Series
	GetSeries(id int) Series
	ListEpisodes(id int) []Episode
	GetPoller(id int, quit <-chan int) func(chan []Episode)
}

func GetSources() map[string]Source {
	return map[string]Source{
		"tvmaze":  NewTVMaze(),
		"anilist": NewAnilist(),
	}
}
