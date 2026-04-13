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

func NewScanner() *Scanner{
}

func (scanner *Scanner) startScan(target []net.IP, portRng []int, strategy Strategy) error{
}

func (scanner *Scanner) stopScan(){
}

func initWorker(){
}

func (scanner *Scanner) resizeWorkerPool(newCount int){
}

// Required setter-methods...

func (scanner *Scanner) setStrategy(strategy Strategy){
}
