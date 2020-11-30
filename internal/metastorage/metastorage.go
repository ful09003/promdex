package metastorage

import (
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

// PromdexStorageType is a custom type used to enumerate possible storage backends
type PromdexStorageType int

func (p PromdexStorageType) String() string {
	switch p {
	case NilStore:
		return "null"
	case SQLiteStore:
		return "sqlite"
	default:
		return ""
	}
}

const (
	//NilStore is a PromdexStorageType which does nothing
	NilStore PromdexStorageType = iota
	//SQLiteStore is a PromdexStorageType which uses a SQLite database
	SQLiteStore
)

//Metastorer represents behavior which allows Promdex to persist Prometheus metrics/metadata
type Metastorer interface {
	Store(k string, v v1.Metadata) error
	SetMetricsVersion(v string)
}

//NullMetastore is a sink that does nothing.
type NullMetastore struct {}

//Store as implemented here does nothing at all
func (n NullMetastore) Store(k string, v v1.Metadata) (error) {
	return nil
}

//SetMetricsVersion here does nothing at all
func (n NullMetastore) SetMetricsVersion(s string) {}