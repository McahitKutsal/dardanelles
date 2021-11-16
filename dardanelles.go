package main

import (
	"flag"
	"fmt"
	"log"
	"main/flagvalue"
	"main/ports"
	scanner "main/scan"

	"golang.org/x/sync/semaphore"
)

var (
	address    string
	portNum    int
	customFlag flagvalue.CustomFlag
)

func main() {
	//Okunacak flag'lerin tanımlanması
	flag.StringVar(&address, "address", "", "address to scan")
	flag.IntVar(&portNum, "port", 0, "port for specified address")
	flag.Parse()
	stringInterval := flag.Args()
	//customFlag objesine okunan flagler set edilir varsa kuyruk set edilir
	customFlag = *customFlag.NewCustomFlag()
	customFlag.SetAddress(address)
	//check interval önemli bir fonksiyon
	numberInterval, err := flagvalue.CheckInterval(stringInterval, portNum)
	if err != nil {
		log.Fatal(err)
	}
	customFlag.SetPort(portNum)
	customFlag.Interval.SetStart(numberInterval[0])
	customFlag.Interval.SetEnd(numberInterval[1])
	//port scanner nesnesi oluşturulur
	portScanner := &scanner.PortScanner{
		Ip: customFlag.GetAddress(),
		// Semaphore eş zamanlı programlama ortamında kaynak yönetimi için kullanılan bir objedir
		// Bir task'tan diğerine sinyal göndermek için de kullanılır.
		Thread: semaphore.NewWeighted(5000),
		Up:     false,
	}
	if customFlag.GetPort() == 0 {
		//Kullanıcı bir port flag'i girmemişse veya bir aralık girmiş ise çalışır
		portScanner.Start(customFlag.Interval.GetStart(), customFlag.Interval.GetEnd())
		portScanner.ScanResult()
	} else {
		//Kullanıcı bir port flag'i girmiş ise çalışır
		scanResult := portScanner.ScanPort("tcp", address, portNum)
		fmt.Println("address:", address, "port:", portNum, "[", scanResult.State, "] ->", ports.PredictPort(portNum))
	}
}
