package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andersjanmyr/badswede"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var PLAYERS = []string{
	"Christian Möller",
	"Isabella Bergman",
	"Isak Sanick",
	"Johanna Sjöstedt",
	"Nils Ihse",
	"Rasmus Janmyr",
	"Tove Rasmusson",
}

func renderFunc(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Host)
	render := render.New(render.Options{Layout: "layout"})
	tournamentName := r.FormValue("tournament")
	log.Println("tournamentName", tournamentName)
	query := badswede.Query{tournamentName, PLAYERS}
	scraper := badswede.NewScraper()
	tournament, err := scraper.Scrape(query)
	if err != nil {
		sendError(w, err)
		return
	}
	render.HTML(w, http.StatusOK, "matches", tournament)
}

func init() {
	router := mux.NewRouter().StrictSlash(false)
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix("/").HandlerFunc(renderFunc)
	http.Handle("/", router)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Server listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
func sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Internal server error", err)
	log.Println("Internal server error", err)
}
