package stoic

import (
	"errors"
	"math/rand"
	"time"

	"github.com/kamchy/stoic/model"
	log "github.com/sirupsen/logrus"
)

// Interface Repository represents abstraction over storage implementations
// for saving and retrieving quotes and stoic thoughts

type Repository interface {
	// Reads all quotes from database
	ReadAllQuotes() ([]model.Quote, error)
	// Save single Thought
	SaveThought(model.Thought) (int64, error)
	// Save all Quotes
	SaveQuotes([]model.Quote) (int64, error)
}

// ReadRandomQuote reads and returns a random quote
func ReadRandomQuote(repo Repository) (*model.Quote, error) {

	quotes, err := repo.ReadAllQuotes()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Print("Read all quotes")
	var qlen = len(quotes)
	log.Printf("Found %d quotes", qlen)
	if qlen > 0 {
		rand.Seed(time.Now().UnixMilli())
		randQuote := quotes[rand.Intn(qlen)]
		return &randQuote, nil
	}
	return nil, errors.New("no new quote found")
}
