package stoicdb

import (
	"database/sql"
	"os"

	"github.com/kamchy/stoic/model"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SqliteRepository struct {
	Name string
	Db   *sql.DB
}

func New(dbpath string) (SqliteRepository, error) {
	db, err := Open(dbpath)
	if err == nil {
		return SqliteRepository{dbpath, db}, nil
	}
	return SqliteRepository{}, err
}

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

const CreateThoughtStatement = `
	create  table if not exists thought (
		id integer primary key autoincrement,
		text varchar(1024) unique not null,
		time ineger not null,
		quoteid integer not null,
		foreign key (quoteid) references quote(id)
	);
`

// ReadAllQuotesQuery reads all quotes
const ReadAllQuotesQuery = `select * from quote;`

// InsertQuoteStatement is sql statement that inserts text and author a quote table row
const InsertQuoteStatement = `insert into quote(text, author) values(?, ?);`

//
const InsertThoughtStatement = `insert into thought(text, time, quoteid) values(?, ?, ?)`

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
	for i, stmt := range []string{CreateQuotesStatement, CreateThoughtStatement} {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Errorf("createTables error for table no %d - %s: %v", i, stmt, err)
			return err
		}
	}
	return nil
}

// Open opens existing db, creating necessary tables
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

func (repo SqliteRepository) SaveQuote(q model.Quote)(lastInsertId int64, err error) {
	log.Infof("Saving %v", q)
	ps, err := repo.Db.Prepare(InsertQuoteStatement)

	if err != nil {
		log.Errorf("Error in SaveQuote: %v", err)
		return -1, err
	}

	defer ps.Close()

	log.Infof("Saving %v", q)
	res, err := ps.Exec(q.Author, q.Text)
	if err == nil {
		lastInsertId, err = res.LastInsertId()
		q.Id = int(lastInsertId)
		log.Infof("Saved as %v", q)
	}
	return 
	
}
// SaveQuotes saves to sb a slice pf Quote structs
func (repo SqliteRepository) SaveQuotes(qs []model.Quote) (int64, error) {
	log.Infof("SaveQuotes gives %d quotes for writing ", len(qs))
	ps, err := repo.Db.Prepare(InsertQuoteStatement)
	count := int64(0)
	if err != nil {
		return count, err
	}

	for _, q := range qs {
		log.Infof("Save iteration for %v", q)
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

func (repo SqliteRepository) SaveThought(th model.Thought) (lastInsertId int64, err error) {
	ps, err := repo.Db.Prepare(InsertThoughtStatement)

	if err != nil {
		log.Errorf("Error in SaveThought: %v", err)
		return -1, err
	}

	defer ps.Close()

	log.Infof("Saving %v", th)
	res, err := ps.Exec(th.Text, th.Time, th.QuoteId)
	if err == nil {
		lastInsertId, err = res.LastInsertId()
	}
	return
}

func (repo SqliteRepository) ReadAllQuotes() ([]model.Quote, error) {
	log.WithField("method", "read")
	log.Print("Read started")
	quotes := make([]model.Quote, 0)
	var err error
	rows, err := repo.Db.Query(ReadAllQuotesQuery)
	if err != nil {
		log.Printf("Error after executing %s: %v", ReadAllQuotesQuery, err)
		return quotes, err
	}
	defer rows.Close()
	var q model.Quote
	for rows.Next() {
		if err := rows.Scan(&q.Id, &q.Text, &q.Author); err == nil {
			quotes = append(quotes, q)
		}
	}
	log.Printf("Read has %v elems", len(quotes))
	return quotes, err

}

