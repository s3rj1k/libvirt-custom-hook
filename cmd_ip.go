package main

import (
	"fmt"
	"strconv"
)

// CreateVxLANInterface - creates VxLAN interface inside host node with specified VNI
func CreateVxLANInterface(name string, vni int64, dev string) error {

	// prefix for errors logging
	const errPrefix = "vxlan config error:"

	// create VxLAN interface
	cmd := RunCommand("ip", "link", "add", "name", SanitizeInput(name),
		"type", "vxlan", "id", strconv.FormatInt(vni, 10),
		"dev", SanitizeInput(dev),
		"group", "239.0.0.1", "dstport", "4789",
	)
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// bring VxLAN interface to UP state
	cmd = RunCommand("ip", "link", "set", "dev", SanitizeInput(name), "up")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// CreateVethInterface - creates Veth pair interface inside host node
func CreateVethInterface(upper, lower string) error {

	// prefix for errors logging
	const errPrefix = "veth config error:"

	// create Veth interface
	cmd := RunCommand("ip", "link", "add", "name", SanitizeInput(upper), "type", "veth", "peer", "name", SanitizeInput(lower))
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// bring upper veth pair to UP state
	cmd = RunCommand("ip", "link", "set", "dev", SanitizeInput(upper), "up")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// bring lower veth pair to UP state
	cmd = RunCommand("ip", "link", "set", "dev", SanitizeInput(lower), "up")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// DestroyVethInterface - deletes previosly created Veth interface
func DestroyVethInterface(dev string) error {

	// prefix for errors logging
	const errPrefix = "veth config error:"

	// get extended information previosly created interface
	cmd := RunCommand("ip", "-o", "-d", "l", "show", SanitizeInput(dev), "type", "veth")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// remove interface
	cmd = RunCommand("ip", "link", "del", SanitizeInput(dev), "type", "veth")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// AddStaticV4Route - adds static route for IPv4/32 to specified interface
func AddStaticV4Route(ip string, dev string) error {

	// prefix for errors logging
	const errPrefix = "route4 config error:"

	// add static v4 route
	cmd := RunCommand("ip", "-4", "route", "add", fmt.Sprintf("%s/32", SanitizeInput(ip)), "dev", SanitizeInput(dev))
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// AddStaticV6Route - adds static route for IPv6/128 to specified interface
func AddStaticV6Route(ip string, dev string) error {

	// prefix for errors logging
	const errPrefix = "route6 config error:"

	// add static v6 route
	cmd := RunCommand("ip", "-6", "route", "add", fmt.Sprintf("%s/128", SanitizeInput(ip)), "dev", SanitizeInput(dev))
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}

// AddVMGatewayForIPv6 - adds network address computed from IPv6/64 network ad gateway for V6 routing used in VM
func AddVMGatewayForIPv6(ip string, dev string) error {

	// prefix for errors logging
	const errPrefix = "vmgw6 config error:"

	// add IPv6 address to device, for VM usage as gateway, used for v6 routing
	cmd := RunCommand("ip", "-6", "addr", "add", fmt.Sprintf("%s/64", GetNetworkAddressFromIPv6(SanitizeInput(ip))),
		"dev", SanitizeInput(dev),
		"noprefixroute", "nodad", "scope", "link",
	)
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}
