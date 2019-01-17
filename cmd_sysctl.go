package main

import (
	"fmt"
)

// EnableIPv4ForwardingOnInterface - enables IPv4 forwarding on specified interface
func EnableIPv4ForwardingOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "forwarding4 config error:"

	// enable IPv4 forwarding
	cmd := RunCommand("sysctl", "-w", fmt.Sprintf("net.ipv4.conf.%s.forwarding=1", SanitizeInput(dev)))
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv6ForwardingOnInterface - enables IPv6 forwarding on specified interface
func EnableIPv6ForwardingOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "forwarding6 config error:"

	// enable IPv6 forwarding
	cmd := RunCommand("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.forwarding=1", SanitizeInput(dev)))
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// enable IPv6 forwarding globally, this differs from IPv4 behavior, consult kernel docs
	cmd = RunCommand("sysctl", "-w", "net.ipv6.conf.all.forwarding=1")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv4ProxyARPOnInterface - enables IPv4 ProxyARP on specified interface
func EnableIPv4ProxyARPOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "proxyARP config error:"

	// enable IPv4 ProxyARP
	cmd := RunCommand("sysctl", "-w", fmt.Sprintf("net.ipv4.conf.%s.proxy_arp=1", SanitizeInput(dev)))
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// EnableIPv6ProxyNDPOnInterface - enables IPv4 ProxyNDP on specified interface
func EnableIPv6ProxyNDPOnInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "proxyNDP config error:"

	// enable IPv6 ProxyNDP
	cmd := RunCommand("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.proxy_ndp=1", SanitizeInput(dev)))
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}
