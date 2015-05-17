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

type Scraper struct {
	browser *browser.Browser
}

func NewScraper() *Scraper {
	scraper := Scraper{surf.NewBrowser()}
	return &scraper
}

const MISSING = "missing"

func (self *Scraper) Scrape(query Query) (*Tournament, error) {
	log.Println(query)
	result, err := self.findTournaments(query.Tournament)
	if err != nil {
		return nil, err
	}
	result.Each(func(_ int, s *goquery.Selection) {
		matches, _ := self.findMatches(s.AttrOr("href", MISSING))
		fmt.Println(matches.Text())
	})
	return nil, nil
}

func (self *Scraper) findTournaments(tournament string) (*goquery.Selection, error) {
	err := self.browser.Open("http://badmintonsweden.tournamentsoftware.com")
	if err != nil {
		return nil, err
	}
	form, err := self.browser.Form("#formBasePage")
	if err != nil {
		return nil, err
	}
	form.Input("tbxSearchQuery", tournament)
	form.Input("ctl00$ctl01$cphPage$ddlSearchType", "1")
	if err := form.Submit(); err != nil {
		return nil, err
	}
	log.Println(self.browser.Title())
	tournamentLinkSelector := "#cphPage_cphPage_tournamentlistpage_maincolumn_ctl06_ctl00_row1 h3 a"
	selection := self.browser.Dom().Find(tournamentLinkSelector)
	return selection, nil
}

func (self *Scraper) findMatches(url string) (*goquery.Selection, error) {
	absoluteUrl, err := self.browser.ResolveStringUrl(url)
	log.Println("tournamentUrl", absoluteUrl)
	err = self.browser.Open(absoluteUrl)
	if err != nil {
		return nil, err
	}
	log.Println(self.browser.Title())
	matchesLinkSelector := "#cphPage_cphPage_tmTournamentMenu li:nth-child(5) a"
	matchesHref := self.browser.Dom().Find(matchesLinkSelector).AttrOr("href", MISSING)
	absoluteUrl, err = self.browser.ResolveStringUrl(matchesHref)
	log.Println("matchesUrl", absoluteUrl)
	err = self.browser.Open(absoluteUrl)
	if err != nil {
		return nil, err
	}
	matchesRowSelector := ".matches tbody tr"
	selection := self.browser.Dom().Find(matchesRowSelector)
	log.Println("selection", selection.Text())
	return selection, nil
}
