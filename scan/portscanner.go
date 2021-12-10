package scan

import (
	"context"
	"fmt"
	"main/ports"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

//Sadece tek bir portu taramak için kullanılan struct (refaktör edilecek)
type ScanResult struct {
	Port  int
	State string
}

//PortScanner struct'ı geniş bir tarama yapılacağı zaman kullanılan struct'tır
type PortScanner struct {
	Ip      string
	Thread  *semaphore.Weighted
	Up      bool
	latency time.Duration
	results []string
}

//Tek bir port için kullanılan fonksiyon
func (ps *PortScanner) ScanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: port}
	adress := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, adress, time.Second*2)
	if err != nil {
		result.State = "Closed"
		return result
	}
	defer conn.Close()
	result.State = "Open"
	return result
}

//geniş tarama için kullanılan fonksiyon
func (ps *PortScanner) ScanOpenPorts(ip string, port int, timeout time.Duration, c bool) {

	target := fmt.Sprintf("%s:%d", ip, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	latency := time.Duration(time.Since(start).Milliseconds())

	//Eğer sunucu too many open files hatası verirse portu atlamamak için timeout kadar bekleyip aynı port'a yeniden istek atıyoruz
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ps.ScanOpenPorts(ip, port, timeout, c)
		} else {
			if c {
				s := "port: " + strconv.Itoa(port) + " [Closed] -> " + ports.PredictPort(port)
				ps.results = append(ps.results, s)
			}
		}
	} else {
		ps.Up = true
		ps.latency = latency
		defer conn.Close()
		//açık olan port default değeri tahmin edilerek results slice'ına eklenir
		s := "port: " + strconv.Itoa(port) + " [Open] -> " + ports.PredictPort(port)
		ps.results = append(ps.results, s)
	}

}

func (ps *PortScanner) Start(initial, final int, c bool) {
	fmt.Println("Scanning ports...")
	//Taranacak ip adresinin portları için bir bekleme grubu oluşturulur
	wg := sync.WaitGroup{}
	//Bekleme grubunun sayacı sıfıra ulaşmadan fonksiyonun dışına çıkılmaz
	// yani döngü bitse bile eşzamanlı fonksiyonlar bitmediği için sonuçlar ekrana yazıdırlmaz
	defer wg.Wait()
	for port := initial; port <= final; port++ {
		//Bir adet semaphore açılır
		ps.Thread.Acquire(context.TODO(), 1)
		//Taranacak her bir port numarası için bekleme grubunun sayacı 1 arttırılır
		wg.Add(1)
		//Aşağıdaki fonksiyon her çalıştırıldığında for döngüsünden kopar ve programın geri kalanının akışından bağımsız çalışır
		go func(port int) {
			//parametre ile alınan port numarası taranır
			ps.ScanOpenPorts(ps.Ip, port, 2*time.Second, c)
			//yukarıda açılan semaphore yeni bir port taramak için serbest bırakılır
			defer ps.Thread.Release(1)
			//bekleme grubunun sayacı 1 azaltılır
			defer wg.Done()
		}(port)

	}
	//konsol temizlenir
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//sonuclar yazdırılır
func (ps *PortScanner) ScanResult() {
	if ps.Up {
		fmt.Println("\nServer is Up With", ps.latency, "Latency\n")
	}
	for _, s := range ps.results {
		fmt.Println(s)
	}
}
