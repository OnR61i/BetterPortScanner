package pkg/ui

import(
)

type Bridge struct{
	scanner *Scanner
	window *Window
	config *Config
	configFile ConfigFile
	totalJobs int
	currentResults []Job
}

func NewBridge() Bridge{
}

func ExecuteScan(targetRange string, PortRange string, strategy int) error{
}

func StopScan(){
}

func listenForResults(){
}

func calculateProgress() float{
}
