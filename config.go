package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// TC - traffic control config, for basic traffic shaping
type TC struct {
	// mbit
	Rate int64 `json:"Rate" validate:"required,min=1"`
	// kb
	Burst int64 `json:"Burst" validate:"required,min=1"`
	// packets
	Limit int64 `json:"Limit" validate:"required,min=10240"`
}

// VM - config per VM
type VM struct {
	Interface struct {
		// 15(max interface length, kernel limit) = 12(max ID length) + 3(max prefix length)
		ID string `json:"ID" validate:"required,printascii,min=1,max=12"`

		VxLAN struct {
			// assume that VNI == 0, no VxLAN
			VNI int64 `json:"VNI" validate:"required,min=0,max=16777214"`
			TC  TC    `json:"TC" validate:"required"`
		} `json:"VxLAN" validate:"required"`

		L3 struct {
			IPv4 []string `json:"IPv4" validate:"required,unique,dive,ipv4"`
			IPv6 []string `json:"IPv6" validate:"required,unique,dive,ipv6"`
			TC   TC       `json:"TC" validate:"required"`
		} `json:"L3" validate:"required"`
	} `json:"Interface" validate:"required"`
}

// Interface - interface prefixes inside libvirt Domain XML and in hypervisor node
type Interface struct {
	Name struct {
		Prefix struct {
			VxLAN struct {
				Source string `json:"Source" validate:"required,min=1,max=3,printascii"`
				Target string `json:"Target" validate:"required,min=1,max=3,printascii"`
			} `json:"VxLAN" validate:"required"`
			Veth struct {
				Upper  string `json:"Upper" validate:"required,min=1,max=3,printascii"`
				Source string `json:"Source" validate:"required,min=1,max=3,printascii"`
				Target string `json:"Target" validate:"required,min=1,max=3,printascii"`
			} `json:"Veth" validate:"required"`
		} `json:"Prefix" validate:"required"`
	} `json:"Name" validate:"required"`
	Uplink struct {
		Name string `json:"Name" validate:"required,min=1,max=15"`
	} `json:"Uplink" validate:"required"`
}

// Config - main hook config
type Config struct {
	Interface Interface     `json:"Interface" validate:"required"`
	VMs       map[string]VM `json:"VMs" validate:"required"`
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
