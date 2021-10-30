package main

import (
	"flag"
	"fmt"
	"log"
	"main/flagvalue"
	scanner "main/scan"

	"golang.org/x/sync/semaphore"
)

var (
	address    string
	portNum    int
	customFlag flagvalue.CustomFlag
)

func main() {

	flag.StringVar(&address, "address", "", "address to scan")
	flag.IntVar(&portNum, "port", 0, "port for specified address")
	flag.Parse()
	stringInterval := flag.Args()

	numberInterval, err := flagvalue.CheckInterval(stringInterval, portNum)
	if err != nil {
		log.Fatal(err)
	}

	customFlag = *customFlag.NewCustomFlag()
	customFlag.SetAddress(address)
	customFlag.SetPort(portNum)
	customFlag.Interval.SetStart(numberInterval[0])
	customFlag.Interval.SetEnd(numberInterval[1])

	portScanner := &scanner.PortScanner{
		Ip:     customFlag.GetAddress(),
		Thread: semaphore.NewWeighted(5000),
	}
	if customFlag.GetPort() == 0 {
		portScanner.Start(customFlag.Interval.GetStart(), customFlag.Interval.GetEnd())
	} else {
		scanResult := portScanner.ScanPort("tcp", address, portNum)
		fmt.Println("address:", address, "port:", portNum, "is", scanResult.State)
	}
}
