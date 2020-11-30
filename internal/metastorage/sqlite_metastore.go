package metastorage

import (
	"os"
	"path/filepath"
	"database/sql"
	
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3" //SQLite goodness
	"github.com/prometheus/client_golang/api/prometheus/v1"

)
const createPromdexMetrics = `
	CREATE TABLE promdexMetrics (metric_name text NOT NULL PRIMARY KEY, metric_source_desc text)
`

// SQLiteMetastore is a Metastorer which uses SQLite.
type SQLiteMetastore struct {
	DBPath string //Full path to SQLite DB. Defaults to ./promdex.sqlite

	d *sql.DB
}

//NewSQLiteMetastore constructs a SQLiteMetastore with sensible defaults. If the provided path is an empty string, this function will default to ./promdex.sqlite. See struct documentation for versioned information. Init, if true, will drop the provided database and recreate it
func NewSQLiteMetastore(path string, init bool) (*SQLiteMetastore, error) {
	var s SQLiteMetastore
	var e error

	if path == "" {
		log.Info("Creating new datastore with default name")
		path = "./promdex.sql"
	}
	
	sanitizedPath, e := filepath.Abs(path)

	if e != nil {
		return &s, e
	}
	
	s.DBPath = sanitizedPath

	if init {
		if err := os.RemoveAll(s.DBPath); err != nil {
			return &s, err
		}
	}

	db, e := sql.Open("sqlite3", s.DBPath)
	defer db.Close() //TODO: Does this work?

	if e != nil {
		return &s, e
	}

	if e != nil {
		//TODO: Determine if this is safe enough? Bail out here?
		s.d = &sql.DB{}
	}

	s.d = db

	if init {
		if err := initDB(db); err != nil {
			log.Panic(err)
			panic(err)
		}
	}

	log.WithField("dbPath", s.DBPath).Info("sqlite metastore successfully (re-)initialized")
	return &s, e
}

//Store upserts a Prometheus metric into the backed SQLite database. This method will immediately error if the storer was set to version metrics without the appropriate call to SetMetricsVersion
func (s *SQLiteMetastore) Store(k string, m v1.Metadata) (error) {
	return nil
}

//TODO: TEST THIS!!! :zzzzzz:
func initDB(s *sql.DB) error {
	if e := s.Ping(); e != nil {
		return e
	}

	tx, e := s.Begin()
	log.Info("Beginning initDB transaction")
	if e != nil {
		return e
	}
	_, e = tx.Exec(createPromdexMetrics)
	if e != nil {
		return e
	}

	return tx.Commit()
}