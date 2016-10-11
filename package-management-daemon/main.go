package main

import (
	"launchpad.net/unity-scope-snappy/package-management-daemon/daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// main is the entry point of the daemon
func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	daemon, err := daemon.New()
	if err != nil {
		log.Fatalf("Unable to create daemon: %s", err)
	}

	err = daemon.Run()
	if err != nil {
		log.Printf("package-management-daemon: Error running daemon: %s", err)
	}

	<-signals // Block here so the daemon can run, exiting if a signal comes in.
}
