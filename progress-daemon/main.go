package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"launchpad.net/unity-scope-snappy/progress-daemon/daemon"
	"log"
)

// main is the entry point of the daemon
func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	webdmAddressParameter := flag.String("webdm", "", "WebDM address[:port]")
	flag.Parse()

	daemon, err := daemon.New(*webdmAddressParameter)
	if err != nil {
		if *webdmAddressParameter == "" {
			log.Fatalf("Unable to create daemon: %s", err)
		} else {
			log.Fatalf(`Unable to create daemon with webdm API URL "%s": %s`, *webdmAddressParameter, err)
		}

	}

	err = daemon.Run()
	if err != nil {
		log.Printf("progress-daemon: Error running daemon: %s", err)
	}

	<-signals // Block here so the daemon can run, exiting if a signal comes in.
}
