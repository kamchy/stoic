package stoicdb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/kamchy/stoic/model"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SqliteRepository struct {
	Name string
	Db   *sql.DB
}

// New creates SqliteRepository for given path
func New(dbpath string) (SqliteRepository, error) {
	db, err := open(dbpath)
	if err == nil {
		return SqliteRepository{dbpath, db}, nil
	}
	return SqliteRepository{}, err
}

//DbName is default database name
const DbName = "quotes.db"

// CreateQuotesStatement is sql statement that creates quote table
const CreateQuotesStatement = `
	create table if not exists quote (
		id integer primary key autoincrement,
		text varchar(1000) unique not null,
		author varchar(100) not null
);
`

// CreateThoughtStatement creates table thought with quoteid as foreign key to quote table
const CreateThoughtStatement = `
	create  table if not exists thought (
		id integer primary key autoincrement,
		text varchar(1024) unique not null,
		time ineger not null,
		quoteid integer not null,
		foreign key (quoteid) references quote(id)
	);
`

// ReadAllQuotesWithCountQuery reads all quotes
const ReadAllQuotesWithCountQuery = `select q.id, q.text, q.author, (select count(*) from thought t where t.quoteid=q.id) tcount from quote q order by tcount desc`

// ReadQuoteByIDQuery selects quote with given id
const ReadQuoteByIDQuery = `select id, text, author from quote where id=?;`

// ThouthsCountForQuoteQuery counts number of thougths for given queryid
const ThouthsCountForQuoteQuery = `select count(*) from thought t where t.quoteid=?`

// ReadAllThoughtsQuery reads all thoughts; see https://sqlite.org/lang_datefunc.html
const ReadAllThoughtsQuery = `select strftime("%s", time) tt, text, quoteid, id from thought order by tt desc;`

// ReadThoughtsForQuote reads all thoughts for given query id
const ReadThoughtsForQuote = `select strftime("%s", time) tt, text, quoteid, id from thought where quoteid=? order by tt desc;`

// InsertQuoteStatement is sql statement that inserts text and author a quote table row
const InsertQuoteStatement = `insert into quote(text, author) values(?, ?);`

// InsertThoughtStatement is sql statement that insterts a thought with time and quote id
const InsertThoughtStatement = `insert into thought(text, time, quoteid) values(?, ?, ?)`

// DeleteQuote is an SQL query to delete quote with given id
const DeleteQuote = `delete from quote where id=?`

// DeleteThought is an SQL query to delete thought with given id
const DeleteThought = `delete from thought where id=?`

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

	db, err := open(uri)
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

// open opens existing db, creating necessary tables
func open(uri string) (*sql.DB, error) {
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

// SaveQuote saves model.Quote to db and returns its db id (and error)
func (repo SqliteRepository) SaveQuote(q model.Quote) (lastInsertID int64, err error) {
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
		lastInsertID, err = res.LastInsertId()
		q.Id = lastInsertID
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

// SaveThought saves model.Thought and returns db id
func (repo SqliteRepository) SaveThought(th model.Thought) (lastInsertID int64, err error) {
	ps, err := repo.Db.Prepare(InsertThoughtStatement)

	if err != nil {
		log.Errorf("Error in SaveThought: %v", err)
		return -1, err
	}

	defer ps.Close()

	log.Infof("Saving %v", th)
	res, err := ps.Exec(th.Text, th.Time, th.QuoteId)
	if err == nil {
		lastInsertID, err = res.LastInsertId()
	}
	return
}

// ReadAllQuotes returns slice of model.QuoteWithCount
func (repo SqliteRepository) ReadAllQuotes() ([]model.QuoteWithCount, error) {
	l := log.WithField("method", "ReadAllQuotes")
	l.Print("Read started")
	quotes := make([]model.QuoteWithCount, 0)
	var err error
	rows, err := repo.Db.Query(ReadAllQuotesWithCountQuery)
	if err != nil {
		l.Printf("Error after executing %s: %v", ReadAllQuotesWithCountQuery, err)
		return quotes, err
	}
	defer rows.Close()
	var q model.QuoteWithCount
	for rows.Next() {
		if err := rows.Scan(&q.Quote.Id, &q.Quote.Text, &q.Quote.Author, &q.ThoughtCount); err == nil {
			quotes = append(quotes, q)
		}
	}
	l.Printf("Read has %v elems", len(quotes))
	return quotes, err

}

// ReadQuote returns model.Quote for given id
func (repo SqliteRepository) ReadQuote(id int64) (model.Quote, error) {
	l := log.WithField("method", "ReadQuote")
	l.Infof("Read quote for id %d", id)
	var q model.Quote
	stmt, err := repo.Db.Prepare(ReadQuoteByIDQuery)
	if err != nil {
		log.Printf("Error after executing %s: %v", ReadQuoteByIDQuery, err)
		return q, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)

	err = row.Scan(&q.Id, &q.Text, &q.Author)
	return q, err

}

// ThoughtsCountForQuote returns count of thoughts for given quote id
func (repo SqliteRepository) ThoughtsCountForQuote(id int64) (count int64, err error) {
	l := log.WithField("method", "ThoughtsCountForQuote")
	l.Infof("id %d", id)
	stmt, err := repo.Db.Prepare(ThouthsCountForQuoteQuery)
	if err != nil {
		log.Printf("Error after executing %s: %v", ThouthsCountForQuoteQuery, err)
		return -1, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)

	err = row.Scan(&count)
	return count, err

}

// RemoveQuote removes quote with given id
func (repo SqliteRepository) RemoveQuote(id int64) (err error) {
	l := log.WithField("id", id)
	l.Infof("removing quote %d", id)
	stmt, err := repo.Db.Prepare(DeleteQuote)
	var num int64 = 0
	if err != nil {
		return
	}
	res, err := stmt.Exec(id)
	if err == nil {
		num, err = res.RowsAffected()
		log.Infof("Removed %d rows", num)
	}
	log.WithField("id", nil)
	return
}

const DateTimeFormatGo = "2006-01-02 15:04:05"
const DateTimeFormatSqlite = "%Y-%m-%d %H:%M:%S"

type RowsFunc func(string) (*sql.Rows, error)

func (repo SqliteRepository) ReadAllThoughts() (ts []model.Thought, err error) {
	l := log.WithField("method", "readAllThoughts")
	rowsf := func(s string) (rs *sql.Rows, err error) {
		rs, err = repo.Db.Query(s)
		return
	}
	ts, err = readThoughts(repo, ReadAllThoughtsQuery, rowsf, l)
	l.Infof("Readall has %v elems", len(ts))
	return ts, err
}

func readThoughts(repo SqliteRepository, query string, rowsfn RowsFunc, l *log.Entry) (ts []model.Thought, err error) {
	querystr := fmt.Sprintf(query, DateTimeFormatSqlite)
	rows, err := rowsfn(querystr)
	if err != nil {
		l.Infof("Error after executing %s: %v", querystr, err)
		return ts, err
	}
	defer rows.Close()
	var q model.Thought
	var t string
	for rows.Next() {
		if err := rows.Scan(&t, &q.Text, &q.QuoteId, &q.Id); err == nil {
			if q.Time, err = time.Parse(DateTimeFormatGo, t); err != nil {
				l.Errorf("Cannot parse time %s as %s", t, time.RFC3339)
				q.Time = time.Now()
			}
			l.Info(fmt.Sprintf("Reading %+v", q))
			ts = append(ts, q)
		} else {
			l.Error(err)
		}
	}
	return
}

// ReadThoughtsForQuote reads and returns all thoughts recorded for given quote id
func (repo SqliteRepository) ReadThoughtsForQuote(qid int64) (ts []model.Thought, err error) {
	l := log.WithField("method", "readThoghtsForQuote")
	l.Infof("Read started for id %d", qid)

	rowsfn := func(s string) (rs *sql.Rows, err error) {
		if st, err := repo.Db.Prepare(s); err == nil {
			rs, err = st.Query(qid)
			defer st.Close()
		}
		return
	}
	ts, err = readThoughts(repo, ReadThoughtsForQuote, rowsfn, l)
	l.Infof("ReadThoughtsForQuote has %v thoughts", len(ts))
	return

}

// DeleteThought deletes a thought with given id
func (repo SqliteRepository) RemoveThought(id int64) (err error) {
	s, err := repo.Db.Prepare(DeleteThought)
	if err != nil {
		return err
	}
	_, err = s.Exec(id)
	return
}
