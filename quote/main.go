package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/kamchy/stoic"
	"github.com/kamchy/stoic/model"
	"github.com/kamchy/stoic/stoicdb"
	log "github.com/sirupsen/logrus"
)

func format(q *model.Quote) {
	color.Yellow(q.Text)
	color.Blue(q.Author)
}
func main() {
	dbpath := os.Args[1]
	log.SetLevel(log.WarnLevel)
	repo, err := stoicdb.New(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	if q, err := stoic.ReadRandomQuote(repo); err == nil {
		format(q)
	} else {
		log.Fatalf("Could not read: %v", err)
	}

}
