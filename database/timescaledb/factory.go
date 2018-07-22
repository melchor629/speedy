package timescaledb

import "github.com/melchor629/speedy/database"

func Factory(host string, table string, _ string, _ string) (database.Database, error) {
	db, err := New(host, table)
	if err != nil {
		return nil, err
	}
	return db, nil
}
