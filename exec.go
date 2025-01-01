package main

import (
	"os/exec"
	"strings"
	"syscall"
)

// SanitizeInput - basic string input sanitisation
func SanitizeInput(s string) string {
	return strings.TrimSpace(strings.TrimSpace(strings.TrimRight(s, "\r\n")))
}

// RunCommandOutput i- output object for RunCommand
type RunCommandOutput struct {
	Command        string
	ReturnCode     int
	CombinedOutput []byte
}

// RunCommand - wrapper to run command
func RunCommand(name string, arg ...string) RunCommandOutput {
	var (
		outputObj RunCommandOutput
		err       error
	)

	// declare cmd object
	cmd := exec.Command(name, arg...)

	// set full command as string, for logging purposes
	if cmd.Args == nil || len(cmd.Args) == 0 {
		outputObj.Command = cmd.Path
	} else {
		outputObj.Command = strings.Join(cmd.Args, " ")
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}

	// execute command, acquire command return code
	outputObj.CombinedOutput, err = cmd.CombinedOutput()
	if err != nil {
		outputObj.ReturnCode = 254 // set undefined return code
		exitError, ok := err.(*exec.ExitError)
		if ok {
			status, ok := exitError.Sys().(syscall.WaitStatus)
			if ok {
				outputObj.ReturnCode = status.ExitStatus()
			}
		}
	}

	return outputObj
}
