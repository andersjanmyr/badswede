package main

import (
	"fmt"

	"github.com/andersjanmyr/badswede"
)

func main() {
	query := badswede.Query{"Gothenburg Open 2015", []string{"Rasmus Janmyr", "Tove", "Simon"}}
	scraper := badswede.NewScraper()
	tournament, err := scraper.Scrape(query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", tournament)
}
