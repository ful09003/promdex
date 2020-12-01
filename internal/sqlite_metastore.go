package internal

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" //SQLite goodness
	log "github.com/sirupsen/logrus"
)

const createPromdexMetrics = `
	CREATE TABLE promdexMetrics (metric_name text NOT NULL,
		metric_instance_job text,
		metric_source_desc text,
		PRIMARY KEY(metric_name,metric_instance_job))
`

const createPromdexDescriptions = `
	CREATE TABLE promdexDescriptions (metric_name text NOT NULL,
		metric_instance_job text,
		createDate text,
		updateDate text,
		desc text)
`

const createPromdexDescriptionsIndices = `
	CREATE INDEX IF NOT EXISTS idx_desc ON promdexDescriptions(metric_name, metric_instance_job)
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
		path = "./promdex.sql"
	}

	log.WithField("storePath", path).Debug("targeting storage path for sqlite metastore")

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
	//defer db.Close() //TODO: Does this work?

	if e != nil {
		return &s, e
	}

	s.d = db

	if init {
		if err := initDB(db); err != nil {
			log.Panic(err)
		}
	}

	log.WithField("dbPath", s.DBPath).Info("sqlite metastore successfully (re-)initialized")
	return &s, e
}

//StoreExporterMeta inserts a Prometheus metric into the backed SQLite database. This method will immediately error if the storer was set to version metrics without the appropriate call to SetMetricsVersion
func (s *SQLiteMetastore) StoreExporterMeta(ex, k, m string) error {
	if e := s.d.Ping(); e != nil {
		return e
	}

	tx, e := s.d.Begin()
	defer tx.Commit()
	log.Debug("beginning store transaction")
	_, e = tx.Exec("INSERT OR IGNORE INTO promdexMetrics(metric_name, metric_instance_job, metric_source_desc) VALUES(?,?,?)", k, ex, m)

	if e != nil {
		return e
	}

	return nil
}

//AddMetricFlavor inserts a row into an SQLiteMetastore's promdexDescriptions table. Takes an exporter name (e.g. job name), a metric name, and an array of bytes which are JSON marshalled and stored using sqlite/json1.
func (s *SQLiteMetastore) AddMetricFlavor(ex, k string, v PromdexFlavor) error {
	if e := s.d.Ping(); e != nil {
		return e
	}

	tx, e := s.d.Begin()
	defer tx.Commit()
	log.Debug("beginning flavorful transaction")

	timeNow := time.Now().UTC().String()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, e = tx.Exec("INSERT INTO promdexDescriptions(metric_name, metric_instance_job, createDate, desc) VALUES(?,?,?,?)", k, ex, timeNow, string(b))

	if e != nil {
		return e
	}

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
	_, e = tx.Exec(createPromdexDescriptions)
	if e != nil {
		return e
	}
	_, e = tx.Exec(createPromdexDescriptionsIndices)
	if e != nil {
		return e
	}

	return tx.Commit()
}
