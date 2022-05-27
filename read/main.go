package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/kamchy/stoic"
	"github.com/kamchy/stoic/stoicdb"
)

// OptionPath is the name of the option for database path name, default is stoicdb.DbName
const OptionPath = "dbpath"

// Read all quotes from stdin and save to -dbpath (by default: stoic.DbName) sqlite3 database
func main() {
	log.SetLevel(log.WarnLevel)

	var dbPath = flag.String(OptionPath, stoicdb.DbName, fmt.Sprintf("Path to a database: %v", stoicdb.DbName))
	flag.Parse()

	var absPath = path.Clean(*dbPath)
	log.Printf("Path to a database passed  as -%s arg: %s", OptionPath, absPath)

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
