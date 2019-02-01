package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// VM - config per VM
type VM struct {
	Interface *Interface `json:"Interface" validate:"required"`
}

// Interface - interfaces configuration for VM
type Interface struct {
	VxLAN  *VxLAN `json:"VxLAN" validate:"omitempty"`
	L3     *L3    `json:"L3" validate:"required"`
	Uplink *Iface `json:"Uplink" validate:"required"`
}

// VxLAN - Private LAN configuration
type VxLAN struct {
	// assume that VNI == 0, no VxLAN
	VNI int64 `json:"VNI" validate:"required,min=1,max=16777214"`
	// usually uplink
	Source *Iface `json:"Source" validate:"required"`
	// created by libvirt
	Target *Iface `json:"Target" validate:"required"`
	TC     *TC    `json:"TC" validate:"required"`
}

// L3 - Internet configuration for VM
type L3 struct {
	// upper peer of Veth pair
	Upper *Iface `json:"Upper" validate:"required"`
	// lower peer of Veth pair
	Source *Iface `json:"Source" validate:"required"`
	// created by libvirt
	Target *Iface   `json:"Target" validate:"required"`
	TC     *TC      `json:"TC" validate:"required"`
	IPv4   []string `json:"IPv4" validate:"required,unique,dive,ipv4"`
	IPv6   []string `json:"IPv6" validate:"unique,dive,ipv6,notGW6"`
}

// Iface - represents interface name
type Iface struct {
	Name string `json:"Name" validate:"required,iface"`
}

// TC - traffic control config, for basic traffic shaping
type TC struct {
	// mbit
	Rate int64 `json:"Rate" validate:"required,min=1"`
	// kb
	Burst int64 `json:"Burst" validate:"required,min=1"`
	// packets
	Limit int64 `json:"Limit" validate:"required,min=10240"`
}

// Config - main hook config
type Config struct {
	VMs map[string]VM `json:"VMs" validate:"required"`
}

// GetConfig - get application configuration
func GetConfig(path string) (*Config, error) {

	// prefix for errors logging
	const errPrefix = "config error:"

	var err error

	// create config object
	c = new(Config)

	// read config file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errPrefix, err)
	}

	// convert config file to object
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errPrefix, err)
	}

	// additional structe validation
	err = Validate.Struct(c)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errPrefix, err)
	}

	return c, nil
}
