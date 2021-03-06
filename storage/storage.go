package storage

import (
	"github.com/melchor629/speedy/capture"
	"github.com/melchor629/speedy/database"
	"log"
	"os"
	"sync"
	"time"
)

// The key is the MAC Address as String (to be easily hasheable in go I suppose)
type Storage struct {
	db map[string]Entry
	mutex sync.RWMutex
}

//Starts capturing the traffic, processing them and then storing it into the database every second. The recommended way
//is to call this function as a gorutine.
func (s *Storage) Start(capturer capture.Context, db database.Database) {
	s.db = make(map[string]Entry)
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

		s.mutex.Lock()
		elem, ok := s.db[packet.SrcMac.String()]
		if !ok {
			elem = Entry{ mac: packet.SrcMac }
		}

		changedMetadata := false
		if packet.IsIP4() {
			changedMetadata = !elem.ipv4.Equal(packet.SrcIp)
			elem.ipv4 = packet.SrcIp
		} else if packet.IsIP6() {
			changedMetadata = !elem.ipv6.Equal(packet.SrcIp)
			elem.ipv6 = packet.SrcIp
		}

		if reversed {
			elem.accumulatedDownload += uint64(packet.Bytes)
		} else {
			elem.accumulatedUpload += uint64(packet.Bytes)
		}

		elem.modified()
		s.db[packet.SrcMac.String()] = elem
		if changedMetadata {
			go s.storeChangeOfMetadata(db, elem)
		}
		s.mutex.Unlock()
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
			go db.Store(s.getCopyAndClearSpeed())
			s.cleanUpOldEntries()
		}
	}
}

func (s *Storage) storeChangeOfMetadata(db database.Database, entry Entry) {
	db.StoreMetadata(database.Entry(&entry))
}

//Cleanup: when some entry has not been modified for a while, it will be deleted
func (s *Storage) cleanUpOldEntries() {
	s.mutex.Lock()
	keysToDelete := make([]string, 0)
	for key, value := range s.db {
		if value.tooOld() {
			keysToDelete = append(keysToDelete, key)
		}
	}
	for _, key := range keysToDelete {
		delete(s.db, key)
	}
	s.mutex.Unlock()
}

func (s *Storage) getCopyAndClearSpeed() []database.Entry {
	s.mutex.RLock()
	newSlice := make([]database.Entry, 0)
	for key, value := range s.db {
		copiedValue := value
		newSlice = append(newSlice, database.Entry(&copiedValue))
		value.ClearSpeed()
		s.db[key] = value
	}
	s.mutex.RUnlock()
	return newSlice
}
