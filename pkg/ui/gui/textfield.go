package main 

import( 
	"gfx"
)

type Textfield struct {
	id string
	x int
	y int
	height int
	width int
	text Text
}

func NewTextfield(x int, y int, height int, width int, text Text) Textfield {
	return Textfield{
		x:	x,
		y:	y,
		height:	height,
		width:	width,
		text:	text,
	}
}

func (textfield *Textfield) Position() (int, int, int, int) { 
	return textfield.x, textfield.y, textfield.height, textfield.width
}

func (textfield *Textfield) Identification() string {
	return id
}

func (textfield *Textfield) Text(text Text) {
	textfield.text = text
}

/*
func (textfield *Textfield) Rearrange(x int, y int, height int, width int) {
	textfield.x 	= x
	textfield.y	= y
	textfield.height= height
	textfield.width	= width
}
*/
