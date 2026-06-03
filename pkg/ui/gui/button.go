package main

import (
	"gfx"
)

type Button struct {
	id string
	x int
	y int
	height int 
	width int
	chan isPressed bool
	text Text	
}

func NewButton(x int, y int, height int, width int, chan isPressed bool, text Text) Button {
	return Button {
		x:		x, 
		y:		y,
		height:		heigth,
		width:		width,
		text:		text,
	}
}

func (button *Button) Position() (int, int, int, int) {
	return button.x, button.y, button.height, button.width
}

func (button *Button) Text() Text {
	return Text
}

func (button *Button) Identification() string {
	return id
}

func (button *Button) Listener(chan isPressed bool) {
	button.isPressed = isPressed
}

func (button *Button) Listen() {
	for btn, stat, x, y := gfx.MausLesen1() {
		if btn == {	// Get to know number for left mouse btn... 
			if(x >= button.x && x <= (x + button.width) && y >= button.y && y <= (y + button.height)) {
				isPressed <- true
			}
		}
	}
}

/*
func (button *Button) Rearrange(x int, y int, height int, width int) {
	button.x 	= x
	button.y 	= y
	button.height 	= heigth
	button.width 	= width
}
*/
