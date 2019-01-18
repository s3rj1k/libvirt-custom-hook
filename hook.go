package main

import (
	"fmt"

	"github.com/libvirt/libvirt-go-xml"
)

// LookupVMConfig - wrapper to get VM configuration by 'UUID' first and if this failes try to lookup by 'Name'
func (c *Config) LookupVMConfig(domCfg *libvirtxml.Domain) (VM, error) {

	var ok bool

	// declare variable for VM configuration description
	var vm VM

	// check config for defined VM by UUID
	vm, ok = c.VMs[domCfg.UUID]
	if ok {
		return vm, nil
	}

	// check config for defined VM by Name
	vm, ok = c.VMs[domCfg.Name]
	if ok {
		return vm, nil
	}

	return VM{}, fmt.Errorf("no VM found in config for UUID='%s' or Name='%s'", domCfg.UUID, domCfg.Name)
}

// PrepareBeginHook - hook for `qemu vm1 prepare begin -`
func (c *Config) PrepareBeginHook(domCfg *libvirtxml.Domain) error {

	// lookup VM config
	vm, err := c.LookupVMConfig(domCfg)
	if err != nil {
		return err
	}

	// Uplink v4
	err = EnableIPv4ForwardingOnInterface(c.Interface.Uplink.Name)
	if err != nil {
		return err
	}

	// Uplink v6
	err = EnableIPv6ForwardingOnInterface(c.Interface.Uplink.Name)
	if err != nil {
		return err
	}

	// upper interface name inside Veth pair, this device is used for L3 routing
	upperVethDev := fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Upper, vm.Interface.ID)

	// VxLAN
	if vm.Interface.VxLAN.VNI != 0 { // skip for Non-Defined VxLAN VNI
		err = CreateVxLANInterface(
			fmt.Sprintf("%s%d", c.Interface.Name.Prefix.VxLAN.Source, vm.Interface.VxLAN.VNI),
			vm.Interface.VxLAN.VNI,
			c.Interface.Uplink.Name,
		)
		if err != nil {
			return err
		}
	}

	// Veth
	err = CreateVethInterface(
		upperVethDev,
		fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Source, vm.Interface.ID),
	)
	if err != nil {
		return err
	}

	// IPv4
	for _, ipv4 := range vm.Interface.L3.IPv4 {
		err = AddStaticV4Route(ipv4, upperVethDev)
		if err != nil {
			return err
		}

		err = EnableIPv4ProxyARPOnInterface(upperVethDev)
		if err != nil {
			return err
		}

		err = EnableIPv4ForwardingOnInterface(upperVethDev)
		if err != nil {
			return err
		}
	}

	// IPv6
	for _, ipv6 := range vm.Interface.L3.IPv6 {
		err = AddStaticV6Route(ipv6, upperVethDev)
		if err != nil {
			return err
		}

		err = AddVMGatewayForIPv6(ipv6, upperVethDev)
		if err != nil {
			return err
		}

		err = EnableIPv6ProxyNDPOnInterface(upperVethDev)
		if err != nil {
			return err
		}

		err = EnableIPv6ForwardingOnInterface(upperVethDev)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartedBeginHook - hook for `qemu vm1 started begin -`
func (c *Config) StartedBeginHook(domCfg *libvirtxml.Domain) error {

	// lookup VM config
	vm, err := c.LookupVMConfig(domCfg)
	if err != nil {
		return err
	}

	// TC on L3
	err = ConfigureTrafficControlOnInterface(
		vm.Interface.L3.TC.Rate,
		vm.Interface.L3.TC.Burst,
		vm.Interface.L3.TC.Limit,
		fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Target, vm.Interface.ID),
	)
	if err != nil {
		return err
	}

	// TC on VxLAN
	if vm.Interface.VxLAN.VNI != 0 { // skip for Non-Defined VxLAN VNI
		err = ConfigureTrafficControlOnInterface(
			vm.Interface.VxLAN.TC.Rate,
			vm.Interface.VxLAN.TC.Burst,
			vm.Interface.VxLAN.TC.Limit,
			fmt.Sprintf("%s%s", c.Interface.Name.Prefix.VxLAN.Target, vm.Interface.ID),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// StoppedEndHook - hook for `qemu vm1 stopped end -`
func (c *Config) StoppedEndHook(domCfg *libvirtxml.Domain) error {

	// lookup VM config
	vm, err := c.LookupVMConfig(domCfg)
	if err != nil {
		return err
	}

	// upper interface name inside Veth pair, this device is used for L3 routing
	upperVethDev := fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Upper, vm.Interface.ID)

	// Veth
	return DestroyVethInterface(upperVethDev)
}
