package ui

import "BetterPortScanner/pkg/scanner"
import "log"

type Window struct{

	width int
	height int
	padding int
	title string
	isRunning bool

	targetRangeInput string
	portRangeInput string
	selectedStrategyIndex int
	workerCount int
	timeOutInput int

	progressPercent float32
	scanResults []scanner.Job
	statusMessage string
	isScanning bool
}


func NewWindow() Window{
	return Window{}
}
/*
func Run(){
}

func Init(){
}




func processInput(){
}

func handleKeyBoard(){
}

func draw() error{
}

func drawInputSection(input string) error{
}

func drawResultSection(input []Job) error{
}

func drawProgressIndicator(float) error{
}
func (window *Window) IsOpen() {}

*/

func (window *Window) IsOpen() (bool){
	return true
}
func (window *Window) Close(){
}

func (window *Window) UpdateScanStatus(string, float32){		// Input: float32, string
}

func (window *Window) RefreshOutcomeList(outcomeArr []scanner.Job){	// Input: []Job
	for _, outcome := range outcomeArr {
		log.Println(outcome)
	}
}

func (window *Window) RefreshLiveView(outcome scanner.Job){
	log.Println("Prints live view")
	log.Println(outcome)
}	// Input: Job

func (window *Window) ShowError(err error){	// Input: error
	log.Println(err)
}



