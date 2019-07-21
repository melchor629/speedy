package influxdb

import (
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/melchor629/speedy/database"
	"log"
	"time"
)

//Implementation of a database using influxdb.
type Database struct {
	client client.Client
	name string
}

//Creates a client to an influxdb with the given address and username/password combination and database name. The
//address must be something like "http://influxdb:8086". If the database has no username nor password, put empty strings
//in both.
func New(addr string, databaseName string, username string, password string) (*Database, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: addr,
		Username: username,
		Password: password,
		UserAgent: "speedy",
	})

	if err != nil {
		return nil, err
	}

	return &Database{
		c,
		databaseName,
	}, nil
}

//Closes the connection to the database.
func (d *Database) Close() {
	_ = d.client.Close()
}

//Store a list of entries in a batch. Supposes no error will occur. If so, the app will stop.
func (d *Database) Store(entries []database.Entry) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: d.name,
		Precision: "ns",
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, entry := range entries {
		tags := map[string]string{"mac": entry.Mac().String()}
		fields := map[string]interface{}{
			"download": int64(entry.GetDownloadSpeed()),
			"upload":   int64(entry.GetUploadSpeed()),
		}

		pt, err := client.NewPoint("measures", tags, fields, time.Now())

		if err != nil {
			log.Fatal(err)
			return
		}
		bp.AddPoint(pt)
	}

	err = d.client.Write(bp)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (d *Database) StoreMetadata(entry database.Entry) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: d.name,
		Precision: "ns",
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	tags := map[string]string{"mac": entry.Mac().String()}
	fields := map[string]interface{}{
		"ipv4": entry.Ipv4(),
		"ipv6": entry.Ipv6(),
	}

	pt, err := client.NewPoint("measures_metadata", tags, fields, time.Now())

	if err != nil {
		log.Fatal(err)
		return
	}
	bp.AddPoint(pt)

	err = d.client.Write(bp)
	if err != nil {
		log.Fatal(err)
		return
	}
}