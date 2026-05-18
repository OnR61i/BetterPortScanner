package ui

import(
	"pkg/scanner"
	"net"
	"errors"
	"sync"
)

type Bridge struct{
	scanner *Scanner
	window *Window
	config *Config
	configFile ConfigFile
	totalJobs int
	currentResults []Job
	wg *sync.WaitGroup
}

func NewBridge(scanner *scanner.Scanner, window *window, config *Config, configFile ConfigFile) Bridge{
	wg := sync.WaitGroup{}	
	return Bridge{
		scanner: 	scanner,
		window:	 	window,
		config:	 	config,
		configFile:	configFile,
		wg:		wg,
	}
}

func ExecuteScan(targetRange string, PortRange string, strategy int) error{
	// Checking if target range is a CIDR...
	ip, ipnet, err := net.parseCIDR(targetRange)
	if err != nil {
		// Checking if target is a single IP...
		ip := net.parseIP(targetRange) 
		if ip = nil{
			return errors.New("Given target range invalid") 
		}

		target := []net.IP{ip}
		scanner.Scan(target, []int, scanner.Strategy, &wg)
		return nil	
	}
	// Iterating over valid target range...	
	var targets []net.IP	
	for ip := ipnet.IP; ipnet.Contains(ip); inc(ip) {
		ipRange := append(ipRange, ip)
	}
	
	scanner.Scan(targets, []int, scanner.Strategy, &wg)
	return nil	
}

func StopScan(){
	scanner.Interrupt()	
}

func listenForResults(){
}

func calculateProgress() float{
}

// Function for iterating over IP range...
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j--{
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
} 
