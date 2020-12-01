package internal

// PromdexStorageType is a custom type used to enumerate possible storage backends. Used for CLI purposes.
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
	StoreExporterMeta(e, k, v string) error
	AddMetricFlavor(e, k string, v PromdexFlavor) error
}

//NullMetastore is a sink that does nothing.
type NullMetastore struct{}

//StoreExporterMeta as implemented here does nothing at all
func (n NullMetastore) StoreExporterMeta(e, k, v string) error {
	return nil
}

//AddMetricFlavor as implemented here does nothing at all
func (n NullMetastore) AddMetricFlavor(e, k string, v PromdexFlavor) error {
	return nil
}
