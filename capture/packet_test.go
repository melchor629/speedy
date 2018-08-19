package capture

import "testing"

func TestIsIp4ShouldReturnTrueWhenItIsIndeed(t *testing.T) {
	packet := Packet{ IpType: 4 }

	isv4 := packet.IsIP4()

	if !isv4 {
		t.Error("Should return true")
	}
}

func TestIsIp4ShouldReturnFalseWhenItIsNot(t *testing.T) {
	packet := Packet{ IpType: 0 }

	isv4 := packet.IsIP4()

	if isv4 {
		t.Error("Should return false")
	}
}

func TestIsIp6ShouldReturnTrueWhenItIsIndeed(t *testing.T) {
	packet := Packet{ IpType: 6 }

	isv6 := packet.IsIP6()

	if !isv6 {
		t.Error("Should return true")
	}
}

func TestIsIp6ShouldReturnFalseWhenItIsNot(t *testing.T) {
	packet := Packet{ IpType: 0 }

	isv6 := packet.IsIP6()

	if isv6 {
		t.Error("Should return false")
	}
}

func TestIsReversedShouldReturnFalseWhenDstMacIsTheSameAsTheArgument(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff },
		DstMac: []byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 },
	}

	isReversed := packet.IsReversed([]byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 })

	if isReversed {
		t.Error("Should be false")
	}
}

func TestIsReversedShouldReturnTrueWhenDstMacIsTheDifferentOfTheArgument(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 },
		DstMac: []byte { 0x11, 0x22, 0x33, 0xdd, 0xee, 0xff },
	}

	isReversed := packet.IsReversed([]byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 })

	if !isReversed {
		t.Error("Should be true")
	}
}

func TestReversedChangesDstAndSrcMac(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 },
		DstMac: []byte { 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff },
	}

	packet.Reverse()

	if !packet.reversed {
		t.Error("packet.reversed should be true")
	}

	if packet.SrcMac.String() == "11:22:33:44:55:66" {
		t.Error("SrcMac is not reversed")
	}

	if packet.DstMac.String() == "aa:bb:cc:dd:ee:ff" {
		t.Error("DstMac is not reversed")
	}
}

func TestReversedChangesDstAndSrcMacIfIsReversedAlready(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0x11, 0x22, 0x33, 0x44, 0x55, 0x66 },
		DstMac: []byte { 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff },
		reversed: true,
	}

	packet.Reverse()

	if packet.reversed {
		t.Error("packet.reversed should be false")
	}

	if packet.SrcMac.String() == "11:22:33:44:55:66" {
		t.Error("SrcMac is not reversed")
	}

	if packet.DstMac.String() == "aa:bb:cc:dd:ee:ff" {
		t.Error("DstMac is not reversed")
	}
}

func TestIsIPv6MulticastShouldReturnTrueWhenItIs(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0x33, 0x33, 0x22, 0x11, 0x00, 0xFF },
	}

	is := packet.IsIPv6Multicast()

	if !is {
		t.Error("Should return true")
	}
}

func TestIsIPv6MulticastShouldReturnFalseWhenItIsNot(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0x44, 0x33, 0x22, 0x11, 0x00, 0xFF },
	}

	is := packet.IsIPv6Multicast()

	if is {
		t.Error("Should return false")
	}
}

func TestIsBroadcastShouldReturnTrueWhenItIs(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0xff, 0x33, 0x22, 0x11, 0x00, 0xFF },
	}

	is := packet.IsBroadcast()

	if !is {
		t.Error("Should return true")
	}
}

func TestIsBroadcastShouldReturnFalseWhenItIsNot(t *testing.T) {
	packet := Packet{
		SrcMac: []byte { 0xee, 0x33, 0x22, 0x11, 0x00, 0xFF },
	}

	is := packet.IsBroadcast()

	if is {
		t.Error("Should return false")
	}
}