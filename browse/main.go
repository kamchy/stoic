package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	_ "strings"
	"time"

	"github.com/kamchy/stoic"
	"github.com/kamchy/stoic/model"
	"github.com/kamchy/stoic/stoicdb"
	log "github.com/sirupsen/logrus"
)

func readTemplate(tmpl string) *template.Template {
	log.Printf("Reading from %v", tmpl)
	tp, e := template.ParseFiles(tmpl)
	if e == nil {
		return tp
	}
	log.Printf("After parsing temptae error is %v", e)
	nt, e := template.New("foo").Parse(`Check logs`)
	if e != nil {
		log.Error(e)
	}
	return nt
}

func quoteHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		q := readQuote(repo)
		template.Execute(writer, q)
	}
}

func thoughtHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		th := request.FormValue("thought")
		quoteid := request.FormValue("quoteid")
		id, err := strconv.ParseInt(quoteid, 10, 32)
		log.Infof("thought: %v for quoteid: %d, error: %v", th, quoteid, err)
		if err == nil {
			th := model.Thought{
				Text: th, Time: time.Now(),
				QuoteId: int(id)}
			saveThought(repo, th)
		}
		template.Execute(writer, readQuote(repo))
	}

}

func svgHandler() func(http.ResponseWriter, *http.Request) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/svg+xml")
		writer.WriteHeader(http.StatusOK)
		stoic.GenerateSvgImage(writer)
	}
	return handler

}

func saveThought(repo stoic.Repository, th model.Thought) {
	repo.SaveThought(th)
}


func readQuote(repo stoic.Repository) *model.Quote {
	randQuote, err := stoic.ReadRandomQuote(repo)
	if err != nil {
		return &model.Quote{Text: err.Error(), Author: "system"}
	}
	return randQuote
}

func main() {
	t := readTemplate("index.html")
	if len(os.Args) < 2 {
		log.Fatalf("Expected path to database file, got: %v", os.Args)
	}
	dbname := os.Args[1]
	repo, err := stoicdb.New(dbname)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", quoteHandler(repo, t))
	http.HandleFunc("/imgsvg", svgHandler())
	http.HandleFunc("/add", thoughtHandler(repo, t))
	err = http.ListenAndServe(":5000", nil)
	log.Info("Listening on port 5000")
	if err != nil {
		log.Println(err)
	}
}
