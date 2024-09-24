package catcher

import (
	"fmt"
	"net"

	"github.com/google/gopacket/pcap"
)

type NetworkDevice struct {
	name        string
	description string
	addresses   []Address
}

func (nd *NetworkDevice) Name() string {
	return nd.name
}

func (nd *NetworkDevice) Description() string {
	return nd.description
}

func (nd *NetworkDevice) Addresses() []Address {
	return nd.addresses
}

type Address struct {
	ip               net.IP
	netmask          net.IPMask
	broadcastAddress net.IP
	p2p              net.IP
}

func (a *Address) IP() net.IP {
	return a.ip
}

func (a *Address) Netmask() net.IPMask {
	return a.netmask
}

func (a *Address) BroadcastAddress() net.IP {
	return a.broadcastAddress
}

func (a *Address) P2P() net.IP {
	return a.p2p
}

func ListAvailableDevices() ([]NetworkDevice, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find network devices: %w", err)
	}
	return toNetworkDevices(interfaces), nil
}

func toNetworkDevices(interfaces []pcap.Interface) []NetworkDevice {
	devices := make([]NetworkDevice, len(interfaces))
	for i, deviceInterface := range interfaces {
		devices[i] = NetworkDevice{
			name:        deviceInterface.Name,
			description: deviceInterface.Description,
			addresses:   toAddresses(deviceInterface.Addresses),
		}
	}

	return devices
}

func toAddresses(interfaceAddresses []pcap.InterfaceAddress) []Address {
	addresses := make([]Address, len(interfaceAddresses))
	for i, a := range interfaceAddresses {
		addresses[i] = Address{
			ip:               a.IP,
			netmask:          a.Netmask,
			broadcastAddress: a.Broadaddr,
			p2p:              a.P2P,
		}
	}
	return addresses
}
