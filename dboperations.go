package stoic

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

// ReadRandomQuote reads and returns a random quote
func ReadRandomQuote(dbname string) (*Quote, error) {
	var uri = fmt.Sprintf("%s", dbname)
	db, err := Open(uri)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	quotes, err := ReadAllQuotes(db)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Print("Read all quotes")
	var qlen = len(quotes)
	log.Printf("Found %d quotes in %s db", qlen, uri)
	if qlen > 0 {
		rand.Seed(time.Now().UnixMilli())
		randQuote := quotes[rand.Intn(qlen)]
		return &randQuote, nil
	}
	return nil, errors.New("no new quote found")
}

// SaveUserThought saves a thought to given database returning
// last insert id or error
func SaveUserThought(dbname string, th Thought) (int64, error) {
	var uri = fmt.Sprintf("%s", dbname)
	db, err := Open(uri)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	res, err := SaveThought(db, th)
	log.Infof("Saving thought %v, result is %v", th, res)

	if err != nil {
		log.Errorf("Error executing %v with %v: %v",
			InsertThoughtStatement, th, err)
	}
	return res.LastInsertId()

}
