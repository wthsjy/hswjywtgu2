package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sync/atomic"
	"time"

	// for register mysql driver
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// ErrStmtNil prepared stmt error
	ErrStmtNil = errors.New("prepare failed and stmt nil")
	// ErrNoRows is returned by Scan when QueryRow doesn't return a row.
	// In such a case, QueryRow returns a placeholder *Row value that defers
	// this error until a Scan.
	ErrNoRows = sql.ErrNoRows
	// ErrTxDone transaction done.
	ErrTxDone = sql.ErrTxDone
)

// log err types
const (
	// MySQLDBErr for db err
	MySQLDBErr = "mysql-db-err"
	// MySQLTxErr for tx err
	MySQLTxErr = "mysql-tx-err"
	// MySQLRowErr for row err
	MySQLRowErr = "mysql-row-err"
	// MySQLRowsErr for rows err
	MySQLRowsErr = "mysql-rows-err"
	// MySQLStmtErr for stmt err
	MySQLStmtErr = "mysql-stmt-err"
)

// Config mysql config.
type Config struct {
	DSN          string        // data source name.Refer Link: https://github.com/go-sql-driver/mysql#dsn-data-source-name
	Active       int           // pool
	Idle         int           // pool
	IdleTimeout  time.Duration // connect max life time.
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
	TranTimeout  time.Duration // transaction sql timeout
}

// DB database connection
type DB struct {
	conf *Config
	conn *sql.DB
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (*DB, error) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transction timeout config")
	}

	d, err := connect(c)
	if err != nil {
		return nil, err
	}

	return &DB{
		conf: c,
		conn: d,
	}, nil
}

func connect(c *Config) (*sql.DB, error) {
	d, err := sql.Open("mysql", c.DSN)
	if err != nil {
		fmt.Printf("open mysql dsn:%s error:%v\n", c.DSN, err)
		return nil, err
	}
	d.SetMaxOpenConns(c.Active)
	d.SetMaxIdleConns(c.Idle)
	d.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	return d, nil
}

// Tx transaction.
type Tx struct {
	db     *DB
	tx     *sql.Tx
	c      context.Context
	cancel func()
}

// Row row.
type Row struct {
	err error
	*sql.Row
	db     *DB
	query  string
	args   []interface{}
	cancel func()
}

// Rows rows.
type Rows struct {
	*sql.Rows
	cancel func()
}

// Stmt prepared stmt.
type Stmt struct {
	db    *DB
	tx    bool
	query string
	stmt  atomic.Value
}

// Begin begin tx
func (db *DB) Begin(c context.Context) (tx *Tx, err error) {
	c, cancel := context.WithTimeout(c, time.Duration(db.conf.TranTimeout))
	rtx, err := db.conn.BeginTx(c, nil)
	if err != nil {
		fmt.Printf("BeginTx err:%v\n", err)
		cancel()
		return
	}
	tx = &Tx{db: db, tx: rtx, c: c, cancel: cancel}
	return
}

// Exec exec
func (db *DB) Exec(c context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	c, cancel := context.WithTimeout(c, time.Duration(db.conf.ExecTimeout))
	res, err = db.conn.ExecContext(c, query, args...)
	cancel()
	if err != nil {
		fmt.Printf("Exec err:%v\n", err)
	}
	return
}

// Ping for check mysql health
func (db *DB) Ping(c context.Context) (err error) {
	c, cancel := context.WithTimeout(c, time.Duration(db.conf.ExecTimeout))
	err = db.conn.PingContext(c)
	cancel()
	if err != nil {
		fmt.Printf("Ping err:%v\n", err)
	}
	return
}

// Prepare prepare
func (db *DB) Prepare(query string) (*Stmt, error) {
	stmt, err := db.conn.PrepareContext(context.Background(), query)
	if err != nil {
		fmt.Printf("Prepare err:%v\n", err)
		return nil, err
	}
	st := &Stmt{query: query, db: db}
	st.stmt.Store(stmt)
	return st, nil
}

// Prepared Prepared
func (db *DB) Prepared(query string) (stmt *Stmt) {
	stmt = &Stmt{query: query, db: db}
	s, err := db.conn.PrepareContext(context.Background(), query)
	if err == nil {
		stmt.stmt.Store(s)
		return
	}
	go func() {
		for {
			s, err := db.conn.PrepareContext(context.Background(), query)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			stmt.stmt.Store(s)
			return
		}
	}()
	return
}

// Query query
func (db *DB) Query(c context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	c, cancel := context.WithTimeout(c, time.Duration(db.conf.QueryTimeout))
	rs, err := db.conn.Query(query, args...)
	if err != nil {
		fmt.Printf("Query err:%v\n", err)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

// QueryRow QueryRow
func (db *DB) QueryRow(c context.Context, query string, args ...interface{}) *Row {
	c, cancel := context.WithTimeout(c, time.Duration(db.conf.QueryTimeout))
	r := db.conn.QueryRowContext(c, query, args...)
	return &Row{db: db, Row: r, query: query, args: args, cancel: cancel}
}

// Close Close.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Commit commits the transaction.
func (tx *Tx) Commit() (err error) {
	err = tx.tx.Commit()
	tx.cancel()
	if err != nil {
		fmt.Printf("Commit err:%v\n", err)
	}
	return
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() (err error) {
	err = tx.tx.Rollback()
	tx.cancel()
	if err != nil {
		fmt.Printf("Rollback err:%v\n", err)
	}
	return
}

// Exec executes a query that doesn't return rows. For example: an INSERT and UPDATE.
func (tx *Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	res, err = tx.tx.ExecContext(tx.c, query, args...)
	if err != nil {
		fmt.Printf("Exec err:%v\n", err)
	}
	return
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...interface{}) (rows *Rows, err error) {
	rs, err := tx.tx.QueryContext(tx.c, query, args...)
	if err == nil {
		rows = &Rows{Rows: rs}
	} else {
		fmt.Printf("Query, err:%v\n", err)
	}
	return
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	r := tx.tx.QueryRowContext(tx.c, query, args...)
	return &Row{Row: r, db: tx.db, query: query, args: args}
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	as, ok := stmt.stmt.Load().(*sql.Stmt)
	if !ok {
		return nil
	}
	ts := tx.tx.StmtContext(tx.c, as)
	st := &Stmt{query: stmt.query, tx: true, db: tx.db}
	st.stmt.Store(ts)
	return st
}

// Prepare creates a prepared statement for use within a transaction.
// The returned statement operates within the transaction and can no longer be
// used once the transaction has been committed or rolled back.
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		fmt.Printf("Prepare, err:%v\n", err)
		return nil, err
	}
	st := &Stmt{query: query, tx: true, db: tx.db}
	st.stmt.Store(stmt)
	return st, nil
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) (err error) {
	if r.err != nil {
		err = r.err
	} else if r.Row == nil {
		err = ErrStmtNil
	}
	if err != nil {
		fmt.Printf("Scan err:%v\n", err)
		return
	}
	err = r.Row.Scan(dest...)
	if r.cancel != nil {
		r.cancel()
	}
	if err != nil && err != ErrNoRows {
		fmt.Printf("Scan Not NO Rows err:%v\n", err)
	}
	return
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result sets,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (rs *Rows) Close() (err error) {
	err = rs.Rows.Close()
	if rs.cancel != nil {
		rs.cancel()
	}
	return
}

// Exec executes a prepared statement with the given arguments and returns a
// Result summarizing the effect of the statement.
func (s *Stmt) Exec(c context.Context, args ...interface{}) (res sql.Result, err error) {
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	c, cancel := context.WithTimeout(c, time.Duration(s.db.conf.ExecTimeout))
	res, err = stmt.ExecContext(c, args...)
	cancel()
	if err != nil {
		fmt.Printf("Exec err:%v\n", err)
	}
	return
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (s *Stmt) Query(c context.Context, args ...interface{}) (rows *Rows, err error) {
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	c, cancel := context.WithTimeout(c, time.Duration(s.db.conf.QueryTimeout))
	rs, err := stmt.QueryContext(c, args...)
	if err != nil {
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRow(c context.Context, args ...interface{}) (row *Row) {
	row = &Row{db: s.db, query: s.query, args: args}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		return
	}
	c, cancel := context.WithTimeout(c, time.Duration(s.db.conf.QueryTimeout))
	row.Row = stmt.QueryRowContext(c, args...)
	row.cancel = cancel
	return
}

// Close closes the statement.
func (s *Stmt) Close() (err error) {
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if ok {
		err = stmt.Close()
	}
	return
}
