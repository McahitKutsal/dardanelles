package main

import (
	"flag"
	"fmt"
	scanner "main/scan"

	"golang.org/x/sync/semaphore"
)

var (
	address string
	portNum int
)

func main() {

	flag.StringVar(&address, "address", "", "")
	flag.IntVar(&portNum, "port", 0, "")
	flag.Parse()

	portScanner := &scanner.PortScanner{
		Ip:     address,
		Thread: semaphore.NewWeighted(5000),
	}
	if address == "" {
		fmt.Println("Please specify an address with '--address' flag")
	} else if portNum == 0 {
		portScanner.Start()
	} else {
		scanResult := portScanner.ScanPort("tcp", address, portNum)
		fmt.Println("address:", address, "port:", portNum, "is", scanResult.State)
	}
}
