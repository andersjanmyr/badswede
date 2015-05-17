package badswede

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type Query struct {
	Tournament string
	Players    []string
}

type Tournament struct {
	Name string
}

func Scrape(query Query) (*Tournament, error) {
	log.Println(query)
	browser := surf.NewBrowser()
	result, err := findTournaments(browser, query.Tournament)
	if err != nil {
		return nil, err
	}
	result.Each(func(_ int, s *goquery.Selection) {
		fmt.Println(s.Attr("href"))
	})
	return nil, nil
}

func findTournaments(browser *browser.Browser, tournament string) (*goquery.Selection, error) {
	err := browser.Open("http://badmintonsweden.tournamentsoftware.com")
	if err != nil {
		panic(err)
	}
	form, err := browser.Form("#formBasePage")
	if err != nil {
		return nil, err
	}
	form.Input("tbxSearchQuery", tournament)
	form.Input("ctl00$ctl01$cphPage$ddlSearchType", "1")
	if err := form.Submit(); err != nil {
		return nil, err
	}
	fmt.Println(browser.Title())
	tournamentLinkSelector := "#cphPage_cphPage_tournamentlistpage_maincolumn_ctl06_ctl00_row1 h3 a"
	selection := browser.Dom().Find(tournamentLinkSelector)
	return selection, nil
}
