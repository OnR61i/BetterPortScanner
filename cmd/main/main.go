package main


import(
	"BetterPortScanner/pkg/scanner"
	"BetterPortScanner/pkg/ui"
	//"pkg/config"
	//"log"
)

func main(){
	scnr := scanner.NewScanner(1000)
	window := ui.NewWindow()
	windowScans := make(chan ui.Scan)
	stopChan := make(chan bool)
	bridge := ui.NewBridge(&scnr, &window, windowScans, stopChan)
	go bridge.Start()

	scan := ui.NewScan(
	"192.168.178.1/24",
	"80",
	"TCP",
	"eno1",
	"1s")

	windowScans <- scan


	for {}
}
