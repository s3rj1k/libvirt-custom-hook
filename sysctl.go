package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// SysctlGet - wrapper to read value from sysctl file
func SysctlGet(path string) (string, error) {

	// read file to memory
	content, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	// split by newlines
	lines := bytes.Split(content, []byte("\n"))

	// get value from first line
	value := string(bytes.TrimSpace(lines[0]))

	return value, nil
}

// SysctlCheckEqual - wrapper to compare supplied value with value in systcl file
func SysctlCheckEqual(path string, value string) (bool, error) {

	// get current value
	currentValue, err := SysctlGet(filepath.Clean(path))
	if err != nil {
		return false, err
	}

	// compare
	if strings.EqualFold(currentValue, value) {
		return true, nil
	}

	return false, nil
}

// SysctlSet - wrapper to set sysctl file value
func SysctlSet(path string, value string) error {

	// check the need to update sysctl file
	ok, err := SysctlCheckEqual(filepath.Clean(path), value)
	if err != nil {
		return err
	}

	// current value is same as supplied, do nothing
	if ok {
		return nil
	}

	return ioutil.WriteFile(filepath.Clean(path), []byte(value), 0)
}

// EnableIPv4ForwardingOnInterface - enables IPv4 forwarding on specified interface
func EnableIPv4ForwardingOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "sysctl config error:"

	// enable IPv4 forwarding: sysctl -w net.ipv4.conf.%s.forwarding=1
	err := SysctlSet(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/forwarding", SanitizeInput(dev)), "1")
	if err != nil {
		e := fmt.Errorf("%s failed to enable IPv4 forwarding for '%s' device: %s", errPrefix, SanitizeInput(dev), err.Error())
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv6ForwardingOnInterface - enables IPv6 forwarding on specified interface
func EnableIPv6ForwardingOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "sysctl config error:"

	// enable IPv6 forwarding: sysctl -w net.ipv6.conf.%s.forwarding=1
	err := SysctlSet(fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/forwarding", SanitizeInput(dev)), "1")
	if err != nil {
		e := fmt.Errorf("%s failed to enable IPv6 forwarding for '%s' device: %s", errPrefix, SanitizeInput(dev), err.Error())
		Logger.Println(e)

		return e
	}

	// enable IPv6 forwarding globally, this differs from IPv4 behavior, consult kernel docs: sysctl -w net.ipv6.conf.all.forwarding=1
	err = SysctlSet("/proc/sys/net/ipv6/conf/all/forwarding", "1")
	if err != nil {
		e := fmt.Errorf("%s failed to enable IPv6 forwarding: %s", errPrefix, err.Error())
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv4ProxyARPOnInterface - enables IPv4 ProxyARP on specified interface
func EnableIPv4ProxyARPOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "sysctl config error:"

	// enable IPv4 ProxyARP: sysctl -w net.ipv4.conf.%s.proxy_arp=1
	err := SysctlSet(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/proxy_arp", SanitizeInput(dev)), "1")
	if err != nil {
		e := fmt.Errorf("%s failed to enable IPv4 ProxyARP for '%s' device: %s", errPrefix, SanitizeInput(dev), err.Error())
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv6ProxyNDPOnInterface - enables IPv6 ProxyNDP on specified interface
func EnableIPv6ProxyNDPOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "sysctl config error:"

	// enable IPv6 ProxyNDP: sysctl -w net.ipv6.conf.%s.proxy_ndp=1
	err := SysctlSet(fmt.Sprintf("/proc/sys/net/ipv6/conf/%s/proxy_ndp", SanitizeInput(dev)), "1")
	if err != nil {
		e := fmt.Errorf("%s failed to enable IPv6 ProxyNDP for '%s' device: %s", errPrefix, SanitizeInput(dev), err.Error())
		Logger.Println(e)

		return e
	}

	return nil
}
