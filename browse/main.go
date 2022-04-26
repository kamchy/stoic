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

func readTemplate(tmpl ...string) *template.Template {
	log.Printf("Reading from %v", tmpl)
	tp, e := template.ParseFiles(tmpl...)
	for idx, tt := range tp.Templates() {
		log.Infof("%d: Template name is %s", idx, tt.Name())
	}
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
				id, err := repo.SaveQuote(*q.Quote)
				if err != nil {
					q.Message["save"] = err.Error()
				} else {
					q.Id = id
				}
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

type AllQuotesData struct {
	Quotes []model.QuoteWithCount
	Len int
	Error  error
}

func allQuotesHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		qs, err := repo.ReadAllQuotes()
		log.Infof("In all quotes handler, found %d quotes", len(qs))
		template.Execute(writer, AllQuotesData{Quotes: qs, Len: len(qs), Error: err})
	}

}
func removeThoughtHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Infof("in removeThoughtHandler: method %s, header: %v, form values: %v", request.Method, request.Header,request.Form)
		if request.Method == http.MethodPost {
			ts := make([]model.Thought, 0)
			tId := request.FormValue("id")
			
			id, err := strconv.ParseInt(tId, 10, 64)
			q := model.Quote{}
			if err == nil {
				err := repo.RemoveThought(id)
				if err == nil {
					qId := request.FormValue("quoteid")
					qidnum, err := strconv.ParseInt(qId, 10, 64)
					if err == nil {
						ts, err = repo.ReadThoughtsForQuote(qidnum)
						q, err = repo.ReadQuote(qidnum)
					} else {
						ts, err = repo.ReadAllThoughts()
					}
				}
			}
			
			template.Execute(writer, AllThoughtsData{Thoughts: ts, Quote: &q, Error: err})
		}
	}

}

func removeQuoteHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		log.Infof("in removeQuoteHandler: method %s, header: %v", request.Method, request.Header)
		if request.Method == http.MethodPost {
			qs := make([]model.QuoteWithCount, 0)
			qId := request.FormValue("id")
			id, err := strconv.ParseInt(qId, 10, 64)
			if err == nil {
				err := repo.RemoveQuote(id)
				if err == nil {
					qs, err = repo.ReadAllQuotes()
				}
			}
			template.Execute(writer, AllQuotesData{qs, len(qs), err})
		}
	}
}

type AllThoughtsData struct {
	Thoughts []model.Thought
	Quote *model.Quote
	Error error
}

func thoughtsHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet: allThoughtsHandler(repo, template).ServeHTTP(writer, request)
		case http.MethodPost: removeThoughtHandler(repo, template).ServeHTTP(writer, request)
		}
	}
}

func allThoughtsHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	l := log.WithField("method", "allThoughtsHandler")
	return func(writer http.ResponseWriter, request *http.Request) {
		var ts []model.Thought
		var q model.Quote
		var err error
		vals := request.URL.Query()
		quoteidparam := vals.Get("quoteid")
		if quoteidparam != "" {
			qid, err := strconv.ParseInt(quoteidparam, 10, 64)
			l.Infof("all thoughts for id %d", qid)
			if err == nil {
				ts, err = repo.ReadThoughtsForQuote(qid)
				if err == nil {
					q, err = repo.ReadQuote(qid)
				}
			}
		} 
		if quoteidparam == "" || err != nil {
			l.Info("all thoughts - quoteid not given")
			ts, err = repo.ReadAllThoughts()
			template.Execute(writer, AllThoughtsData{Thoughts: ts, Error: err, Quote: nil})
		} else {
			l.Infof("Rendering %d thoughts for id %s", len(ts), quoteidparam)
			template.Execute(writer, AllThoughtsData{Thoughts: ts, Error: err, Quote: &q})
		}
	}

}

func thoughtHandler(repo stoic.Repository, template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		q := NewAddQuote(false)
		th := request.FormValue("thought")
		quoteid := request.FormValue("quoteid")
		id, err := strconv.ParseInt(quoteid, 10, 64)
		log.Infof("thought: %v for quoteid: %d, error: %v", th, id, err)
		if err == nil {
			th := model.Thought{
				Text: th, Time: time.Now(),
				QuoteId: id}
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
	t := readTemplate("index.html", "quotes.html", "thoughts.html")
	idxtemplate := t.Lookup("index.html")
	qlistTemplate := t.Lookup("quotes.html")
	thoughtTemplate := t.Lookup("thoughts.html")

	if len(os.Args) < 2 {
		log.Fatalf("Expected path to database file, got: %v", os.Args)
	}
	dbname := os.Args[1]
	repo, err := stoicdb.New(dbname)
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir("./static/"))
	http.HandleFunc("/", quoteHandler(repo, idxtemplate, false))
	http.HandleFunc("/addquote", quoteHandler(repo, idxtemplate, true))
	http.HandleFunc("/quotes", allQuotesHandler(repo, qlistTemplate))
	http.HandleFunc("/rmquote", removeQuoteHandler(repo, qlistTemplate))
	http.HandleFunc("/thoughts", thoughtsHandler(repo, thoughtTemplate))
	http.HandleFunc("/imgsvg", svgHandler())
	http.HandleFunc("/add", thoughtHandler(repo, idxtemplate))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	err = http.ListenAndServe(":5000", nil)
	log.Info("Listening on port 5000")
	if err != nil {
		log.Println(err)
	}
}
