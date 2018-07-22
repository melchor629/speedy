package database

import "net"

//Entry with the information to store in the database.
type Entry interface {
	Ipv6() net.IP
	Ipv4() net.IP
	Mac() net.HardwareAddr
	GetDownloadSpeed() uint64
	GetUploadSpeed() uint64
}

//How a database should look like.
type Database interface {
	Store(entry []Entry)
	Close()
}
