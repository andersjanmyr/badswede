package badswede

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type Query struct {
	Tournament string
	Players    []string
}

type Match struct {
	PlannedTime string
	Draw        string
	Left        string
	Right       string
	Result      string
}

type Tournament struct {
	Name    string
	Matches []Match
}

type Scraper struct {
	browser *browser.Browser
}

func NewScraper() *Scraper {
	scraper := Scraper{surf.NewBrowser()}
	return &scraper
}

const MISSING = "missing"

func (self *Scraper) Scrape(query Query) (tournament *Tournament, err error) {
	log.Println(query)
	tournament = nil
	result, err := self.findTournaments(query.Tournament)
	if err != nil {
		return
	}
	result.Each(func(_ int, s *goquery.Selection) {
		tournament = &Tournament{s.Text(), make([]Match, 0)}
		matches, _ := self.findMatches(s.AttrOr("href", MISSING))
		matches.Each(func(_ int, s *goquery.Selection) {
			if hasPlayer(s, query.Players) {
				match := Match{
					PlannedTime: s.Find(".plannedtime").Text(),
					Draw:        s.Find("td:nth-child(3)").Text(),
					Left:        s.Find("td:nth-child(4)").Text(),
					Right:       s.Find("td:nth-child(6)").Text(),
					Result:      s.Find("td:nth-child(7)").Text(),
				}
				html, _ := s.Html()
				fmt.Println("html", html)
				tournament.Matches = append(tournament.Matches, match)
			}
		})
	})
	return
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

func hasPlayer(selection *goquery.Selection, players []string) bool {
	text := selection.Find("td:nth-child(4)").Text() + " " + selection.Find("td:nth-child(6)").Text()
	fmt.Println("text", text, players)
	if len(strings.TrimSpace(text)) == 0 {
		return false
	}
	for _, player := range players {
		if strings.Contains(text, player) {
			return true
		}
	}
	return false
}
