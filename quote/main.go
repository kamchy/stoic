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
	if q, err := stoic.ReadRandomQuote(os.Args[1]); err == nil {
		format(q)
	} else {
		log.Fatalf("Could not read: %v", err)
	}

}
