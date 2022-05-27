package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/kamchy/stoic"
	"github.com/kamchy/stoic/model"
	"github.com/kamchy/stoic/stoicdb"
)

func handle(idx int, q model.Quote, f func(int, model.Quote) error) error {
	if e := f(idx, q); e != nil {
		return e
	}
	return nil
}

func printQuote(idx int, quote model.Quote) error {
	fmt.Printf("Quote %d\n%s\n%s\n\n",
		idx, quote.Text, quote.Author)
	return nil
}

func serial(fs ...func(int, model.Quote) error) func(int, model.Quote) error {
	return func(idx int, q model.Quote) error {
		for _, f := range fs {
			if err := f(idx, q); err != nil {
				log.Fatal(err)
				return err
			}
		}
		return nil
	}
}

// Read all quotes from stdin and save to -dbpath (by default: stoic.DbName) sqlite3 database
func main() {
	log.SetLevel(log.WarnLevel)
	var dbPath = flag.String("dbpath", stoicdb.DbName, fmt.Sprintf("Path to a database; default: %v", stoicdb.DbName))
	flag.Parse()
	var absPath = path.Clean(*dbPath)
	log.Printf("Path to a database as --dbname arg: %s", absPath)

	repo, err := stoicdb.New(absPath)
	if err != nil {
		log.Fatal(err)
	}
	if quotes, err := stoic.ReadQuotes(os.Stdin); err == nil {
		cnt, err := repo.SaveQuotes(quotes)
		if err != nil {
			log.Fatalf("Cannot save quotes: %v", err.Error())
		}
		log.Printf("Saved %d rows in %s", cnt, *dbPath)
	}
}
