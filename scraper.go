package badswede

import (
	"fmt"

	"github.com/headzoo/surf"
)

type Query struct {
	Tournament string
	Players    []string
}

type Tournament struct {
	Name string
}

func Scrape(query Query) (*Tournament, error) {
	bow := surf.NewBrowser()
	err := bow.Open("http://badmintonsweden.tournamentsoftware.com")
	if err != nil {
		panic(err)
	}

	fmt.Println(bow.Title())
	return nil, nil
}
