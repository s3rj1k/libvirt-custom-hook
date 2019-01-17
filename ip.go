package main

import (
	"net"
)

// GetNetworkAddressFromIPv6 - compute network address of IPv6/64 address, for gateway usage inside VM
func GetNetworkAddressFromIPv6(ip string) string {
	// mask corresponds to a /64 subnet for IPv6.
	return net.ParseIP(ip).Mask(net.CIDRMask(64, 128)).String()
}
