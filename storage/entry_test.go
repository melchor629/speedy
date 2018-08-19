package storage

import (
	"testing"
	"time"
	"net"
)

var entry = Entry{
	mac: []byte{0, 1, 2, 3, 4, 5},
	ipv4: []byte{127, 0, 0, 1},
	ipv6: []byte{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
	accumulatedDownload: 123,
	accumulatedUpload: 321,
	lastModified: time.Unix(0, 0),
}

func TestGetsCorrectIPv6(t *testing.T) {
	if !entry.Ipv6().Equal(net.ParseIP("fe80::1")) {
		t.Error("IPv6 is not what was expected: got", entry.Ipv6(), "but expecting fe80::1")
	}
}

func TestGetsCorrectIPv4(t *testing.T) {
	if !entry.Ipv4().Equal(net.ParseIP("127.0.0.1")) {
		t.Error("IPv4 is not what was expected: got", entry.Ipv4(), "but expecting 127.0.0.1")
	}
}

func TestGetsCorrectMAC(t *testing.T) {
	if entry.Mac().String() != "00:01:02:03:04:05" {
		t.Error("MAC is not what was expected: got", entry.Mac(), "but expecting 00:01:02:03:04:05")
	}
}

func TestGetsCorrectAccumulatedDownload(t *testing.T) {
	if entry.GetDownloadSpeed() != 123 {
		t.Error("GetDownloadSpeed is not what was expected: got", entry.GetDownloadSpeed(), "but expecting 123")
	}
}

func TestGetsCorrectAccumulatedUpload(t *testing.T) {
	if entry.GetUploadSpeed() != 321 {
		t.Error("GetUploadSpeed is not what was expected: got", entry.GetUploadSpeed(), "but expecting 321")
	}
}

func TestTooOldGetsNotified(t *testing.T) {
	if !entry.tooOld() {
		t.Error("tooOld is not what was expected: got false but expecting true")
	}
}

func TestModifiedModifiesTime(t *testing.T) {
	entry.modified()
	if entry.tooOld() {
		t.Error("tooOld is not what was expected: got true but expecting false")
	}
}

func TestClearSpeedClearsAccumulatedVariables(t *testing.T) {
	entry.accumulatedDownload = 123
	entry.accumulatedUpload = 123

	entry.ClearSpeed()

	if entry.accumulatedUpload != 0 {
		t.Error("ClearSpeed didn't set accumulatedUpload to 0")
	}
	if entry.accumulatedDownload != 0 {
		t.Error("ClearSpeed didn't set accumulatedDownload to 0")
	}
}