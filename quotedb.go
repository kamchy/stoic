package stoic

import (
	"database/sql"
	
	"os"
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

// /DbName is default database name
const DbName = "quotes.db"


// CreateQuotesStatement is sql statement that creates quote table
const CreateQuotesStatement = `
	create table if not exists quote (
		id integer primary key autoincrement,
		text varchar(1000) unique not null,
		author varchar(100) not null
);
`
const ReadAllQuotesQuery = `select * from quote;`

// InsertQuoteStatement is sql statement that inserts text and author a quote table row
const InsertQuoteStatement = `insert into quote(text, author) values(?, ?);`

// Create creates db with given uri string
func Create(uri string) error {
	log.Printf("Removing %s", uri)
	os.Remove(uri)
	log.Printf("Creating %v...", uri)
	file, err := os.Create(uri) // Create SQLite file
	defer file.Close()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("%v created", uri)

	db, err := Open(uri)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	
	return nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(CreateQuotesStatement)
	return err
}

// Open opens existing db
func Open(uri string) (*sql.DB, error) {
	log.Printf("Opening db with uri %s", uri)
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		log.Fatalf("Open failed with %v", err.Error())
		return nil, err
	}
	err = createTables(db)
	if err != nil {
		log.Fatalf("Table creation failed: %s", err.Error())
	}
	return db, err
}

// Save saves to sb a slice pf Quote structs
func Save(db *sql.DB, qs []Quote) (int64, error) {
	ps, err := db.Prepare(InsertQuoteStatement)
	count := int64(0)
	if err != nil {
		return count, err
	}

	for _, q := range qs {
		res, err := ps.Exec(q.Text, q.Author)
		log.Printf("Saving %s %s", q.Text, q.Author)
		if err != nil {
			log.Printf("Save: eror %v while saving %v", err, q)
			continue
		}
		num, err := res.RowsAffected()
		if err != nil {
			log.Printf("Save: rows aff: %v", num)
			return count, err
		}
		count += num

	}
	return count, nil
}

func Read(db *sql.DB) ([]Quote, error) {
	log.WithField("method", "read")
	log.Print("Read started")
	quotes := make([]Quote, 0)
	var err error
	rows, err := db.Query(ReadAllQuotesQuery); 
	if err != nil {
		log.Printf("Error after executing %s: %v", ReadAllQuotesQuery, err)
		return quotes, err
	}
	defer rows.Close()
	var q Quote
	var id = 0 
	for rows.Next() {
		if err := rows.Scan(&id, &q.Text, &q.Author); err == nil {
			quotes = append(quotes, q)
		} 	
	}
	log.Printf("Read has %v elems", len(quotes))
	return quotes, err
	
}