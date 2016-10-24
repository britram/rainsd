package main

import (
	"fmt"
	"net"
)

const RAINSD_PORT = 1228

// Switchboard: manage TCP connections

type Switchboard struct {
	connectedPeers map[string]net.Conn
}

func (sb *Switchboard) sendMessage(peerName string, message string) error {
	var c net.Conn
	var ok bool

	c, ok = sb.connectedPeers[peerName]

	if !ok {
		// open a connection and add to the map
		c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peerName))
		if err != nil {
			return err
		}
		sb.connectedPeers[peerName] = c
	}

	// Framing goes here when we define it
	_, err := fmt.Fprintf(c, "%s\n", message)
	if err != nil {
		sb.reapConnection(peerName)
		return err
	}

	return nil
}

func (sb *Switchboard) reapConnection(peerName string) {
	c, ok := sb.connectedPeers[peerName]
	if ok {
		delete(sb.connectedPeers, peerName)
		c.Close()
	}
}

func main() {

}
