package storage

import (
	"speedy/capture"
	"speedy/database"
	"time"
	"log"
	"os"
)

// The key is the MAC Address as String (to be easily hasheable in go I suppose)
type Storage map[string]Entry

//Starts capturing the traffic, processing them and then storing it into the database every second. The recommended way
//is to call this function as a gorutine.
func (s Storage) Start(capturer capture.Context, db database.Database) {
	stop := make(chan bool)
	go capturer.StartCapturing()
	go s.storeInDB(db, stop)

	for packet := range capturer.Packets() {
		reversed := false
		if packet.IsReversed(capturer.GetMAC()) {
			packet.Reverse()
			reversed = true
		}

		if packet.IsBroadcast() || packet.IsIPv6Multicast() {
			continue
		}

		elem, ok := s[packet.SrcMac.String()]
		if !ok {
			elem = Entry{ mac: packet.SrcMac }
		}

		if packet.IsIP4() {
			elem.ipv4 = packet.SrcIp
		} else if packet.IsIP6() {
			elem.ipv6 = packet.SrcIp
		}

		if reversed {
			elem.accumulatedDownload += uint64(packet.Bytes)
		} else {
			elem.accumulatedUpload += uint64(packet.Bytes)
		}

		elem.modified()
		s[packet.SrcMac.String()] = elem
	}

	stop <- true
}

//Every second, gets a copy of the memory db and stores them into the good old db. Also cleans the unused entries.
func (s *Storage) storeInDB(db database.Database, stop chan bool) {
	logger := log.New(os.Stdout, "[Storage]: ", 0)
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()
	logger.Println("Starting storeInDB gorutine")
	itsTimeToStop := false
	for !itsTimeToStop {
		select {
		case <- stop:
			logger.Println("Stopping storeInDB gorutine...")
			itsTimeToStop = true
		case <- timer.C:
			db.Store(s.getCopyAndClearSpeed())

			//Cleanup: when some entry has not been modified for a while, it will be deleted
			for key, value := range *s {
				if value.tooOld() {
					delete(*s, key)
				}
			}
		}
	}
}

func (s Storage) getCopyAndClearSpeed() []database.Entry {
	newSlice := make([]database.Entry, 0)
	for key, value := range s {
		copiedValue := value
		newSlice = append(newSlice, database.Entry(&copiedValue))
		value.ClearSpeed()
		s[key] = value
	}
	return newSlice
}
