package pkg/scanner

import(
	"log"
	"time"
	"sync"
	"net"

)

func main() {
	scanner := NewScanner(500)
	var strategy Strategy = &TCP{}
	intrfc, _ := net.InterfaceByName("eno1")
	wg := &sync.WaitGroup{}
	targets := []net.IP{net.IP{192, 168, 178, 38}, net.IP{192, 168, 178, 1}, net.IP{192, 168, 178, 46},}
	ports := []int{int(80), int(53), int(8009)}
	outcome, _ := scanner.Scan(targets, ports, strategy, intrfc, (11 * time.Second), wg)

	for oc := range outcome {
		log.Println("check")
		log.Println(oc.State)
		log.Println(oc)
	}
}
