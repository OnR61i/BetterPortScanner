package scanner

import(
	"net"
	"sync"
)

type Scanner struct{
	jobs chan Job
	results chan Job
	interruptChan chan bool
	wg *sync.WaitGroup
	activeWorkers int
 	strategy Strategy 
	status string
}

func NewScanner(activeWorkers int) Scanner{
	return Scanner{
		jobs: 		make(chan Job),
		results:	make(chan Job),
		interruptChan: 	make(chan bool),
		wg: 		&sync.WaitGroup{},
		status:		"waiting",
	}
}

func (scanner *Scanner) Scan(target []net.IP, portRng []int, strategy Strategy) error{
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

func (scanner *Scanner) Interrupt(){
	for job := range jobs {
		wg.Done()
	}
	scanner.status = "waiting"	
}

func (scanner *Scanner) ResizeWorkerPool(newAmount int){
	close(scanner.interruptChan)
	for i := 0; i < newAmount; i++{
		go initWorker(i)
	}
}		

func (scanner *Scanner) initWorker(workerId int){
	for(workerId > activeWorkers){
		select {
		case interrupt := <- scanner.interruptChan:
			break	
		case job := <- scanner.jobs:
			err := scanner.strategy.execute(&job)
			if err != nil {
				continue
			}

			results <- job
		}
	}
}




	




