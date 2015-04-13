package source

import "time"

type Dummy struct {
}

func NewDummy() Source {
	return Dummy{}
}

func (Dummy) Name() string {
	return "dummy"
}

func (Dummy) SearchSeries(title string) []Series {
	return []Series{
		Series{
			Id:    1,
			Title: "Hello, world!",
		},
		Series{
			Id:    2,
			Title: "Agence Tous Risques",
		},
	}
}

func (Dummy) GetSeries(id int) Series {
	return Series{
		Id:    1,
		Title: "Hello, world!",
	}
}

func (Dummy) ListEpisodes(id int) []Episode {
	return []Episode{
		Episode{
			Id:     1,
			Season: 1,
			Number: 1,
			Title:  "First episode",
		},
		Episode{
			Id:     2,
			Season: 1,
			Number: 2,
			Title:  "Second episode",
		},
		Episode{
			Id:     3,
			Season: 1,
			Number: 3,
			Title:  "Third episode",
		},
		Episode{
			Id:     4,
			Season: 1,
			Number: 4,
			Title:  "Fourth episode",
		},
	}
}

func (s Dummy) GetPoller(id int, quit <-chan int) func(chan []Episode) {
	return func(c chan []Episode) {
	Poller:
		for {
			select {
			case <-quit:
				break Poller
			case <-time.After(5 * time.Second):
				c <- s.ListEpisodes(id)
			}
		}
	}
}
