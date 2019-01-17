package main

import (
	"fmt"
	"log"
	"os"
)

// nolint: gocyclo
func main() {

	// closing logfile
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			log.Fatalf("error closing log file: %v", err)
		}
	}(Fd)

	// GracefullExit logs error to defined logger and exits gracefully
	GracefullExit := func(err error) {

		// log hook exit
		if err != nil {
			Logger.Println("Graceful exit for libvirt, but errors occurred")
		} else {
			Logger.Println("Graceful exit for libvirt, no errors occurred")
		}

		// exit with 0 code, else libvirt daemon will fail to start VM
		os.Exit(0)
	}

	switch os.Args[2] {

	// switch on: `qemu vm1 {prepare} begin -`
	case "prepare":

		switch os.Args[3] {

		// switch on: `qemu vm1 prepare {begin} -`
		case "begin":

			// get Libvirt Domain XML as object
			domCfg, err := GetDomainXML(os.Stdin)
			if err != nil {
				GracefullExit(err)
			}

			if val, ok := c.VMs[domCfg.Name]; ok {

				// Uplink v4
				err = EnableIPv4ForwardingOnInterface(c.Interface.Uplink.Name)
				if err != nil {
					GracefullExit(err)
				}

				// Uplink v6
				err = EnableIPv6ForwardingOnInterface(c.Interface.Uplink.Name)
				if err != nil {
					GracefullExit(err)
				}

				// upper interface name inside Veth pair, this device is used for L3 routing
				upperVethDev := fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Upper, val.Interface.ID)

				// VxLAN
				if val.Interface.VxLAN.VNI != 0 { // skip for Non-Defined VxLAN VNI
					err = CreateVxLANInterface(
						fmt.Sprintf("%s%d", c.Interface.Name.Prefix.VxLAN.Source, val.Interface.VxLAN.VNI),
						val.Interface.VxLAN.VNI,
						c.Interface.Uplink.Name,
					)
					if err != nil {
						GracefullExit(err)
					}
				}

				// Veth
				err = CreateVethInterface(
					upperVethDev,
					fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Source, val.Interface.ID),
				)
				if err != nil {
					GracefullExit(err)
				}

				// IPv4
				for _, ipv4 := range val.Interface.L3.IPv4 {
					err = AddStaticV4Route(ipv4, upperVethDev)
					if err != nil {
						GracefullExit(err)
					}

					err = EnableIPv4ProxyARPOnInterface(upperVethDev)
					if err != nil {
						GracefullExit(err)
					}

					err = EnableIPv4ForwardingOnInterface(upperVethDev)
					if err != nil {
						GracefullExit(err)
					}
				}

				// IPv6
				for _, ipv6 := range val.Interface.L3.IPv6 {
					err = AddStaticV6Route(ipv6, upperVethDev)
					if err != nil {
						GracefullExit(err)
					}

					err = AddVMGatewayForIPv6(ipv6, upperVethDev)
					if err != nil {
						GracefullExit(err)
					}

					err = EnableIPv6ProxyNDPOnInterface(upperVethDev)
					if err != nil {
						GracefullExit(err)
					}

					err = EnableIPv6ForwardingOnInterface(upperVethDev)
					if err != nil {
						GracefullExit(err)
					}
				}

			}

			GracefullExit(err)

		default:
			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {started} begin -`
	case "started":

		switch os.Args[3] {

		// switch on: `qemu vm1 started {begin} -`
		case "begin":

			// get Libvirt Domain XML as object
			domCfg, err := GetDomainXML(os.Stdin)
			if err != nil {
				GracefullExit(err)
			}

			if val, ok := c.VMs[domCfg.Name]; ok {

				// TC on L3
				err = ConfigureTrafficControlOnInterface(
					val.Interface.L3.TC.Rate,
					val.Interface.L3.TC.Burst,
					val.Interface.L3.TC.Limit,
					fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Target, val.Interface.ID),
				)
				if err != nil {
					GracefullExit(err)
				}

				// TC on VxLAN
				if val.Interface.VxLAN.VNI != 0 { // skip for Non-Defined VxLAN VNI
					err = ConfigureTrafficControlOnInterface(
						val.Interface.VxLAN.TC.Rate,
						val.Interface.VxLAN.TC.Burst,
						val.Interface.VxLAN.TC.Limit,
						fmt.Sprintf("%s%s", c.Interface.Name.Prefix.VxLAN.Target, val.Interface.ID),
					)
					if err != nil {
						GracefullExit(err)
					}
				}

			}

		default:
			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {stopped} end -`
	case "stopped":

		switch os.Args[3] {

		// switch on: `qemu vm1 stopped {end} -`
		case "end":

			// get Libvirt Domain XML as object
			domCfg, err := GetDomainXML(os.Stdin)
			if err != nil {
				GracefullExit(err)
			}

			if val, ok := c.VMs[domCfg.Name]; ok {

				// upper interface name inside Veth pair, this device is used for L3 routing
				upperVethDev := fmt.Sprintf("%s%s", c.Interface.Name.Prefix.Veth.Upper, val.Interface.ID)

				// Veth
				err := DestroyVethInterface(upperVethDev)
				if err != nil {
					GracefullExit(err)
				}

				GracefullExit(err)
			}

		default:
			GracefullExit(nil)
		}

	default:
		GracefullExit(nil)
	}

	GracefullExit(nil)
}
