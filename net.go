package main

import (
	"net"
)

// IsInterfaceExists - check interface existence
func IsInterfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)
	if err != nil {
		return false
	}

	return true
}
