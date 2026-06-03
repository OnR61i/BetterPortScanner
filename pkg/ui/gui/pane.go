package main

import( 
	"gfx"
)

type Pane struct {
	id string 	// Naming-convention: use "Background" for the basic pane...
	x int
	y int
	height int
	width int

	rgb1 int
	rgb2 int
	rgb3 int
}

func NewPane(x int, y int, height int, width int, rgb1, rgb2, rgb3) Pane {
	return Pane {
		x:	x,
		y:	y,
		height:	height,
		width:	width,
		rgb1:	rgb1,
		rgb2:	rgb2,
		rgb3: 	rgb3,
	}
}

func (pane *Pane) Position() (int, int, int, int) {
	return pane.x, pane.y, pane.height, pane.width
}

func (pane *Pane) Identification() string {
	return pane.id
}

/*
func (pane *Pane) Rearange(x int, y int, height int, width int) {
	pane.x 		= x
	pane.y 		= y
	pane.height	= height
	pane.width	= width
}
*/
