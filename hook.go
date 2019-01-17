package main

import (
	"fmt"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

// LookupVMConfig - wrapper to get VM configuration by 'UUID' first and if this failes try to lookup by 'Name'
func (c *Config) LookupVMConfig(domCfg *libvirtxml.Domain) (VM, error) {
	var ok bool

	// declare variable for VM configuration description
	var vm VM

	// prefix for errors logging
	const errPrefix = "vm config error:"

	// check config for defined VM by UUID
	vm, ok = c.VMs[domCfg.UUID]
	if ok {
		// run validator on VM config
		err := Validate.Struct(vm)
		if err != nil {
			// log error, invalid config
			e := fmt.Errorf("%s %s", errPrefix, err.Error())
			Logger.Println(e)

			return VM{}, e
		}

		return vm, nil
	}

	// check config for defined VM by Name
	vm, ok = c.VMs[domCfg.Name]
	if ok {
		// run validator on VM config
		err := Validate.Struct(vm)
		if err != nil {
			// log error, invalid config
			e := fmt.Errorf("%s %s", errPrefix, err.Error())
			Logger.Println(e)

			return VM{}, e
		}

		return vm, nil
	}

	// log error, no VM in config
	e := fmt.Errorf("%s no VM found in config for UUID='%s' or Name='%s'", errPrefix, domCfg.UUID, domCfg.Name)
	Logger.Println(e)

	return VM{}, e
}

// PrepareBeginHook - hook for `qemu vm1 prepare begin -`
func (c *Config) PrepareBeginHook(domCfg *libvirtxml.Domain) error {
	// lookup VM config
	vm, err := c.LookupVMConfig(domCfg)
	if err != nil {
		return err
	}

	// Validate Uplink interface existence
	if !IsInterfaceExists(vm.Interface.Uplink.Name) {
		return fmt.Errorf("hook: uplink interface '%s' does not exist", vm.Interface.Uplink.Name)
	}

	// Uplink v4
	err = EnableIPv4ForwardingOnInterface(vm.Interface.Uplink.Name)
	if err != nil {
		return err
	}

	// Uplink v6
	err = EnableIPv6ForwardingOnInterface(vm.Interface.Uplink.Name)
	if err != nil {
		return err
	}

	// VxLAN
	if vm.Interface.VxLAN != nil { // skip for Non-Defined VxLAN
		err = CreateVxLANInterface(
			vm.Interface.VxLAN.Source.Name,
			vm.Interface.VxLAN.VNI,
			vm.Interface.Uplink.Name,
		)
		if err != nil {
			return err
		}
	}

	// Veth
	err = CreateVethInterface(
		vm.Interface.L3.Upper.Name,
		vm.Interface.L3.Source.Name,
	)
	if err != nil {
		return err
	}

	// IPv4
	for _, ipv4 := range vm.Interface.L3.IPv4 {
		err = AddStaticV4Route(ipv4, vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}

		err = EnableIPv4ProxyARPOnInterface(vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}

		err = EnableIPv4ForwardingOnInterface(vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}
	}

	// IPv6
	for _, ipv6 := range vm.Interface.L3.IPv6 {
		err = AddStaticV6Route(ipv6, vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}

		err = AddVMGatewayForIPv6(ipv6, vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}

		err = EnableIPv6ProxyNDPOnInterface(vm.Interface.L3.Upper.Name)
		if err != nil {
			return err
		}

		err = EnableIPv6ForwardingOnInterface(vm.Interface.L3.Upper.Name)
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
		vm.Interface.L3.Target.Name,
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
			vm.Interface.VxLAN.Target.Name,
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

	// Veth
	return DestroyVethInterface(vm.Interface.L3.Upper.Name)
}
