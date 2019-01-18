package main

import (
	"log"
	"os"
	"strings"
)

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
			Logger.Println("graceful exit for libvirt, but errors occurred")
		} else {
			Logger.Println("graceful exit for libvirt, no errors occurred")
		}

		// exit with 0 code, else libvirt daemon will fail to start VM
		os.Exit(0)
	}

	// get Libvirt Domain XML as object
	domCfg, err := GetDomainXML(os.Stdin)
	if err != nil {
		GracefullExit(err)
	}

	switch os.Args[2] {

	// switch on: `qemu vm1 {prepare} begin -`
	case "prepare":

		// switch on: `qemu vm1 prepare {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Println("hook: prepare, begin")

			GracefullExit(c.PrepareBeginHook(domCfg))
		}

	// switch on: `qemu vm1 {started} begin -`
	case "started":

		// switch on: `qemu vm1 started {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Println("hook: started, begin")

			GracefullExit(c.StartedBeginHook(domCfg))
		}

	// switch on: `qemu vm1 {stopped} end -`
	case "stopped":

		// switch on: `qemu vm1 stopped {end} -`
		if strings.EqualFold(os.Args[3], "end") {
			Logger.Println("hook: stopped, end")

			GracefullExit(c.StoppedEndHook(domCfg))
		}
	}

	GracefullExit(nil)
}
