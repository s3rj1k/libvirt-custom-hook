package main

import (
	"fmt"
	"strconv"
)

// ConfigureTrafficControlOnInterface - enables TC magic on specified interface
func ConfigureTrafficControlOnInterface(rate, burst, limit int64, dev string) error {

	// prefix for errors logging
	const errPrefix = "tc config error:"

	// remove old TC config
	cmd := RunCommand("tc", "qdisc", "del", "dev", SanitizeInput(dev), "root")
	if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: No such file or directory
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// set tbf qdisk
	cmd = RunCommand("tc", "qdisc", "add", "dev", SanitizeInput(dev),
		"root", "handle", "1:", "tbf",
		"rate", fmt.Sprintf("%dmbit", rate),
		"burst", fmt.Sprintf("%dkb", burst),
		"limit", strconv.FormatInt(limit, 10))
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	// set fq_codel
	cmd = RunCommand("tc", "qdisc", "add", "dev", SanitizeInput(dev), "parent", "1:1", "handle", "10:", "fq_codel")
	if cmd.ReturnCode != 0 {
		e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		Logger.Println(e)

		return e
	}

	return nil
}
