package timescaledb

import (
	"fmt"
	"log"
	"strings"
	"database/sql"
	"github.com/melchor629/speedy/database"

	_ "github.com/lib/pq"
)

//Implementation of a database using timescaledb (postgresql).
type Database struct {
	client *sql.DB
	table string
}

//Creates a client to an timescaleDB (aka PostgreSQL with an extension). The address must be an URL with the following
//format: `postgres://USER:PASS@ADDRESS/DATABASE[?...extraParams]`. See https://godoc.org/github.com/lib/pq for
//the full format of the URL.
func New(addr string, table string) (*Database, error) {
	c, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	return &Database{
		c,
		table,
	}, nil
}

//Closes the connection to the database.
func (d *Database) Close() {
	d.client.Close()
}

//Store a list of entries in a batch. Supposes no error will occur. If so, the app will stop.
func (d *Database) Store(entries []database.Entry) {
	txn, err := d.client.Begin()
	if err != nil {
		log.Fatal(err)
	}

	//From https://stackoverflow.com/questions/21108084/golang-mysql-insert-multiple-data-at-once
	sqlStr := "INSERT INTO ?(time, mac, download, upload, ipv4, ipv6) VALUES"
	var valuesStr []string
	var values []interface{}

	for _, entry := range entries {
		valuesStr = append(valuesStr, "(NOW(), ?, ?, ?, ?, ?)")
		values = append(values, entry.Mac(), entry.GetDownloadSpeed(), entry.GetUploadSpeed(), entry.Ipv4(), entry.Ipv6())
	}

	sqlStr = fmt.Sprintf("%s\n%s", sqlStr, strings.Join(valuesStr, ","))
	stmt, _ := txn.Prepare(sqlStr)
	_, err = stmt.Exec(values)

	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}
	txn.Commit()
}