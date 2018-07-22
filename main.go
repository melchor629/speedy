//Utility to know how much internet consume people does
package main

import (
	"log"
	"os"
	"os/signal"
	"flag"
	"fmt"
	"./capture/pcap"
	"./storage"
	"./database"
	"./database/influxdb"
	"./database/timescaledb"
)

type dbImplFactoryFunction func (host string, dbName string, user string, pass string) (database.Database, error)
var dbImpl = map[string]dbImplFactoryFunction{
	"influxdb": influxdb.Factory,
	"timescaledb": timescaledb.Factory,
}

func main() {
	deviceArg := flag.String("device", "", "Selects the NIC where to listen to and grab statistics")
	dbHostArg := flag.String("db-url", "http://localhost:8086", "The URL to the database")
	dbUserArg := flag.String("db-user", "", "The username to the database, empty for nothing")
	dbPassArg := flag.String("db-pass", "", "The password to the database, empty for nothing")
	dbNameArg := flag.String("db-name", "speedy", "Name of the database")
	dbImplArg := flag.String("db", "influxdb", "Type of the db implementation")
	help := flag.String("help", "", "More help over a command")
	flag.Parse()

	nics, err := pcap.GetActiveInterfaces()

	if help != nil && *help != "" {
		switch *help {
		case "device":
			fmt.Println("Selects the network interface in which the utility will inspect to.")
			fmt.Println("Here you have a list of network interfaces:")
			for _, nic := range nics {
				fmt.Println("  -", nic)
			}
		case "db":
			fmt.Println("Selects an implementation of a database.")
			fmt.Println("Available implementations are:")
			for key := range dbImpl {
				fmt.Println("  -", key)
			}
		default:
			flag.PrintDefaults()
		}
		os.Exit(0)
	}

	if deviceArg == nil || *deviceArg == "" {
		fmt.Println("No device specified.")
		if err != nil {
			log.Fatalf("Could not retreive the network interfaces:\n%s", err)
		}

		fmt.Println("Here you have a list of network interfaces:")
		for _, nic := range nics {
			fmt.Println("  -", nic)
		}
		os.Exit(1)
	} else {
		found := false
		for _, nic := range nics {
			if *deviceArg == nic {
				found = true
			}
		}

		if !found {
			log.Fatal("The interface", *deviceArg, "doesn't exist or is inactive")
		}
	}

	//Creates the capturer using libpcap
	context, err := pcap.New(*deviceArg)
	if err != nil {
		log.Fatal(err)
	}
	defer context.Close() //Close the context, but only when we decide to end the main

	dbImplFactory, ok := dbImpl[*dbImplArg]
	if !ok {
		fmt.Println("Invalid database implementation:", *dbImplArg)
		fmt.Println("Available ones are:")
		for key := range dbImpl {
			fmt.Println("  -", key)
		}
		os.Exit(1)
	}

	db, err := dbImplFactory(*dbHostArg, *dbNameArg, *dbUserArg, *dbPassArg)
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}
	defer db.Close() //Same as before

	//Temporal storage
	mem := storage.Storage{}
	go mem.Start(context, db)

	//Wait for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c
	log.Println("Received SIGINT, closing...")

	//Now all defer statements will be executed :)
}
