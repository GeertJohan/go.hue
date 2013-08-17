package hue

import (
	"net"
)

// Bridge represents a Hue Bridge
type Bridge struct {
	id        string
	ipAddress *net.IPAddr
}

// ID returns the ID of the Bridge as string
func (b *Bridge) ID() string {
	return b.id
}

// IP returns the IP address of the Bridge as *net.IP
func (b *Bridge) IP() *net.IP {
	//++
	return nil
}
