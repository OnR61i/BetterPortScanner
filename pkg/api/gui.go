package api

import(
	"BetterPortScanner/pkg/scanner"
)

type GUI interface {
	Run() error
	// --> Starts the graphical user interface...

	Close()
	// --> Closes the graphical user interface...

	RefreshLiveView(string)
	// --> Allows ui bridge to print scan information on the second screen...

	UpdateScanStatus(string, float32)
	// --> Updates scan status and how far network has been scanned in percent...

	RefreshOutcomeList([]scanner.Job)
	// --> Updates screen for printout of final scan outcome...

	ShowError(error)
	// --> Displays errors while scanning or malformed inputs...
}


