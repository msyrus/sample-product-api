package infra

// DB represents the DB infrastructure
type DB interface {
	Exec(stmt string, args ...interface{}) error
	Query(stmt string, args ...interface{}) (Row, error)
}

// Row represents a data row of DB
type Row interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}
