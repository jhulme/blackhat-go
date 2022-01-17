package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100) //second parameter here makes this channel buffered. Indicates here that receiver can take 100 items before blocking.
	results := make(chan int)

	var openports []int //create a slice to store the openports

	for i := 0; i <= cap(ports); i++ { //spin up 100 workers each as a separate go routine
		go worker(ports, results)
	}

	go func() { //send to the workers in a separate thread, because the result gathering loop needs to start before more than 100 items of work can continue
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ { //results gathering loop receives on the results channel 1024 times, if the port doesn't equal 0 then append it to our slice
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports) //close off the ports and results threads
	close(results)

	sort.Ints(openports) //sort the final slice of open ports

	for _, port := range openports { //print out our sorted results
		fmt.Printf("%d open\n", port)
	}
}
