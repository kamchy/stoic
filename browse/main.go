package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type MMap = map[string]string
type AddQuote struct {
	*model.Quote
	Addquote bool
	Message  MMap
}

func NewAddQuote(shouldAdd bool) *AddQuote {
	return &AddQuote{Quote: new(model.Quote),
		Addquote: shouldAdd,
		Message:  make(MMap)}
}
func (q *AddQuote) isValid() bool {
	q.validate()
	log.Infof("isValid called for %v, message is %s", q.Quote, q.Message)
	return len(q.Message) == 0
}
func (q *AddQuote) validate() {
	if strings.TrimSpace(q.Text) == "" {
		q.Message["text"] = "Text cannot be empty"
	}

	if strings.TrimSpace(q.Author) == "" {
		q.Message["author"] = "Author cannot be empty"
	}
}

func quoteHandler(repo stoic.Repository, template *template.Template, shouldAdd bool) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		q := NewAddQuote(shouldAdd)
		if request.Method == http.MethodPost {
			q.Text = request.FormValue("quote")
			q.Author = request.FormValue("author")

			valid := q.isValid()
			log.Printf("Posted quote: %v is valid: %v message: %s", q.Quote, valid, q.Message)
			if valid {
				repo.SaveQuote(*q.Quote)
				q.Addquote = false
			}

		} else {
			if !shouldAdd {
				q.Quote = readQuote(repo)
			}
		}

		template.Execute(writer, q)

	}
}

func thoughtHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		q := NewAddQuote(false)
		th := request.FormValue("thought")
		quoteid := request.FormValue("quoteid")
		id, err := strconv.ParseInt(quoteid, 10, 32)
		log.Infof("thought: %v for quoteid: %d, error: %v", th, id, err)
		if err == nil {
			th := model.Thought{
				Text: th, Time: time.Now(),
				QuoteId: int(id)}
			saveThought(repo, th)
		}
		q.Quote = readQuote(repo)
		template.Execute(writer, q)
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
	fileServer := http.FileServer(http.Dir("./static/"))
	http.HandleFunc("/", quoteHandler(repo, t, false))
	http.HandleFunc("/addquote", quoteHandler(repo, t, true))
	http.HandleFunc("/imgsvg", svgHandler())
	http.HandleFunc("/add", thoughtHandler(repo, t))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	err = http.ListenAndServe(":5000", nil)
	log.Info("Listening on port 5000")
	if err != nil {
		log.Println(err)
	}
}
