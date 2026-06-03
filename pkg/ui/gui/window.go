package main

import(
	"gfx"
)

type Window struct {
	height int
	width int
	container []Container
	chan keyboard string

}

func NewWindow(height int, width int) Window {
	return Window {
		height: 	height,
		width:		width,
	}
}

func (window *Window) Run() {
	gfx.Fenster(uint16(width), uint16(height))
}

func (window *Window) Closed() bool {
	return gfx.FensterOffen()
}

func (window *Window) Paint() {
	for _, container := range window.container {
		container.Paint()
	}
}

func (window *Window) Add(container Container) (error) {
	if(container.id == nil){
		return errors.New("Container has no id")
	}

	x, y, height, width := container.Position()
	if (x == nil | y == nil | height == nil | width == nil) {
		return errors.New("Not all positioning fields are declared")
	}
	if (height > window.height | width > window.width) {

		return errors.New("Position out of bounds")
	}
	
	window.container = append(window.container, container)
	return nil
}

func (window *Window) Delete(id string) (error) {
	for i, container := range window.container {
		if(container.id == id) {
			// Remove pane from panes...
			part1 := window.container[:i]
			part2 := window.container[i + 1:]
			container := append(part1, part2)

			window.container = container
			return
		}
	}

	return errors.New("No container with this id")
}

func (window *Window) Container(id string) (Container, error) {
	for _, container := range window.container{
		if(container.id == id) {
			return container, nil 
	}

	return nil, errors.New("No container with this id")
}

func listenKeyboard(chan keyboard string) {
	for {
		keyboard <- gfx.TastaturLesen1()
	}
}
