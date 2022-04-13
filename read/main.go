package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/kamchy/stoic"
)

func handle(idx int, q stoic.Quote, f func(int, stoic.Quote) error) error {
	if e := f(idx, q); e != nil {
		return e
	}
	return nil
}

func printQuote(idx int, quote stoic.Quote) error {
	fmt.Printf("Quote %d\n%s\n%s\n\n",
		idx, quote.Text, quote.Author)
	return nil
}

func serial(fs ...func(int, stoic.Quote) error) func(int, stoic.Quote) error {
	return func(idx int, q stoic.Quote) error {
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
	var dbPath = flag.String("dbpath", stoic.DbName, fmt.Sprintf("Path to a database; default: %v", stoic.DbName))
	flag.Parse()
	var absPath = path.Clean(*dbPath)
	log.Printf("Path to a database as --dbname arg: %s", absPath)
	db, err := stoic.Open(absPath)
	if err != nil {
		log.Fatalf("Cannot open %v", dbPath)
	}
	if quotes, err := stoic.ReadQuotes(os.Stdin); err == nil {
		cnt, err := stoic.SaveQuotes(db, quotes)
		if err != nil {
			log.Fatalf("Cannot save quotes: %v", err.Error())
		}
		log.Printf("Saved %d rows in %s", cnt, *dbPath)
	}
}
