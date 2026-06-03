package main

import(
	"gfx"
)

type Text struct {
	text string
	font string	
	arrangement int 	// 0 = left, 1 = center, 2 = right...

	// Color...
	rgb1 int
	rgb2 int
	rgb3 int
}

func NewText(text string, font string, arrangement int, rgb1 int, rgb2 int, rgb3 int) Text {
	return Text {
		text:		text,
		font:		font,
		arrangement: 	arrangement,
		rgb1:		rgb1,
		rgb2:		rgb2,
		rgb3:		rgb3,
	}
}

func (text *Text) Text() string {
	return text.text
}

func (text *Text) Color() (int, int, int) {
	return text.rgb1, text.rgb2, text.rgb3
}

func (text *Text) Font() int {
	return text.font
}

func (text *Text) Arrangement() int {
	return text.arrangement
}

func (text *Text) Refresh(text string) {
	text.text = text
}
