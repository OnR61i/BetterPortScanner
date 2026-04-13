package pkg/ui

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

	progressPercent float
	scanResults []Job
	statusMessage string
	isScanning bool
}

func NewWindow(configFile ConfigFile) Window{
}

func Run(){
}

func Init(){
}

func Close(){
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


func (window *Window) UpdateScanStatus(percent float, status string){
}

func (window *Window) RefreshResultList([]Job){
}

func (window *Window) ShowError(error){
}



