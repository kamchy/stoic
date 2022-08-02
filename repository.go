package stoic

import (
	"github.com/kamchy/stoic/model"
	log "github.com/sirupsen/logrus"
)

// Repository represents abstraction over storage implementations
// for saving and retrieving quotes and stoic thoughts
type Repository interface {
	//ReadAllQuotes reads all quotes from database
	ReadAllQuotes() ([]model.QuoteWithCount, error)
	//ReadQuote Read quote by id
	ReadQuote(int64) (model.Quote, error)
	//SaveThought save single Thought
	SaveThought(model.Thought) (int64, error)
	//SaveQuotes saves all Quotes; returns number of added rows and error
	SaveQuotes([]model.Quote) (int64, error)
	// SaveQuote saves single quote, returns id and error
	SaveQuote(model.Quote) (int64, error)
	// RemoveQuote removes quote with given id, returning error
	RemoveQuote(int64) error
	//RemoveThought Removes thought with given id
	RemoveThought(int64) error
	//ReadAllThoughts reads all thoughts, returning slice og model.ThoughtWithQuote and error
	ReadAllThoughts() ([]model.ThoughtWithQuote, error)
	//ReadThoughtsForQuote Reads thoughts for given quote id
	ReadThoughtsForQuote(int64) ([]model.ThoughtWithQuote, error)
	// ThoughtsCountForQuote returns number of thougths recorded for given query
	ThoughtsCountForQuote(int64) (int64, error)
	//ReadRandomQuote reads random quote from db
	ReadRandomQuote() (model.Quote, error)
}

// ReadRandomQuote reads and returns a random quote
func ReadRandomQuote(repo Repository) (*model.Quote, error) {
	quote, err := repo.ReadRandomQuote()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &quote, nil
}
