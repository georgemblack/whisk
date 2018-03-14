package main

import (
	"fmt"
	"github.com/georgemblack/whisk"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("Cleaning up...\n")
		whisk.Cleanup()
		os.Exit(0)
	}()

	whisk.Launch()
}
