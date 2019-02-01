package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		caseDescription string
		vm              VM    //in
		err             error //out
	}{
		{
			caseDescription: "valid config",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: nil,
		},
		{
			caseDescription: "valid config w/o L3.IPv6",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: nil,
		},
		{
			caseDescription: "valid config w/o VM.Interface.VxLAN",
			vm: VM{
				Interface: &Interface{
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: nil,
		},
		{
			caseDescription: "valid config w/o VM.Interface.VxLAN and w/o L3.IPv6",
			vm: VM{
				Interface: &Interface{
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: nil,
		},
		{
			caseDescription: "valid config. multiple IPs",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4: []string{
							"195.177.117.1",
							"195.177.118.10",
							"195.177.118.100",
						},
						IPv6: []string{
							"2a02:2278:100:a1::1",
							"2a02:2278:100:a2::1",
							"2a02:2278:100:a3::1",
						},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: nil,
		},
		{
			caseDescription: "not unique IPv4",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4: []string{
							"195.177.117.1",
							"195.177.117.1",
						},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv4' Error:Field validation for 'IPv4' failed on the 'unique' tag"),
		},
		{
			caseDescription: "invalid IPv4",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4: []string{
							"195.177.117.1000",
						},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv4[0]' Error:Field validation for 'IPv4[0]' failed on the 'ipv4' tag"),
		},
		{
			caseDescription: "invalid IPv4 #2",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4: []string{
							"195.177.117.100",
							"195.177.117.1000",
							"195.177.117.200",
						},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv4[1]' Error:Field validation for 'IPv4[1]' failed on the 'ipv4' tag"),
		},
		{
			caseDescription: "invalid IPv4 #2",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4: []string{
							"195.177.117.100",
							"195.177.117.200",
							"195.177.117.1000",
						},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv4[2]' Error:Field validation for 'IPv4[2]' failed on the 'ipv4' tag"),
		},
		{
			caseDescription: "IPv6 network address",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv6[0]' Error:Field validation for 'IPv6[0]' failed on the 'notGW6' tag"),
		},
		{
			caseDescription: "missing VM.Interface.L3",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3' Error:Field validation for 'L3' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.Interface.Uplink",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
				},
			},
			err: errors.New("Key: 'VM.Interface.Uplink' Error:Field validation for 'Uplink' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.Interface.VxLAN.VNI",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.VxLAN.VNI' Error:Field validation for 'VNI' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.Interface.VxLAN.Source",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.VxLAN.Source' Error:Field validation for 'Source' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.Interface.VxLAN.Target",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.VxLAN.Target' Error:Field validation for 'Target' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.Interface.VxLAN.TC",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.VxLAN.TC' Error:Field validation for 'TC' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.L3.IPv4",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.IPv4' Error:Field validation for 'IPv4' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.L3.Upper",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.Upper' Error:Field validation for 'Upper' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.L3.Source",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Target: &Iface{"if-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.Source' Error:Field validation for 'Source' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.L3.Target",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.Target' Error:Field validation for 'Target' failed on the 'required' tag"),
		},
		{
			caseDescription: "missing VM.L3.TC",
			vm: VM{
				Interface: &Interface{
					VxLAN: &VxLAN{
						VNI:    42,
						Source: &Iface{"x-42"},
						Target: &Iface{"vx-9a0201"},
						TC: &TC{
							Rate:  250,
							Burst: 256,
							Limit: 10240,
						},
					},
					L3: &L3{
						IPv4:   []string{"195.177.117.1"},
						IPv6:   []string{"2a02:2278:100:a1::1"},
						Upper:  &Iface{"vu-9a0201"},
						Source: &Iface{"vl-9a0201"},
						Target: &Iface{"if-9a0201"},
					},
					Uplink: &Iface{"bond-wan"},
				},
			},
			err: errors.New("Key: 'VM.Interface.L3.TC' Error:Field validation for 'TC' failed on the 'required' tag"),
		},
	}

	for _, testCase := range cases {
		err := Validate.Struct(testCase.vm)
		if err != nil && testCase.err != nil {
			if !strings.EqualFold(err.Error(), testCase.err.Error()) {
				t.Errorf("TestCase: %s\n Got : %s\n Want: %s\n ", testCase.caseDescription, err, testCase.err)
			}
		}
		if err == nil && testCase.err != nil {
			t.Errorf("TestCase: %s\n Got : nil\n Want: %s\n", testCase.caseDescription, testCase.err)
		}
		if err != nil && testCase.err == nil {
			t.Errorf("TestCase: %s\n Got : %s\n Want: nil", testCase.caseDescription, err)
		}
	}
}
