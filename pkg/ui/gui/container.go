package main

import (
	"gfx"
	"errors"
)

const FONT_SIZE = 5
const DEFAULT_BOARDER_DISTANCE = 3
const FONT_PIXEL_LEN = 

type Container struct {
	id string 		// Naming-convention: use "Topbar" & "Sidebar" for top- and sidebar in standart window...
	x int
	y int
	height int
	width int
	panes []Pane
	buttons []Button
	textfields []Textfield
}

func NewContainer(id string, x int, y int, height int, width int) Container {
	return Container {
		id:	id,	
		x: 	x,
		y:	y,
		height:	height,
		width:	width,
		}
}

func (container *Container) Position() (int, int, int, int) {
	return container.x, container.y, container.height, container.width
}

func (container *Container) Identification() string {
	return container.id
}

func (container *Container) Paint() {
	for _, pane := range panes {
		// Set color...
		rgb1, rgb2, rgb3 := pane.Color()
		gfx.Stiftfarbe(rgb1, rgb2, rgb3)

		// Paint...
		x, y, height, width := pane.Position()
		gfx.Vollrechteck(x, y, width, height)
	}
	
	for _, button := range buttons {
		// Set color + fonts...
		rgb1, rgb2, rgb3 := button.Text().Color() 
		gfx.Stiftfarbe(rgb1, rgb2, rgb3)	
		gfx.SetzeFont(button.Text().Font(), FONT_SIZE)

		// Paint...
		x, y, height, width := button.Position() 
		x, y = arrange(x, y, height, width, len(button.Text().Text()), button.Text().Arrangement())
		gfx.SchreibeFont(x, y, button.Text().Text()) 

	for _, textfield : range textfields {
		// Set color + fonts...
		rgb1, rgb2, rgb3 := textfield.Text().Color()
		gfx.Stiftfarbe(rgb1, rgb2, rgb3)
		gfx.SetzeFont(textfield.Text().Font() , FONT_SIZE)

		// Paint...
		x, y, height, width := textfield.Position()
		x, y = arrange(x, y, height, width, len(textfield.Text().Text()), textfield.Text().Arrangement())

		gfx.SchreibeFont(x, y, textfield.Text().Text())
	}
}

func (container *Container) Add(element any) (error) {
	switch element = element.(type) {
	case Pane:
		pane := element.(Pane)
		if(pane.id == nil){
			return errors.New("Pane has no id")
		}

		x, y, height, width := pane.Position()
		if (x == nil | y == nil | height == nil | width == nil | rgb1 == nil | rgb2 == nil | rgb3 == nil) {
			return errors.New("Not all positioning fields are declared")
		}
		if (x < container.x |
			y < container.y |
				(x + width) > (container.x + container.width) |
					y > (container.y + container.width) {

			return errors.New("Position out of bounds")
		}
		
		container.panes = append(container.panes, pane)
		return nil
	
	case Button:
		button := element.(Button)
		if(button.id == nil){
			return errors.New("Button has no id")
		}

		x, y, height, width := button.Position()
		if (x == nil | y == nil | height == nil | width == nil | rgb1 == nil | rgb2 == nil | rgb3 == nil) {
			return errors.New("Not all positioning fields are declared")
		}
		if (x < container.x |
			y < container.y |
				(x + width) > (container.x + container.width) |
					y > (container.y + container.width) {

			return errors.New("Position out of bounds")
		}
		if (button.Text == nil) {
			return errors.New("Text is not declared")
		}
		
		container.buttons = append(container.buttons, button)
		return nil

	case TextField:
		textfield := element.(TextField)
		if(textfield.id == nil){
			return errors.New("Textfield has no id")
		}
		
		x, y, height, width := button.Position()
		if (x == nil | y == nil | height == nil | width == nil | rgb1 == nil | rgb2 == nil | rgb3 == nil) {
			return errors.New("Not all positioning fields are declared")
		}
		if (x < container.x |
			y < container.y |
				(x + width) > (container.x + container.width) |
					y > (container.y + container.width) {

			return errors.New("Position out of bounds")
		}
		if (textField.Text == nil) {
			return errors.New("Text is not declared")
		}
		
		container.textfields = append(textfields, textfield)
		return nil
	}
}

func (container *Container) Delete(id string) (error) {
	for i, pane := range container.panes {
		if(pane.id == id) {
			// Remove pane from panes...
			part1 := container.panes[:i]
			part2 := container.panes[i + 1:]
			panes := append(part1, part2)

			container.panes = panes
			return
		}
	}

	for i, button := range container.buttons {
		if(button.id == id) {
			// Remove button from buttons...
			part1 := container.button[:i]
			part2 := container.button[i + 1:]
			buttons := append(part1, part2)

			container.buttons = buttons
			return
		}
	}

	for i, textfield := range container.textfield {
		if(textfield.id == id) {
			// Remove textfield from textfields...
			part1 := container.textfields[:i]
			part2 := container.textfields[i + 1:]
			textfields := append(part1, part2)

			container.textfields = textfields
			return
		}
	}

	return errors.New("No element with given id")
}

func arrange(x int, y int, height int, width int, textlen int, arrangement int) (int, int) {
	switch arrangement {
	case 0:
		return (x + DEFAULT_BOARDER_DISTANCE), (y + DEFAULT_BOADER_DISTANCE)
		
	case 1: 
		totalTextLen := textlen * FONT_PIXEL_LEN
		halfTextLen := totalTextLen / 2

		totalWidth := width 
		halfWidth := totalWidth / 2

		return (x + halfWidth - halfTextLen), (y + DEFAULT_BOARDER_DISTANCE)

	case 2:
		totalTextLen := textlen * FONT_PIXEL_LEN
		
		totalWidth := width

		return (x + totalWidth - totalTextLen), (y + DEFAULT_BOARDER_DISTANCE)
	}
}



