package pkg/scanner

import(
	"net"
	"sync"
	"time"
	"errors"
	//"log"
)

type Scanner struct{
	jobPool chan Job	// --> Current job pool used...
	worker int		// --> Count of active Workers...
}


func NewScanner(worker int) Scanner{
	return Scanner {
		worker: 	worker,
	}
}

func (scanner *Scanner) Scan(targetRange []net.IP, ports []int, strategy Strategy, intrfc *net.Interface, timeout time.Duration, wg *sync.WaitGroup) (chan Job, error) {
	// Initiates the network scan...
	jobPool := make(chan Job)
	outcome := make(chan Job)

	// --> starting worker routines...
	for i := 0; i < scanner.worker; i++ {
		go initWorker(jobPool, outcome) 
	}

	// --> Filling job pool...
	for _, ip := range targetRange {

		for _, port := range ports {
			jobPool <- NewJob(ip, port, strategy, intrfc, timeout)
			wg.Add(1)
		}
	}

	return outcome, nil
}

func (scanner *Scanner) Interrupt(){
	close(scanner.jobPool)
	scanner.jobPool = nil
}

func (scanner *Scanner) ResizeWorkerPool(newCount int){
	scanner.worker = newCount	
}

// Need better error handling for workers...
func initWorker(jobPool chan Job, outcome chan Job){
	for job := range jobPool {
		// --> Getting source ip address...
		srcIp, err := resolveSrcIp(job.UsedInterface)
		if err != nil {
			outcome <- Job{ State : "error" }
			continue
		}

		// --> Scanning target...
		isOpen, err := job.Strategy.Scan(job.Target, srcIp, job.Port, job.UsedInterface, job.Timeout)
		if err != nil { 
			outcome <- Job{ State : "error" }
			continue
		}

		// --> Returning outcome...
		if isOpen {
			job.State = "open"
			outcome <- job
		} else {
			// Has to be specified in the future...
			job.State = "Not Open"
			outcome <- job
		}
	}
}

// Code from: https://go.dev/play/p/BDt3qEQ_2H
func resolveSrcIp(intrfc *net.Interface) (net.IP, error) {
	// Missing error handling...
	addrs, _ := intrfc.Addrs()
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}

		return ip, nil 
    }

    return nil, errors.New("Didn't find ip address for given interface")
}


