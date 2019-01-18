package main

import (
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// IsValidInterfaceName - validates name for interface (ascii, nospaces, max-length 15)
func IsValidInterfaceName(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^([a-zA-Z0-9]{1,15})$").MatchString(fl.Field().String())
}

// IsNotIPv6NetworkAddress - validates that IPv6 is not network address for /64 of it self
func IsNotIPv6NetworkAddress(fl validator.FieldLevel) bool {
	// store input IPv6
	ipv6 := fl.Field().String()
	// compute GW for /64
	gw6 := GetNetworkAddressFromIPv6(ipv6)
	// compare
	return !strings.EqualFold(ipv6, gw6)
}
