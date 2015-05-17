package badswede

import (
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
	s, err := self.findTournamentUrl(query.Tournament)
	if err != nil {
		return
	}
	tournamentUrl := s.AttrOr("href", MISSING)
	matchUrls, err := self.findMatchPages(tournamentUrl)
	if err != nil {
		return
	}

	tournament = &Tournament{s.Text(), make([]Match, 0)}
	for _, url := range matchUrls {
		matches, _ := self.findMatches(url)
		matches.Each(func(_ int, s *goquery.Selection) {
			if hasPlayer(s, query.Players) {
				match := Match{
					PlannedTime: strings.TrimSpace(s.Find(".plannedtime").Text()),
					Draw:        strings.TrimSpace(s.Find("td:nth-child(3)").Text()),
					Left:        strings.TrimSpace(s.Find("td:nth-child(4)").Text()),
					Right:       strings.TrimSpace(s.Find("td:nth-child(6)").Text()),
					Result:      strings.TrimSpace(s.Find("td:nth-child(7)").Text()),
				}
				tournament.Matches = append(tournament.Matches, match)
			}
		})
	}
	return
}

func (self *Scraper) findTournamentUrl(tournament string) (*goquery.Selection, error) {
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

func (self *Scraper) findMatchPages(url string) ([]string, error) {
	absoluteUrl, _ := self.browser.ResolveStringUrl(url)
	log.Println("tournamentUrl", absoluteUrl)
	err := self.browser.Open(absoluteUrl)
	if err != nil {
		return nil, err
	}
	log.Println(self.browser.Title())
	matchesLinkSelector := ".tournamentcalendar a"
	matches := self.browser.Dom().Find(matchesLinkSelector)
	urls := make([]string, 0)
	matches.Each(func(_ int, s *goquery.Selection) {
		var url, _ = self.browser.ResolveStringUrl(s.AttrOr("href", MISSING))
		urls = append(urls, url)
	})
	log.Println("urls", urls)
	return urls, nil
}

func (self *Scraper) findMatches(url string) (*goquery.Selection, error) {
	log.Println("matchesUrl", url)
	err := self.browser.Open(url)
	if err != nil {
		return nil, err
	}
	matchesRowSelector := ".matches tbody tr"
	selection := self.browser.Dom().Find(matchesRowSelector)
	return selection, nil
}

func hasPlayer(selection *goquery.Selection, players []string) bool {
	text := selection.Find("td:nth-child(4)").Text() + " " + selection.Find("td:nth-child(6)").Text()
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
