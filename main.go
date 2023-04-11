package main

import (
	"flag"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	DEFAULT_MAX_PORT_RANGE = 65535 // 2^256
	DEFAULT_PORT           = 8080
	DEFAULT_HOST           = "scanme.nmap.org"
	DEFAULT_MAX_WORKERS    = 200
)

var (
	host = flag.String("host", DEFAULT_HOST, "host to be scanned(default=localhost)")
	totalWorkers = flag.Int("workers", DEFAULT_MAX_WORKERS, "max amount of workers (default=200)")
)

func CheckEnviroment() (string, error) {
	out, err := exec.Command("ping", "-c", "1", *host).Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`ttl=(.?).[\S]`)
	ttl := fmt.Sprintf("%s", re.FindString(string(out)))
	ttl = strings.Split(ttl, "=")[1]
	ttlNum, err := strconv.Atoi(ttl)
	if err != nil {
		return "", err
	}
	if ttlNum <= 64 {
		return "\n\t[+] Linux system\n", nil
	} else if ttlNum >= 127 {
		return "\n\t[+] Windows system\n", nil
	} else {
		return "\n\t[-] the time to the life of the target system doesn't exists\n", nil
	}
}

func workers(id int, jobs <-chan int, result chan<- string) {
	for jobs_port := range jobs {
		// fmt.Printf("Worker #%d started. Job(port): %d;\n", id, jobs_port)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, jobs_port))
		if err != nil {
			str := fmt.Sprintf("[ERROR] Port: %d failed to scan\n", jobs_port)
			result <- str
			continue
		}
		conn.Close()
		// fmt.Printf("Worker #%d finished. Job(port): %d;\n", id, jobs_port)
		str := fmt.Sprintf("[SUCCESS] Port: %d is Open\n", jobs_port)
		result <- str
	}
}

func main() {
	flag.Parse()

	env, _ := CheckEnviroment()
	fmt.Println(env)
	jobs := make(chan int, DEFAULT_MAX_PORT_RANGE)
	results := make(chan string, DEFAULT_MAX_PORT_RANGE)

	// load worker function
	for i := 0; i < *totalWorkers; i++ {
		go workers(i, jobs, results)
	}

	// send them to work
	for currentPort := 1; currentPort <= DEFAULT_MAX_PORT_RANGE; currentPort++ {
		jobs <- currentPort
	}
	close(jobs)

	// recover result from workers
	for currentPort := 1; currentPort <= DEFAULT_MAX_PORT_RANGE; currentPort++ {
		anwser := <-results
		if strings.Contains(anwser, "[SUCCESS]") {
			fmt.Println(anwser)
		}
	}
}
