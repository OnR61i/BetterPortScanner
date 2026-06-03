package ui

import(
	"BetterPortScanner/pkg/scanner"
	//"BetterPortScanner/pkg/api"
	"net"
	"errors"
	"sync"
	"strings"
	"strconv"
	"time"
	"fmt"
	"log"
)

type Bridge struct{
	scnr *scanner.Scanner
	window *Window
// 	config *Config
//	configFile ConfigFile
	windowScans chan Scan
	stopChan chan bool
	wg *sync.WaitGroup
	interfaceNames []string
}

type ScanParams struct{
	TargetRange []net.IP
	Ports []int
	Strategy scanner.Strategy
	UsedInterface net.Interface
	Timeout time.Duration
}

func NewBridge(scanner *scanner.Scanner, window *Window, windowScans chan Scan, stopChan chan bool) Bridge{		// Missing inputs: config *Config, configFile ConfigFile
	wg := &sync.WaitGroup{}
	return Bridge{
		scnr: 	scanner,
		window:	 	window,
//		config:	 	config,
//		configFile:	configFile,
		windowScans:	windowScans,
		stopChan:	stopChan,
		wg:		wg,
	}
}

const PORT_RANGE_ERROR_MSG = "Given port range invalid"

func (bridge *Bridge) Start() {
	// Main-Loop of whole application...
	// --> Starting graphical user interface...
	// Activate gui with: Run()
	// --> Waiting for tasks from window...
	for scan := range bridge.windowScans {
		// All actions when window task arrives...
		// --> Validate given params...
		params, err := validateScan(scan)

		if err != nil {
			log.Println(err)
			bridge.window.ShowError(errors.New("Given params not valid: Please check your spelling"))
			continue
		}

		// --> Send task to Scanner...
		outcomeChan, err := bridge.scnr.Scan(params.TargetRange, params.Ports, params.Strategy, &params.UsedInterface, params.Timeout, bridge.wg)
		log.Println(".")
		if err != nil {
			bridge.window.ShowError(err)
			continue
		}

		log.Println("Validation complete")
		// --> Setting up update channel for scan percentage...
		updateChan := make(chan bool)
		go bridge.updateScanStatus(params, updateChan)
		// --> Collect Outcome until Scanner terminates Scan...
		go func() {
			bridge.wg.Wait()
			log.Println("closing channel")
			outcomeChan <- scanner.Job{ State: "-1"}
			close(outcomeChan)
		}()

		var outcome []scanner.Job

		out:
		for {
			select {
			case <- bridge.stopChan:
				bridge.scnr.Interrupt()
				updateChan <- false
				break out

			case oc, ok := <- outcomeChan:
				// --> Breaks if channel was closed...
				if !ok {
					updateChan <- false
					log.Println("Breaks")
					break out
				}
				outcome = append(outcome, oc)
				updateChan <- true
				bridge.window.RefreshLiveView(oc)
				//log.Println(oc)
			}
		}

		log.Println("End of loop")
		// --> Sort Outcome and pass it to window...
		arrayedOutcome := arrayOutcome(outcome)
		log.Println("Stuck at array outcome")
		bridge.window.RefreshOutcomeList(arrayedOutcome)
	}

	// --> Main-Loop ends with closing window if not already done...
	bridge.window.Close()
}

func (bridge *Bridge) updateScanStatus(params ScanParams, updateChan chan bool) {
	// Constantly refreshes percentange of completed scans...
	bridge.window.UpdateScanStatus("scanning", 0.0)
	// --> Getting number of all scans in total...
	totalJobCount := 0
	for _, _ = range params.TargetRange {
		for _, _ = range params.Ports {
			totalJobCount++
		}
	}

	// --> For every completed Scan update percentage
	doneJobCount := 0
	for update := range updateChan {
		 if (update == false) {
			 bridge.window.UpdateScanStatus("waiting", 0.0)
			 break
		 }

		 // --> Calculate and shorten float...
		 k := doneJobCount / totalJobCount
		 i := fmt.Sprintf("%.2f", k)
		 f, _ := strconv.ParseFloat(i, 2)
		 bridge.window.UpdateScanStatus("scanning", float32(f))
	 }
 }

/*
func executeScan(params ScanParams, wg *sync.WaitGroup) (chan []scanner.Job, error) {
	// Starts scan...
	outcome, err := scanner.Scan(params.TargetRange, params.PortRange, params.Strategy, params.UsedInterface, params.Timeout, &wg)
	if err != nil {
		return nil, err
	}
}
*/

func validateScan(scan Scan) (ScanParams, error) {
	// Validates given params for scan...
	targetRangeStr := scan.TargetRange
	portRangeStr := scan.PortRange
	strategyStr := scan.Strategy
	intrfcStr := scan.NetInterface

	//log.Println("Start validation")
	var params ScanParams
	// 1. Checking value TargetRange...
	// --> Checking if target range is a CIDR...
	ip, ipnet, err := net.ParseCIDR(targetRangeStr)
	if err != nil {
		log.Println(err)
		// --> Checking if target is a single IP...
		ip := net.ParseIP(targetRangeStr)
		if ip == nil{
			return ScanParams{}, errors.New("Given target range invalid")
		}

		params.TargetRange = []net.IP{ip}

	} else {
		// --> Iterating over valid target range...
		var ipRange []net.IP = []net.IP{}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			ipCopy := make(net.IP, len(ip))
			copy(ipCopy, ip)
			// log.Println(ip)
			ipRange = append(ipRange, ipCopy)
			//log.Println(ipRange)
		}

		params.TargetRange = ipRange
		//log.Println(ipRange)

	}

	//log.Println("TargetRange valid")
	// 2. Checking value PortRange...
	if (strings.Contains(portRangeStr, ",")) {
		// --> Checking if PortRange is a list...
		portStrs := strings.Split(portRangeStr, ",")
		var ports []int
		for _, portStr := range portStrs {

			portStr = strings.TrimSpace(portStr)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}
			if (port > 65535 || port < 0) {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}

			ports = append(ports, port)

		}

		params.Ports = ports

	} else if (strings.Contains(portRangeStr, "-")) {
		// --> Checking if PortRange is an actual range...
		rangeStrs := strings.Split(portRangeStr, " ")
		var ports []int
		for _, rangeStr := range rangeStrs {
			vals := strings.Split(rangeStr, "-")

			// --> Validating first value...
			valStr1 := vals[0]
			valStr1 = strings.TrimSpace(valStr1)

			val1, err := strconv.Atoi(valStr1)
			if err != nil {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}
			if val1 > 65535 ||  val1 < 0 {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}

			// --> Validating second value...
			valStr2 := vals[1]
			valStr2 = strings.TrimSpace(valStr2)

			val2, err := strconv.Atoi(valStr1)
			if err != nil {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}
			if val2 > 65535 ||  val2 < 0 {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}

			// --> Validating range...
			if val2 >= val1 {
				return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
			}

			// --> Creating range...
			for i := val1; i <= val2; i++ {
				ports = append(ports, i)
			}
		}

		params.Ports = ports


	} else {
		// --> Checking if PortRange is single value...
		portRangeStr = strings.TrimSpace(portRangeStr)
		port, err := strconv.Atoi(portRangeStr)
		if err != nil {
			return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
		}
		if port > 65535 || port < 0 {
			return ScanParams{}, errors.New(PORT_RANGE_ERROR_MSG)
		}

		ports := []int{port}

		params.Ports = ports

	}
	// log.Println("PortRange valid")
	// 3. Checking Strategy...
	if (strings.EqualFold(strategyStr, "TCP")) {
		// --> Checking if strategy is TCP...
		params.Strategy = &scanner.TCP{}

	} else if (strings.EqualFold(strategyStr, "UDP")) {
		// --> Checking if strategy is UDP...
		params.Strategy = &scanner.UDP{}

	} else {
		return ScanParams{}, errors.New("Scanning strategy invalid")
	}

	// log.Println("Strategy valid")
	// 4. Checking chosen interface...
	interfaces, err := net.Interfaces()
	for _, intrfc := range interfaces {
		if (strings.EqualFold(intrfc.Name, intrfcStr)) {
			params.UsedInterface = intrfc
			if (err != nil) {
				return ScanParams{}, errors.New("Given interface invalid")
			}
		}
	}


	// log.Println("Interface valid")
	// 5. Checking timeout...
	timeout, err := time.ParseDuration(scan.Timeout)
	if err != nil {
		return ScanParams{}, errors.New("Given timeout invalid")
	}

	params.Timeout = timeout

	return params, nil

}

func arrayOutcome(unarrayedSlice []scanner.Job) []scanner.Job {
	// Sorts outcome chronologically...
	var arrayedOutcome []scanner.Job

	// --> Puts first reference in the beginning...
	x := unarrayedSlice[0]
	var xCopy scanner.Job
	xCopy = x
	arrayedOutcome = append(arrayedOutcome, xCopy)

	out:
	for _, x = range unarrayedSlice {
		// log.Println("Large iteration ---------------------------------------")
		// log.Println(x)

		for _, y := range arrayedOutcome {
			//log.Println(arrayedOutcome)
			//time.Sleep(time.Second)
			val, err := isGreaterThan(x.Target, y.Target)
			if err != nil {
				// Error Handling here...
			}
			switch val {
			case 1: // --> x is greater than y...
				// log.Println("Case1")
				break
			case 2: // --> x is smaller than y...
				// log.Println("Case2")
				var xCopy scanner.Job
				xCopy = x

				arrayedOutcome = append(arrayedOutcome, xCopy)
				// log.Println("Appended:", xCopy)
				continue out


			case 3: // --> x is equal to y...
				// log.Println("Case3")
				if y.Port > x.Port {
					break
				}
				if y.Port < x.Port {
					var xCopy scanner.Job
					xCopy = x

					arrayedOutcome = append(arrayedOutcome, xCopy)
					// log.Println("Appended:", xCopy)
					break
				}
				if y.Port == x.Port {
					break 	// --> Is an error, just ignoring it...
				}
			}
		}

		val, err := isGreaterThan(x.Target, arrayedOutcome[len(arrayedOutcome) - 1].Target)
		if err != nil {
			// Error Handling here...
		}
		if (val == 1) {
			// --> Checks if entry is greater than everything so far...
			// log.Println("Case4")
			var xCopy scanner.Job
			xCopy = x

			arrayedOutcome = append(arrayedOutcome, xCopy)
			// log.Println("Appended:", xCopy)
		}
	}
	// log.Println("Ready")
	return arrayedOutcome
}

// --> gfx-specific function...
func (bridge *Bridge) watchDog() {
	// --> Checks if Window is still running...
	for(bridge.window.IsOpen()){
		// --> Check is every fifth of a second...
		time.Sleep(200 * time.Millisecond)
	}

	// --> Interrupts Main-Loop...
	close(bridge.windowScans)
}

func inc(ip net.IP) {
	// Function for iterating over IP range...
	for j := len(ip) - 1; j >= 0; j--{
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

func isGreaterThan(ip1 net.IP, ip2 net.IP) (int, error) { // --> -1: error, 1: true, 2: false, 3: same
	if (len(ip1) != len(ip2)){
		return -1, errors.New("Passed Arrays have different length")
	}

	for j := len(ip1) - 1; j >= 0; j--{
		if (ip1[j] > ip2[j]) {
			return 1, nil

		} else if (ip1[j] == ip2[j]) {
			continue

		} else {
			return 2, nil
		}
	}

	return 3, nil
}


