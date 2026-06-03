package main

import(
	"gfx"
)

type Window struct {
}

type Button struct {
	y int
	x int 
	width int
	height int
}
	
func NewWindow() (Window) {
	return Window{
	}
}

func NewButton(y int, x int, width int, height int) (Button) {
	return Button{
		y: 		y,
		x: 		x,
		width: 		width, 
		height: 	height,
		}
}

func (window *Window) paintWindow() {
	// Paints the basic window...
	// --> Paints background...
	gfx.Stiftfarbe(31, 30, 30)
	gfx.Vollrechteck(
	0,
	0,
	900,
	600)
	
	// --> Paints sidebar...
	gfx.Stiftfarbe(51, 51, 51)
	gfx.Vollrechteck(
	0,
	0,
	30,
	600)
	
	// Boarderline...
	gfx.Stiftfarbe(69, 68, 68)
	gfx.Vollrechteck(
	30,
	0,
	2,
	600)
	
	// --> Paints Middle Screen divide...
	gfx.Stiftfarbe(69, 68, 68)
	gfx.Vollrechteck(
	450,
	0,
	2,
	600)
	
	// --> Paints topbar...
	gfx.Stiftfarbe(38, 38, 38)
	gfx.Vollrechteck(
	0,
	0,
	900,
	20)
	
	// Boarderline...
	gfx.Stiftfarbe(69, 68, 68)
	gfx.Vollrechteck(
	0,
	20,
	900,
	2)
	
	// --> Paints buttonpane...
	gfx.Stiftfarbe(101, 118, 163)
	gfx.Vollrechteck(
	0,
	490,
	300,
	110)
	
	// ------------------------------------------------
	// Buttons without function...
	y := uint16(10)
	x := uint16(532)
	
	
	// --> Ipnet...
	gfx.Stiftfarbe(107, 131, 242)
	gfx.Vollrechteck(
	y,
	(x - 26),
	280,
	20)
	
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Schreibe((y + 3), (x - 20), "IpNet / CIDR")
	
	// --> Portrange...
	gfx.Stiftfarbe(107, 131, 242)
	gfx.Vollrechteck(
	y,
	x,
	86,
	20)
	
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Schreibe((y + 3), (x + 6), "Port-Range")
	
	// --> Protocol...
	gfx.Stiftfarbe(107, 131, 242)
	gfx.Vollrechteck(
	107,
	x,
	86,
	20)
	
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Schreibe((y + 103), (x + 6), "Protocol")
	
	// --> Interface...
	gfx.Stiftfarbe(107, 131, 242)
	gfx.Vollrechteck(
	203,
	x,
	87,
	20)
	
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Schreibe((y + 197), (x + 6), "Interface")
	
	// --> Start...
	gfx.Stiftfarbe(107, 131, 242)
	gfx.Vollrechteck(
	y,
	565,
	280,
	20)
	
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Schreibe((y + 113), (x + 38), "Start")
}

func (button *Button) runListener(bool chan signalChan) {
	for {
		btn, stat, x, y := gfx.MausLesen1()
		if(btn == 1) {
			if(x > button.x && x < button.height && y > button.y && y > button.width){
				signalChan <- true
			}
		}
	}
}


func main(){
	gfx.Fenster(900, 600)
	window := NewWindow()
	paintWindow(window)
	
	for{
		}
}

