package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

func calc_portrange(min, max int) []int {
	var portrange []int
	for i := min; i <= max; i++ {
		portrange = append(portrange, i)
	}

	return portrange
}

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

func portsparser(ports []string) []int {
	var parsedports []int
	for _, item := range ports {
		if strings.Contains(item, "-") {
			splitrange := strings.Split(item, "-")
			min, _ := strconv.Atoi(splitrange[0])
			max, _ := strconv.Atoi(splitrange[len(splitrange)-1])
			portrange := calc_portrange(min, max)
			parsedports = append(parsedports, portrange...)
		} else {
			p, _ := strconv.Atoi(item)
			parsedports = append(parsedports, p)
		}
	}

	return parsedports
}

func main() {

	//trying to add a commandline flag so we can pass ports like -p 80, 100-109, 213 etc.
	var portFlag = flag.String("p", "ports", "ports to scan")

	flag.Parse()

	fmt.Println("portFlag has value ", *portFlag)

	parsedports := portsparser(strings.Split(*portFlag, ","))

	for _, v := range parsedports {
		fmt.Printf("port to scan %d\n", v)
	}

	ports := make(chan int, 100) //second parameter here makes this channel buffered. Indicates here that receiver can take 100 items before blocking.
	results := make(chan int)

	var openports []int //create a slice to store the openports

	for i := 0; i <= cap(ports); i++ { //spin up 100 workers each as a separate go routine
		go worker(ports, results)
	}

	go func() { //send to the workers in a separate thread, because the result gathering loop needs to start before more than 100 items of work can continue
		// for i := 1; i <= 1024; i++ {
		// 	ports <- i
		// }
		for _, i := range parsedports {
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
