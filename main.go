package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	fileName, nPorts := handleArgs()
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(nPorts)
		for port := 1; port <= nPorts; port++ {
			address := net.JoinHostPort(scanner.Text(), fmt.Sprintf("%d", port))
			go scanHost(address, &wg)
		}
		wg.Wait()
	}
}

func scanHost(address string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err == nil && conn != nil {
		defer conn.Close()
		fmt.Println(address)
	}
}

func handleArgs() (string, int) {
	var fileName string
	var nPorts int
	flag.StringVar(&fileName, "f", "", "file name that contains host addresses")
	flag.IntVar(&nPorts, "n", 1024, "number of ports to scan (starting from 1)")
	flag.Parse()
	if fileName == "" {
		fmt.Println("file name is required")
		os.Exit(1)
	}
	return fileName, nPorts
}
