package timescaledb

import (
	"database/sql"
	"fmt"
	"github.com/melchor629/speedy/database"
	"log"

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
	_ = d.client.Close()
}

//Store a list of entries in a batch. Supposes no error will occur. If so, the app will stop.
func (d *Database) Store(entries []database.Entry) {
	if len(entries) == 0 {
		return
	}

	txn, err := d.client.Begin()
	if err != nil {
		log.Fatal(err)
	}

	//From https://stackoverflow.com/questions/21108084/golang-mysql-insert-multiple-data-at-once
	sqlStr := fmt.Sprintf("INSERT INTO %s(time, mac, download, upload) VALUES\n" +
		"(NOW(), $1, $2, $3);", d.table)

	stmt, _ := txn.Prepare(sqlStr)

	for _, entry := range entries {
		_, err = stmt.Exec(
			toString(entry.Mac()),
			entry.GetDownloadSpeed(),
			entry.GetUploadSpeed(),
		)

		if err != nil {
			stmt.Close()
			txn.Rollback()
			log.Fatal(err)
		}
	}

	txn.Commit()
	stmt.Close()
}

func (d *Database) StoreMetadata(entry database.Entry) {
	sqlStr2 := fmt.Sprintf("INSERT INTO %s_metadata(mac, ipv4, ipv6) VALUES ($1, $2, $3)\n" +
		"ON CONFLICT (mac) DO\n" +
		"UPDATE SET ipv4 = $2, ipv6 = $3", d.table)
	stmt, err := d.client.Prepare(sqlStr2)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		toString(entry.Mac()),
		toString(entry.Ipv4()),
		toString(entry.Ipv6()),
	)

	if err != nil {
		log.Fatal(err)
	}
}

//Converts an object with .String() method into a NullString for database
func toString(a interface{ String() string }) sql.NullString {
	if a == nil {
		return sql.NullString{ Valid: false }
	}

	str := a.String()
	if str == "" || str == "<nil>" {
		return sql.NullString{ Valid: false }
	}

	return sql.NullString{
		String: str,
		Valid:  true,
	}
}
