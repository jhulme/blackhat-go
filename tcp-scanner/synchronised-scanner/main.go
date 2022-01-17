package main

import (
	"fmt"
	"net"
	"sync"
)

//This version of the scanner can fail due if excessive number of hosts or posts are scanned simultaneously
//To avoid this the next version implements a worker pool of goroutines to manage the concurrent scans.

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("%d open\n", j)
		}(i)
		wg.Wait()
	}
}
