package storage

import (
	"net"
	"time"
)

//An entry of data.
type Entry struct {
	mac net.HardwareAddr
	ipv4 net.IP
	ipv6 net.IP

	accumulatedDownload uint64
	accumulatedUpload uint64

	lastModified time.Time
}

//Get the IPv6 address for this entry (if given).
func (e *Entry) Ipv6() net.IP {
	return e.ipv6
}

//Get the IPv4 address for this entry (if given).
func (e *Entry) Ipv4() net.IP {
	return e.ipv4
}

//Get the MAC address for this entry.
func (e *Entry) Mac() net.HardwareAddr {
	return e.mac
}

//Gets the download speed for this entry (or the accumulated download)
func (e *Entry) GetDownloadSpeed() uint64 {
	return e.accumulatedDownload
}

//Gets the upload speed for this entry (or the accumulated upload)
func (e *Entry) GetUploadSpeed() uint64 {
	return e.accumulatedUpload
}

//Clear the accumulated upload and download speeds
func (e *Entry) ClearSpeed() {
	e.accumulatedUpload = 0
	e.accumulatedDownload = 0
}

func (e *Entry) tooOld() bool {
	return time.Now().Sub(e.lastModified) > time.Hour
}

func (e *Entry) modified() {
	e.lastModified = time.Now()
}
