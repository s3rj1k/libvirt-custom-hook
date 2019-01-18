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
			Logger.Printf("hook: '%s' prepare, begin -\n", os.Args[1])

			GracefullExit(c.PrepareBeginHook(domCfg))
		}

	// switch on: `qemu vm1 {start} begin -`
	case "start":

		// switch on: `qemu vm1 start {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' start, begin -\n", os.Args[1])

			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {started} begin -`
	case "started":

		// switch on: `qemu vm1 started {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' started, begin -\n", os.Args[1])

			GracefullExit(c.StartedBeginHook(domCfg))
		}

	// switch on: `qemu vm1 {stopped} end -`
	case "stopped":

		// switch on: `qemu vm1 stopped {end} -`
		if strings.EqualFold(os.Args[3], "end") {
			Logger.Printf("hook: '%s' stopped, end -\n", os.Args[1])

			GracefullExit(c.StoppedEndHook(domCfg))
		}

	// switch on: `qemu vm1 {release} end -`
	case "release":

		// switch on: `qemu vm1 release {end} -`
		if strings.EqualFold(os.Args[3], "end") {
			Logger.Printf("hook: '%s' release, end -\n", os.Args[1])

			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {migrate} begin -`
	case "migrate":

		// switch on: `qemu vm1 migrate {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' migrate, begin -\n", os.Args[1])

			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {restore} begin -`
	case "restore":

		// switch on: `qemu vm1 restore {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' restore, begin -\n", os.Args[1])

			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {reconnect} begin -`
	case "reconnect":

		// switch on: `qemu vm1 reconnect {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' reconnect, begin -\n", os.Args[1])

			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {attach} begin -`
	case "attach":

		// switch on: `qemu vm1 attach {begin} -`
		if strings.EqualFold(os.Args[3], "begin") {
			Logger.Printf("hook: '%s' attach, begin -\n", os.Args[1])

			GracefullExit(nil)
		}

	}

	GracefullExit(nil)
}
