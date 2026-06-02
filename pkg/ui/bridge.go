package pkg/ui

import(
	"pkg/scanner"
	"pkg/api"
	"net"
	"errors"
	"sync"
	"strings"
	"strconv"
)

type Bridge struct{
	scanner *Scanner
	window *Window
	config *Config
	configFile ConfigFile
	WindowScans chan Scan
	StopChan chan bool
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

const PORT_RANGE_ERROR_MSG = "Given port range invalid"

func Start() {
	// Main-Loop of whole application...
	// --> Starting graphical user interface...
	Run()

	// --> Waiting for tasks from window...
	for scan range WindowScans {
		// All actions when window task arrives...
		// --> Validate given params...
		params, err := validateScan(scan)
		if err != nil {
			window.ShowError(errors.New("Given params not valid: Please check your spelling"))
			continue
		}

		// --> Send task to Scanner...
		outcomeChan, err = scanner.executeScan(params, &wg)
		if err != nil {
			window.ShowError(err)
		}
		
		// --> Setting up update channel for scan percentage...
		updateChan := make(chan bool)
		go updateScanStatus(params, updateChan)
		// --> Collect Outcome until Scanner terminates Scan...
		go func() { 
			wg.Wait()
			close(outcome)
		}()
		
		outcome := []Job
		for { 
			select {
			case <- StopChan:
				scanner.Interrupt()
				updateChan <- false
				continue

			case oc := <- outcome:
				// --> Breaks if channel was closed...
				if oc == nil {
					updateChan <- false
					break
				}
				outcome = append(outcome, oc)
				updateChan <- true
				window.RefreshLiveView(oc)
			}
		}
		
		// --> Sort Outcome and pass it to window...
		arrayedOutcome := arrayOutcome(outcome) 
		window.RefreshOutcomeList(arrayedOutcome)	
	}
	
	// --> Main-Loop ends with closing window if not already done...
	window.Close()
}

func updateScanStatus(params ScanParams, updateChan chan bool) {
	// Constantly refreshes percentange of completed scans...
	window.UpdateScanStatus("scanning", 0.0)
	// --> Getting number of all scans in total...
	totalJobCount := 0
	for _, ip := range params.TargetRange {
		for _, port := range params.PortRange {
			totalJobCount++
		}
	}

	// --> For every completed Scan update percentage
	doneJobCount = 0
	for update := range updateChan {
		 if (update == false) {
			 window.UpdateScanStatus("waiting", 0,0)
			 break
		 }

		 // --> Calculate and shorten float...
		 k := doneJobCount / totalJobCount
		 i := fmtSprintf("%.2f", k)
		 f, _ := strconv.ParseFloat(i, 2)
		 window.UpdateScanStatus("scanning", f)
	 }
 }


func executeScan(params ScanParams, wg *sync.WaitGroup) (chan []Job, error) {
	// Starts scan...
	outcome, err := scanner.InitiateScan(params.TargetRange, params.PortRange, params.Strategy, params.UsedInterface, params.Timeout, &wg)
	if err != nil {
		return nil, err
	}
}

func validateScan(scan Scan) (ScanParams, error) {
	// Validates given params for scan...
	targetRangeStr := scan.TargetRange
	portRangeStr := scan.PortRange
	strategyStr := scan.Strategy
	intrfcStr := scan.NetInterface

	var params ScanParams
	// 1. Checking value TargetRange...
	// --> Checking if target range is a CIDR...
	ip, ipnet, err := net.parseCIDR(targetRange)
	if err != nil {
		// --> Checking if target is a single IP...
		ip := net.parseIP(targetRange) 
		if ip = nil{
			return nil, errors.New("Given target range invalid") 
		}
		
		params.TargetRange = []net.IP{ip}

	} else {
		// --> Iterating over valid target range...	
		var targetRange []net.IP	
		for ip := ipnet.IP; ipnet.Contains(ip); inc(ip) {
			ipRange := append(ipRange, ip)
		}

		params.TargetRange = targetRange
	}

	// 2. Checking value PortRange...
	if (strings.Contains(portRangeStr, ",")) {
		// --> Checking if PortRange is a list...
		portStrs := strings.Split(portRangeStr, ",")
		var ports []int
		for _, portStr := range portStrs {

			portStr = strings.TrimSpace(portStr)
			port, err := strconv.Atoi(portStr)	
			if err != nil {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}
			if port > 65535 | port < 0 {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}
			
			ports = append(ports, port)

		}

		params.Ports = ports

	} else if (strings.Contains(portRangeStr, "-")) {
		// --> Checking if PortRange is an actual range...
		rangeStrs := strings.Split(portRangeStr, " ")
		var ports []int
		for _, rangeStr := rangeStrs {
			vals := strings.Split(rangeStr, "-")

			// --> Validating first value...
			valStr1 := vals[0]
			valStr1 := strings.TrimSpace(valStr1)

			val1, err := strconv.Atoi(valStr1)
			if err != nil {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}
			if val1 > 65535 |  val1 < 0 {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}

			// --> Validating second value...
			valStr2 := vals[1]
			valStr2 := strings.TrimSpace(valStr2)
			
			val1, err := strconv.Atoi(valStr1)
			if err != nil {
				return errors.New(PORT_RANGE_ERROR_MSG)
			}
			if val2 > 65535 |  val2 < 0 {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}
			
			// --> Validating range...
			if val2 >= val1 {
				return nil, errors.New(PORT_RANGE_ERROR_MSG)
			}
			
			// --> Creating range...
			for i := val1; i <= val2; i++ {
				ports = append(ports, i)
			}
		}

		params.Ports = ports

	} else {
		// --> Checking if PortRange is single value...
		portRangeStr = strings.TrimSpaces(portRangeStr)
		port, err = strconv.Atoi(portRangeStr)
		if err != nil {
			return nil, errors.New(PORT_RANGE_ERROR_MSG)
		}
		if port > 65535 | port < 0 {
			return nil, errors.New(PORT_RANGE_ERROR_MSG)
		}

		ports := []int{port}

		params.Ports = ports

	}
	
	// 3. Checking Strategy... 
	if (strings.EqualFold(strategyStr, "TCP")) {
		// --> Checking if strategy is TCP...
		params.Strategy = scanner.TCP

	} else if (strings.EqualFold(strategyStr, "UDP")) {
		// --> Checking if strategy is UDP...
		params.Strategy = scanner.UDP
	
	} else {
		return nil, errors.New("Scanning strategy invalid")
	}

	// 4. Checking chosen interface...
	for _, intrfc := range interfaceNames {
		if (strings.EqualFold(intrfc, intrfcStr) {
			params.UsedInterface := net.InterfaceByName(intrfc)
		}
	}

	if (params.UsedInterface == nil) {
		return nil, errors.New("Given interface invalid")
	}

	// 5. Checking timeout...
	timeout, err := time.ParseDuration(scan.Timeout)
	if err != nil {
		return nil, errors.New("Given timeout invalid")
	}

	params.Timeout = timeout

	return params, nil

}

func arrayOutcome(unarrayedSlice []Job) []Job {
	// Sorts outcome chronologically...	
	var arrayedOutcome []Job
	for _, x := range unarrayedSlice {	
		if arrayedOutcome == nil {
			// --> Puts first reference in the beginning...
			arrayedOutcome = append(arrayedOutcome, x)
		}
		for _, y := range arrayOutcome {
			switch isGreaterThan(x, y) {
			case 1: // --> x is greater than y...
				continue

			case 2: // --> x is smaller than y...
				arrayedOutcome = append(arrayedOutcome, x)

			case 3: // --> x is equal to y...
				if y.port > x.port {
					continue
				} 
				if y.port < x.port {
					arrayedOutcome = append(arrayedOutcome, x)
				}
				if y.port == x.port {
					break 	// --> Is an error, just ignoring it...
				}
			}
		}

		if (isGreaterThan(x, [len(arrayOutcome) - 1]arrayOutcome)) {
			// --> Checks if entry is greater than everything so far...
			arrayedOutcome = append(arrayedOutcome, x)
		}
}

// --> gfx-specific function...
func watchDog() {
	// --> Checks if Window is still running...
	for(window.FensterOffen()){
		// --> Check is every fifth of a second...
		time.Sleep(200 * time.Millisecond)
	}

	// --> Interrupts Main-Loop...
	close(windowScans)
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

func isGreaterThan(ip1, ip2) (int, error) { // --> -1: error, 1: true, 2: false, 3: same
	if (len(ip1) != len(ip2)){
		return -1, errors.New("Passed Arrays have different length")
	}

	for i := len(ip1) - 1; j >= 0, j--{
		if [j]ip1 > [j]ip2 {
			return 1, nil
		
		} else if [j]ip1 == [j]ip2 {
			continue
		
		} else {
			return 2
		}
	}

	return 3, nil
}
