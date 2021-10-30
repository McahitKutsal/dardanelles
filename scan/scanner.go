package scan

import (
	"context"
	"fmt"
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
	Ip     string
	Thread *semaphore.Weighted
}

func (ps *PortScanner) ScanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: port}
	adress := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, adress, time.Second*1)
	if err != nil {
		result.State = "Closed"
		return result
	}
	defer conn.Close()
	result.State = "Open"
	return result
}

func ScanOpenPorts(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		//sistemin gönderilen istekleri karşılayamaması durumunda alınan hata
		//portu atlamamak için programı uyutup aynı porta tekrar istek atıyoruz
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanOpenPorts(ip, port, timeout)
		} else {
			//fmt.Println(port, "closed")
		}
		return
	}
	defer conn.Close()
	fmt.Println(port, "is open")
}

func (ps *PortScanner) Start() {
	//Go rutinlerin bitmesini beklemek için bir waitgroup nesnesi
	wg := sync.WaitGroup{}
	//wait group counter 0 olana kadar bekletir
	defer wg.Wait()
	for port := 1; port <= 6553; port++ {
		//5000 thread olarak belirlenmiş semafor nesnesine her port numarası için 1 context eklenir
		ps.Thread.Acquire(context.TODO(), 1)
		//wait group nesnesinin sayacına 1 ekliyoruz
		wg.Add(1)
		go func(port int) {
			ScanOpenPorts(ps.Ip, port, 2*time.Second)
			//işi biten semafor serbest kalıyor
			defer ps.Thread.Release(1)
			//wait group nesnesinin sayacından 1 çıkarıyoruz
			defer wg.Done()
		}(port)
	}
}
