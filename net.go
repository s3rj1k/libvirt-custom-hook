package main

import (
	"net"
)

// IsInterfaceExists - check interface existence
func IsInterfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)
	if err != nil { // nolint: megacheck
		return false
	}

	return true
}
