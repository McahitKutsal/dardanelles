package scan

import (
	"context"
	"fmt"
	"main/ports"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type ScanResult struct {
	Port  int
	State string
}
type PortScanner struct {
	Ip      string
	Thread  *semaphore.Weighted
	Up      bool
	latency time.Duration
}

func (ps *PortScanner) ScanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: port}
	adress := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, adress, time.Second*2)
	if err != nil {
		result.State = "Closed"
		return result
	}
	defer conn.Close()
	result.State = "Open"
	return result
}

func (ps *PortScanner) ScanOpenPorts(ip string, port int, timeout time.Duration, arr *[]string) {

	target := fmt.Sprintf("%s:%d", ip, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	latency := time.Duration(time.Since(start).Milliseconds())

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ps.ScanOpenPorts(ip, port, timeout, arr)
		} else {
			//fmt.Println(port, "closed")
		}
	} else {
		ps.Up = true
		ps.latency = latency
		defer conn.Close()
		s := "port:" + strconv.Itoa(port) + "[Open] -> " + ports.PredictPort(port)
		*arr = append(*arr, s)
	}

}

func (ps *PortScanner) Start(initial, final int) {
	wg := sync.WaitGroup{}
	var arr []string
	defer wg.Wait()
	for port := initial; port <= final; port++ {
		ps.Thread.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			ps.ScanOpenPorts(ps.Ip, port, 2*time.Second, &arr)
			defer ps.Thread.Release(1)
			defer wg.Done()
		}(port)

	}
	boolean := ps.Up
	if boolean {
		fmt.Println("\nServer is Up With", ps.latency, "Latency\n")
	}
	for _, s := range arr {
		fmt.Println(s)
	}

}
