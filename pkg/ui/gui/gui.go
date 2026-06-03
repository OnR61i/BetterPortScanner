package main

import(
	"gfx"
	"errors"
)

// Fonts...
const BASIC_FONT = ""

// Containers...
const TOP_BAR_ID = "Topbar"
const SIDE_BAR_ID = "Sidebar"
const TEXT_PANEL_ID = "Textpanels"

// Panes...
const BACKGROUND_ID = "Background"
const DIVIDING_LINE_ID = "Dividingline"

// Textfields...
const RESULT_PANEL_ID = "Resultpanel"
const LIVE_VIEW_PANEL_ID = "Liveview"

// Buttons...
const START_BUTTON_ID = "Startbutton"
const IP_TARGET_BUTTON_ID = "Iptargetbutton"
const PORT_RANGE_BUTTON_ID = "Portrangebutton"
const PROTOCOL_BUTTON_ID = "Protocolbutton"
const INTERFACE_BUTTON_ID = "Interfacebutton"

// (Button panes)
const START_BUTTON_PANE_ID = "Startbuttonpane"
const IP_TARGET_BUTTON_PANE_ID = "Iptargetbuttonpane"
const PORT_RANGE_BUTTON_PANE_ID = "Portrangebuttonpane"
const PROTOCOL_BUTTON_PANE_ID = "Protocolbuttonpane"
const INTERFACE_BUTTON_ID = "Interfacebuttonpane"

func Run(window *Window) error {
	// Initialize basic user interface...
	err = InitWindow(window)
	if err != nil {
		return err
	}

	// Get all button listeners...
	startBtn := make(chan bool)
	ipTargetBtn := make(chan bool)
	portRangeBtn := make(chan bool)
	protocolBtn := make(chan bool)
	interfaceBtn := make(chan bool)

	for _, container := range window.container {
		for _, button := range container.buttons {
			switch button.Identification() {
			case START_BUTTON_ID:
				button.Listener(startBtn)
				button.Listen()
				break	
			case IP_TARGET_BUTTON_ID:
				button.Listener(ipTargetBtn)
				button.Listen()
				break
			case PORT_RANGE_BUTTON_ID:
				button.Listener(portRangeBtn)
				button.Listen()
				break
			case PROTOCOL_BUTTON_ID:
				button.Listener(protocolBtn)
				button.Listen()
				break
			case INTERFACE_BUTTON_ID:
				button.Listener(interfaceBtn)
				button.Listen()
				break
			}
		}
	}

	// Initialize keyboard listener...
	keyboard := make(chan string)
	go window.ListenKeyboard(keyboard)
	
	// Stop chan for button events...
	stopEventChan := make(chan bool)

	// Initialize scan object with default values...
	scan := Scan {
		Target: 	"192.168.178.1",
		PortRange:	"1",
		Protocol:	"TCP",
		Interface:	"1",
	}

	for {
		select {
		case <- startBtn:
			
		case <- ipTargetBtn:
			go listenIpTarget(keyboard, stopEventChan, &scan, &window)
			break
		case <- portRangeBtn:
			go listenPortRange(keyboard, stopEventChan,  &scan, &window)
			break		
		case <- protocolBtn:
			go listenProtocol(keyboard, stopEventChan, &scan, &window)
			break
		case <- interfaceBtn:
			go listenInterface(keyboard, stopEventChan, &scan, &window)
			break
	}

func listenIpTarget(chan keyboard string, chan stopEventChan bool, scan *Scan, window *Window) {
	window.RefreshLiveView("Select ip-target(-range).: f.e. 192.168.178.113")
	for {
		select {
		case <- stopEventChan:
			return
		case char <- keyboard:
			scan.Target += char
			window.RefreshLiveView(scan.Target)
		}
	}
}

func listenPortRange(chan keyboard string, chan stopEventChan bool, scan *Scan, window *Window) {
	window.RefreshLiveView("Select ip-port(-range).: f.e. 80-129 or 80, 129...")
	for {
		select {
		case <- stopEventChan:
			return
		case char <- keyboard:
			scan.PortRange += char
			window.RefreshLiveView(scan.PortRange)
		}
	}
}

func listenProtocol(chan keyboard string, chan stopEventChan bool, scan *Scan, window *Window) {
	window.RefreshLiveView("Select protocol.: f.e. TCP or UDP")
	for {
		select {
		case <- stopEventChan:
			return
		case char <- keyboard:
			scan.Protocol += char
			window.RefreshLiveView(scan.Protocol)
		}
	}
}

func listenInteface(chan keyboard string, chan stopEventChan bool, scan *Scan, window *Window) {
	window.RefreshLiveView("Select interface (give number of interface).: f.e. 2")	
	for {
		select {
		case <- stopEventChan:
			return
		case char <- keyboard:
			scan.Interface += char
			window.RefreshLiveView(scan.Interface)
		}
	}
}

// Inefficient...
func (window *Window) RefreshLiveView(string) {
	for _, container := range window.container {
		for _, textfield := range container.textfields {
			if(textfield.Identification ==  LIVE_VIEW_PANEL_ID) {
				textfield.Text(Text{
					text:		string,
					font:		BASIC_FONT,
					arrangement:	0,
					rgb1:		0xff,
					rgb2:		0xff,
					rgb3:		0xff,
				}
			}
		}
	}
}

func InitWindow(window *Window) error {
	// Build topbar...

	topbar := NewContainer(
	// Identifiction...
	TOP_BAR_ID,	
	// Position...
	0,
	0,
	900,
	20
	)
	
	topbar.Add(NewPane(
		// Identification...
		BACKGROUND_ID,
		// Position...
		0,
		0,
		900,
		20,
		// Color...
		38, 
		38, 
		38
		)
	)
	
	// Add it...
	err := window.Add(topbar)
	if err != nil {
		return errors.New("Wasn't able to set topbar")
	}


	// Build sidebar...

	sidebar := NewContainer(
	// Identification...
	SIDE_BAR_ID,
	// Position...
	0,
	0,
	30,
	600
	)
		
	sidebar.Add(NewPane(
		// Identification...
		BACKGROUND_ID,
		// Position...
		0,
		0,
		30,
		600,		
		// Color...
		51, 
		51, 
		51
		)
	)
	
	// Add it...
	err := window.Add(topbar)
	if err != nil {
		return errors.New("Wasn't able to set sidebar")
	}


	// Build text panels (where results an live view are displayed)...

	textpanels := NewContainer(
	// Identification...
	TEXT_PANEL_ID,
	// Position...
	30,
	20,
	870,
	680,
	)

	textpanels.Add(NewPane(
		// Identification...
		BACKGROUND_ID,
		// Position...
		30,
		20,
		870,
		680,
		// Color...
		31,
		30,
		30
		)
	)

	textpanels.Add(NewPane(
		// Identification...
		DIVIDING_LINE_ID,
		// Position...
		450,
		0,
		2,
		600
		// Color...
		69, 
		68, 
		68
		)
	)

	textpanels.Add(NewTextfield(
		// Identification...
		RESULT_PANEL_ID,
		// Position...
		30, 
		20,
		450,
		490,
		Text{
			text: 		"(Final results are presented here)",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,
		},)
	)
	
	textpanels.Add(NewTextfield(
		// Identification...
		LIVE_VIEW_PANEL_ID,
		// Position...
		450, 
		20,
		900,
		600,
		// Text...
		Text{
			text: 		"(This is the live view)",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,
		})
	)

	// Add it...
	err = window.Add(panels)
	if err != nil {
		return errors.New("Wasn't able to set text panels")
	}

	
	// Build button panel...

	buttonpanel := NewContainer(
	// Identification...
	BUTTON_PANEL_ID,
	// Position...
	0, 
	490,
	303,	
	110,
	)

	buttonpanel.Add(NewPane(
		// Identification...
		BACKGROUND_ID,
		// Position...
		0,
		490,
		303,
		110,
		// Color...
		101, 
		118, 
		163
		)
	)
	
	// Startbutton background...
	buttonpanel.Add(NewButton(
		// Identification...
		START_BUTTON_PANE_ID,
		// Position...
		565,
		280,
		20,
		// Color...
		107, 
		131, 
		242
		)
	)

	//Startbutton function...
	buttonpanel.Add(NewButton(
		// Identification...
		START_BUTTON_ID,
		// Position...
		565,
		280,
		20,
		// Text...
		Text{
			text: 		"Start",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,

		},)
	)

	// IPtargetbutton background...
	buttonpanel.Add(NewButton(
		// Identification...
		IP_TARGET_BUTTON_PANE_ID,
		// Position...
	 	10,
		506,
		20,
		280,
		// Color...
		107, 
		131, 
		242
		)
	)
	
	// IPtargetbutton function...
	buttonpanel.Add(NewButton(
		// Identification...
		IP_TARGET_BUTTON_ID,
		// Position...
	 	10,
		506,
		20,
		280,
		// Text...
		Text{
			text: 		"IP-Target / CIDR",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,

		},)
	)

	// Portrangebutton background...
	buttonpanel.Add(NewButton(
		// Identification...
		PORT_RANGE_BUTTON_PANE_ID,
		// Position...
		10,
		532,
		20,
		86,
		// Color...
		107, 
		131, 
		242
		)
	)

	// Portrangebutton function...
	buttonpanel.Add(NewButton(
		// Identification...
		PORT_RANGE_BUTTON_ID,
		// Position...
		10,
		532,
		20,
		86,
		// Text...
		Text{
			text: 		"PortRange",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,

		},)
	)

	// Protocolbutton background...
	buttonpanel.Add(NewButton(
		// Identification...
		PROTOCOL_BUTTON_PANE_ID,
		// Position...
		107,
		532,
		20,
		86,
		// Color...
		107, 
		131, 
		242
		)
	)

	// Protocolbutton function...
	buttonpanel.Add(NewButton(
		// Identification...
		PROTOCOL_BUTTON_ID,
		// Position...
		107,
		532,
		20,
		86,
		// Text...
		Text{
			text: 		"Protocol",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,

		},)
	)
	
	// Interfacebutton background...
	buttonpanel.Add(NewButton(
		// Identification...
		INTERFACE_BUTTON_PANE_ID,
		// Position...
		203,
		532,
		20,
		86,
		// Color...
		107, 
		131, 
		242
		)
	)

	// Interfacebutton function...
	buttonpanel.Add(NewButton(
		// Identification...
		INTERFACE_BUTTON_ID,
		// Position...
		203,
		532,
		20,
		86,
		// Text...
		Text{
			text: 		"Interface",
			font: 		BASIC_FONT,
			arrangement:	0,
			rgb1:		0xff,
			rgb2:		0xff,
			rgb3:		0xff,

		},)
	)
	
	// Add it...
	err = window.Add(buttonpanel)
	if err != nil {
		return errors.New("Wasn't able to set button panel")
	}

	return nil
}

