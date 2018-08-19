package storage

import (
	"testing"
	"time"
	"github.com/melchor629/speedy/database"
	"github.com/melchor629/speedy/capture"
	"net"
)

//TESTS FOR: getCopyAndClearSpeed

func TestGetCopyAndClearSpeedWithoutEntries(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }

	l := s.getCopyAndClearSpeed()

	if len(l) != 0 {
		t.Error("List was expected to be empty, instead received:", l)
	}
}

func TestGetCopyAndClearSpeedWithEntries(t *testing.T) {
	s := Storage{
		db: map[string]Entry{
			"00:11:22:33:44:55": {
				mac: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				accumulatedUpload: 123,
				accumulatedDownload: 456,
			},
			"aa:bb:cc:dd:ee:ff": {
				mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
				accumulatedUpload: 789,
				accumulatedDownload: 150,
			},
		},
	}

	l := s.getCopyAndClearSpeed()

	if len(l) != 2 {
		t.Error("List was expected to have 2 elements, instead received:", l)
	}

	for _, entry := range l {
		entry2 := s.db[entry.Mac().String()]
		if entry.Mac().String() != entry2.Mac().String() {
			t.Error("Macs not match:", entry.Mac(), entry2.Mac())
		}
		if entry.GetDownloadSpeed() == 0 {
			t.Error("Copied entry has 0 download speed")
		}
		if entry.GetUploadSpeed() == 0 {
			t.Error("Copied entry has 0 upload speed")
		}
		if entry2.GetDownloadSpeed() != 0 {
			t.Error("Original entry has not 0 download speed")
		}
		if entry2.GetUploadSpeed() != 0 {
			t.Error("Original entry has not 0 upload speed")
		}
	}
}

//TESTS FOR: cleanUpOldEntries

func TestCleanUpOldEntriesWithoutEntries(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }

	s.cleanUpOldEntries()

	if len(s.db) != 0 {
		t.Error("Entries appeared magically")
	}
}

func TestCleanUpOldEntriesWithEntriesButNoOneOld(t *testing.T) {
	s := Storage{
		db: map[string]Entry{
			"00:11:22:33:44:55": {
				mac: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				accumulatedUpload: 123,
				accumulatedDownload: 456,
				lastModified: time.Now(),
			},
			"aa:bb:cc:dd:ee:ff": {
				mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
				accumulatedUpload: 789,
				accumulatedDownload: 150,
				lastModified: time.Now().Add(-1000000),
			},
		},
	}

	s.cleanUpOldEntries()

	if len(s.db) != 2 {
		t.Error("Some entry disappeared or appeared magically")
	}
}

func TestCleanUpOldEntriesWithEntriesAndOneOld(t *testing.T) {
	s := Storage{
		db: map[string]Entry{
			"00:11:22:33:44:55": {
				mac: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				accumulatedUpload: 123,
				accumulatedDownload: 456,
				lastModified: time.Now(),
			},
			"aa:bb:cc:dd:ee:ff": {
				mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
				accumulatedUpload: 789,
				accumulatedDownload: 150,
				lastModified: time.Now().Add(-time.Hour * 10),
			},
		},
	}

	s.cleanUpOldEntries()

	if len(s.db) != 1 {
		t.Error("Some entry didn't disappeared or appeared magically")
	}
}

//TESTS FOR: storeInDB

type dumbDB struct {
	storeCalled bool
	entries []database.Entry
}

func (db *dumbDB) Store(entry2 []database.Entry) {
	db.storeCalled = true
	db.entries = entry2
}

func (db *dumbDB) Close() {}

func TestStoreInDb(t *testing.T) {
	s := Storage{
		db: map[string]Entry{
			"00:11:22:33:44:55": {
				mac: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				accumulatedUpload: 123,
				accumulatedDownload: 456,
				lastModified: time.Now(),
			},
			"aa:bb:cc:dd:ee:ff": {
				mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
				accumulatedUpload: 789,
				accumulatedDownload: 150,
				lastModified: time.Now().Add(-time.Hour * 10),
			},
		},
	}
	db := dumbDB{}
	stop := make(chan bool)

	go s.storeInDB(&db, stop)
	<- time.NewTimer(1 * time.Second + 500 * time.Millisecond).C
	stop <- true

	if !db.storeCalled {
		t.Error("Database was not called")
		t.FailNow()
	}

	if len(db.entries) != 2 {
		t.Error("Entries seems to not to be correct")
		t.FailNow()
	}

	dbMac1 := db.entries[0].Mac().String()
	dbMac2 := db.entries[1].Mac().String()
	if dbMac1 != "00:11:22:33:44:55" && dbMac2 != "00:11:22:33:44:55" {
		t.Error("00:11:22:33:44:55 not found in entries,", dbMac1, dbMac2)
	}

	if dbMac1 != "aa:bb:cc:dd:ee:ff" && dbMac2 != "aa:bb:cc:dd:ee:ff" {
		t.Error("00:11:22:33:44:55 not found in entries,", dbMac1, dbMac2)
	}
}

//TESTS FOR: Start

type dumbCapturer struct {
	p chan *capture.Packet
}

func (c *dumbCapturer) GetMAC() net.HardwareAddr { return []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff} }
func (c *dumbCapturer) Packets() chan *capture.Packet { return c.p }
func (c *dumbCapturer) StartCapturing() {}
func (c *dumbCapturer) Close() { close(c.p) }

func TestStartNoPackets(t *testing.T) {
	s := Storage{}
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.Close()

	if d.storeCalled {
		t.Error("Something appeared in DB from nowhere")
	}
}

func TestStartPacketNotReversedIPv4(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		SrcMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		DstMac: c.GetMAC(),
		IpType: 4,
		SrcIp: []byte{ 127, 0, 0, 1 },
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedUpload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedUpload)
	}

	if e.accumulatedDownload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedDownload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv4.String() != "127.0.0.1" {
		t.Error("IPv4 should be 127.0.0.1, but is", e.ipv4.String())
	}
}

func TestStartPacketReversedIPv4(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		DstMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		SrcMac: c.GetMAC(),
		IpType: 4,
		DstIp: []byte{ 127, 0, 0, 1 },
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedDownload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedDownload)
	}

	if e.accumulatedUpload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedUpload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv4.String() != "127.0.0.1" {
		t.Error("IPv4 should be 127.0.0.1, but is", e.ipv4.String())
	}
}

func TestStartPacketIPv4Broadcast(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		DstMac: []byte{0xFF, 0x22, 0x33, 0x44, 0x55, 0x66},
		SrcMac: c.GetMAC(),
		IpType: 4,
	}
	c.Close()

	if len(s.db) != 0 {
		t.Error("Should not be entries in the in-memory DB")
	}
}

func TestStartPacketNotReversedIPv6(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		SrcMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		DstMac: c.GetMAC(),
		IpType: 6,
		SrcIp: []byte{ 0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01 },
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedUpload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedUpload)
	}

	if e.accumulatedDownload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedDownload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv6.String() != "fe80::1" {
		t.Error("IPv6 should be fe80::1, but is", e.ipv6.String())
	}
}

func TestStartPacketReversedIPv6(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		DstMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		SrcMac: c.GetMAC(),
		IpType: 6,
		DstIp: []byte{ 0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01 },
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedDownload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedDownload)
	}

	if e.accumulatedUpload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedUpload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv6.String() != "fe80::1" {
		t.Error("IPv6 should be fe80::1, but is", e.ipv6.String())
	}
}

func TestStartPacketIPv6Broadcast(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		DstMac: []byte{0x33, 0x33, 0x33, 0x44, 0x55, 0x66},
		SrcMac: c.GetMAC(),
		IpType: 4,
	}
	c.Close()

	if len(s.db) != 0 {
		t.Error("Should not be entries in the in-memory DB")
	}
}

func TestStartPacketNotReversedNoNet(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		SrcMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		DstMac: c.GetMAC(),
		IpType: 0,
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedUpload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedUpload)
	}

	if e.accumulatedDownload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedDownload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv6 != nil || e.ipv4 != nil {
		t.Error("There should not be IP")
	}
}

func TestStartPacketReversedNoNet(t *testing.T) {
	s := Storage{ db: make(map[string]Entry) }
	c := dumbCapturer{ p: make(chan *capture.Packet) }
	d := dumbDB{}

	go s.Start(&c, &d)
	c.p <- &capture.Packet{
		Bytes: 100,
		DataBytes: 20,
		DstMac: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
		SrcMac: c.GetMAC(),
		IpType: 0,
	}
	c.Close()

	if len(s.db) != 1 {
		t.Error("Should be one entry in the in-memory DB")
		t.FailNow()
	}

	if _, ok := s.db["11:22:33:44:55:66"]; !ok {
		t.Error("The entry 11:22:33:44:55:66 is not there")
		t.FailNow()
	}

	e := s.db["11:22:33:44:55:66"]
	if e.accumulatedDownload != 100 {
		t.Error("Accumulated upload is not 100, is", e.accumulatedDownload)
	}

	if e.accumulatedUpload != 0 {
		t.Error("Accumulated download is not 0, is", e.accumulatedUpload)
	}

	if e.mac.String() != "11:22:33:44:55:66" {
		t.Error("MAC should be 11:22:33:44:55:66, but is", e.mac.String())
	}

	if e.ipv6 != nil || e.ipv4 != nil {
		t.Error("There should not be IP")
	}
}
