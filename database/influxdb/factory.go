package influxdb

import "speedy/database"

func Factory(host string, dbName string, user string, pass string) (database.Database, error) {
	db, err := New(host, dbName, user, pass)
	if err != nil {
		return nil, err
	}
	return db, nil
}
