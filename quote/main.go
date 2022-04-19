package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/kamchy/stoic"
	log "github.com/sirupsen/logrus"
)

func format(q *stoic.Quote) {
	color.Yellow(q.Text)
	color.Blue(q.Author)
}
func main() {
	dbpath := os.Args[1]
	log.SetLevel(log.WarnLevel)
	if q, err := stoic.ReadRandomQuote(dbpath); err == nil {
		format(q)
	} else {
		log.Fatalf("Could not read: %v", err)
	}

}
