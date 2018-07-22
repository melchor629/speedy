package capture

import "net"

//Context for a capture session.
type Context interface {
	//Gets the MAC address of the captured interface.
	GetMAC() net.HardwareAddr
	//Gets the channel of the processed packets.
	Packets() chan *Packet
	//Tells the capture context to start the capture session.
	StartCapturing()
	//Ends the capture session.
	Close()
}
