package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	_ "strings"
	"time"

	"github.com/kamchy/stoic"
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

func quoteHandler(dbname string, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		q := readQuote(dbname)
		template.Execute(writer, q)
	}
}

func thoughtHandler(dbname string, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		th := request.FormValue("thought")
		quoteid := request.FormValue("quoteid")
		id, err := strconv.ParseInt(quoteid, 10, 32)
		log.Infof("thought: %v for quoteid: %d, error: %v", th, quoteid, err)
		if err == nil {
			th := stoic.Thought{
				Text: th, Time: time.Now(),
				QuoteId: int(id)}
			saveThought(dbname, th)
		}
		template.Execute(writer, readQuote(dbname))
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

func saveThought(dbname string, th stoic.Thought) {
	stoic.SaveUserThought(dbname, th)
}
func readQuote(dbname string) *stoic.Quote {
	// var r io.Reader = strings.NewReader("some io.Reader stream to be read\n")
	var randQuote *stoic.Quote
	randQuote, err := stoic.ReadRandomQuote(dbname)
	if err != nil {
		randQuote.Text = err.Error()
		randQuote.Author = "whoa!"
	}
	return randQuote
}

func main() {
	t := readTemplate("index.tmpl")
	dbname := os.Args[1]
	http.HandleFunc("/", quoteHandler(dbname, t))
	http.HandleFunc("/imgsvg", svgHandler())
	http.HandleFunc("/add", thoughtHandler(dbname, t))
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Println(err)
	}
}
