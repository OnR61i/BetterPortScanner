package pkg/scanner

import(
	"net"
	"sync.WaitGroup"
)

type Scanner struct{
	targetRange []net.IP
	portRange []int
	activeWorkers int
 	strategy Strategy 
	chan jobs Job
	chan results Job
	chan stopChan struct{}
	wg *sync.WaitGroup
	status string
}


func NewScanner(activeWorkers int) Scanner{

}

func initWorker(){
}

func (scanner *Scanner) startScan(target []net.IP, portRng []int, strategy Strategy) error{
	for _, ip := range target{
		for _, port := range portRng{
			if scanner.status != "stopping"{
				job := NewJob(ip, port)
				jobs <- job
				wg.Add(1)
			}
		}
	}
}

func (scanner *Scanner) stopScan(){
}

func (scanner *Scanner) resizeWorkerPool(newCount int){
}

// Required setter-methods...

func (scanner *Scanner) setStrategy(strategy Strategy){
}

