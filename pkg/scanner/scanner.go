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
	chan stopChan bool
	wg *sync.WaitGroup
	status string
}


func NewScanner(activeWorkers int) Scanner{
	return Scanner{
			jobs: 		make(chan, Job),
			results:	make(chan, Job),
			stopChan:	make(chan, bool),
			wg: 		*sync.WaitGroup{},
			status:		"waiting",
	}
}

func (scanner *Scanner) initWorker(workerId int){
	for(workerId >= scanner.activeWorkers){
		// Executing request...
	}
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

func (scanner *Scanner) resizeWorkerPool(newAmount int){

	formerAmount := scanner.activeWorkers

	scanner.activeWorkers = newAmount	
	
	// Creation count of workers are their id...	
	id := formerAmount	
	for(formerAmount < newAmount){
		id++
		go initWorker(id)
	}
}



