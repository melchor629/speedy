package capture

import "net"

//Filtered/processed packet from a capture context. Holds the necessary information for the app to show the data.
type Packet struct {
	Bytes uint16 //The length of the packet (without ethernet header)
	DataBytes uint16 //If possible, the length of the application layer (TCP/UDP's payload)
	SrcMac net.HardwareAddr //The source MAC address
	SrcIp net.IP //The source IP address (if available), could be IPv4 or IPv6
	DstMac net.HardwareAddr //The destination MAC address
	DstIp net.IP //The destination IP address (if available), could be IPv4 or IPv6
	IpType uint8 //Type of IP: 4, 6 or 0 (for nothing)
	reversed bool //Stores if Reverse() was called
}

//Returns true if the packet was an IPv4 one.
func (p Packet) IsIP4() bool {
	return p.IpType == 4
}

//Returns true if the packet was an IPv6 one.
func (p Packet) IsIP6() bool {
	return p.IpType == 6
}

//If the SrcMac is the same as dstMac, then is reversed. dstMac should be the MAC of the capturing device.
func (p Packet) IsReversed(dstMac net.HardwareAddr) bool {
	for i := 0; i < len(dstMac); i++  {
		if dstMac[i] != p.SrcMac[i] {
			return false
		}
	}
	return true
}

//Reverses the packet (changes the source with the destination).
func (p *Packet) Reverse() {
	p.reversed = !p.reversed
	tmp := p.DstMac
	p.DstMac = p.SrcMac
	p.SrcMac = tmp
	tmp2 := p.DstIp
	p.DstIp = p.SrcIp
	p.SrcIp = tmp2
}

//Returns true if the packet is IPv6 multicast.
func (p Packet) IsIPv6Multicast() bool {
	return p.DstMac[0] == 0x33 && p.DstMac[1] == 0x33
}

//Returns true if the packet is ethernet broadcast.
func (p Packet) IsBroadcast() bool {
	return p.DstMac[0] == 0xFF
}
