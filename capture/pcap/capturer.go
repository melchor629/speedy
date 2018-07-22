//Traffic capture using libpcap
package pcap

import (
	"github.com/google/gopacket/pcap"
	"speedy/capture"
	"net"
	"log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"os"
)

//Capture the trafic using libpcap
type CaptureContext struct {
	device string
	handle *pcap.Handle
	stop chan bool
	packetsChan chan *capture.Packet
	mac net.HardwareAddr
	logger *log.Logger
}

type layersPack struct {
	parser *gopacket.DecodingLayerParser
	decoded []gopacket.LayerType
	eth layers.Ethernet
	ip4 layers.IPv4
	ip6 layers.IPv6
	tcp layers.TCP
	udp layers.UDP
}

//Creates a capture context using libpcap implementation and opens the device to capture. Ensure that the process has
//permission con capture traffic through the device.
func New(device string) (*CaptureContext, error) {
	handle, err := pcap.OpenLive(device, 65535, true, 0)
	if err != nil {
		return nil, err
	}

	mac, err := getMacAddr(device)
	if err != nil {
		return nil, err
	}

	return &CaptureContext{
		device,
		handle,
		make(chan bool),
		make(chan *capture.Packet),
		mac,
		log.New(os.Stdout, "[Context]: ", log.LstdFlags),
	}, nil
}

//Ends with the capture session
func (c *CaptureContext) Close() {
	c.logger.Println("Closing...")
	c.stop <- true
	<- c.stop
	c.handle.Close()
	close(c.packetsChan)
	close(c.stop)
}

//Starts the capture session. Use Packets() to grab the packets channel.
func (c *CaptureContext) StartCapturing() {
	c.logger.Println("Starting capture gorutine")
	c.logger.Println("Capturing", c.device, "with MAC", c.mac.String())
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	lp := layersPack{}
	lp.parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &lp.eth, &lp.ip4, &lp.ip6, &lp.tcp, &lp.udp)

	itsTimeToStop := false
	for !itsTimeToStop {
		select {
		case <- c.stop:
			c.logger.Println("Stopping capturer gorutine")
			itsTimeToStop = true
		case packet := <- packetSource.Packets():
			c.packetsChan <- parsePacket(packet, &lp)
		}
	}

	c.stop <- true
}

//Returns the packets channel where all the packets will be passed through.
func (c *CaptureContext) Packets() chan *capture.Packet {
	return c.packetsChan
}

//Gets the MAC Address of the device being captured.
func (c *CaptureContext) GetMAC() net.HardwareAddr {
	return c.mac
}

func parsePacket(packet gopacket.Packet, lp *layersPack) *capture.Packet {
	lp.parser.DecodeLayers(packet.Data(), &lp.decoded)
	var ppacket capture.Packet
	for _, layerType := range lp.decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			ppacket.Bytes = uint16(len(lp.eth.Payload))
			ppacket.SrcMac = lp.eth.SrcMAC
			ppacket.DstMac = lp.eth.DstMAC
		case layers.LayerTypeIPv4:
			ppacket.SrcIp = lp.ip4.SrcIP
			ppacket.DstIp = lp.ip4.DstIP
			ppacket.IpType = 4
		case layers.LayerTypeIPv6:
			ppacket.SrcIp = lp.ip6.SrcIP
			ppacket.DstIp = lp.ip6.DstIP
			ppacket.IpType = 6
		case layers.LayerTypeTCP:
			if ppacket.IpType == 4 {
				ppacket.DataBytes = uint16(len(lp.ip4.Payload)) - uint16(lp.tcp.DataOffset * 4)
			} else {
				ppacket.DataBytes = uint16(len(lp.ip6.Payload)) - uint16(lp.tcp.DataOffset * 4)
			}
		case layers.LayerTypeUDP:
			ppacket.DataBytes = lp.udp.Length
		}
	}

	return &ppacket
}

//Based on https://gist.github.com/rucuriousyet/ab2ab3dc1a339de612e162512be39283
func getMacAddr(name string) (net.HardwareAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		if i.Name == name {
			return i.HardwareAddr, nil
		}
	}
	return nil, nil
}

//Gets all the active network interfaces (only their names).
func GetActiveInterfaces() ([]string, error) {
	//Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	validDevices := make([]string, 0)
	for _, device := range devices {
		if len(device.Addresses) != 0 {
			validDevices = append(validDevices, device.Name)
		}
	}
	return validDevices, nil
}
