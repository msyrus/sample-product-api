package pgsql

import (
	"database/sql"

	"github.com/msyrus/simple-product-inv/infra"
	"github.com/msyrus/simple-product-inv/log"
)

// DB is the postgres db
type DB struct {
	conn *sql.DB
	lgr  log.Logger
}

// NewDB returns a new postgres DB with conn
func NewDB(conn *sql.DB, lgr log.Logger) *DB {
	return &DB{
		conn: conn,
		lgr:  lgr,
	}
}

func (d *DB) println(stmt string, args ...interface{}) {
	if d.lgr != nil {
		d.lgr.Println(args...)
	}
}

// Exec executes a sql command
func (d *DB) Exec(stmt string, args ...interface{}) error {
	d.println(stmt, args...)
	_, err := d.conn.Exec(stmt, args...)
	return err
}

// Query executes a db query and return row
func (d *DB) Query(stmt string, args ...interface{}) (infra.Row, error) {
	d.println(stmt, args...)
	rows, err := d.conn.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	r := &Row{
		rows: rows,
	}
	return r, nil
}

// Row is an implementation of infra.Row
// it holds postgres query result rows
type Row struct {
	rows *sql.Rows
}

// Scan scans rows and store data into dest
func (r *Row) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

// Next checks if anymore data is available to scan from rows
func (r *Row) Next() bool {
	return r.rows.Next()
}

// Close closes row to scan
func (r *Row) Close() error {
	return r.rows.Close()
}
