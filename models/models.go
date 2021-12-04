package models

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var WindowWidth float32 = 800
var WindowHeight float32 = 800

var DefaultWidth float32 = 225
var DefaultHeight float32 = 40

var DefaultSize fyne.Size = fyne.Size{Width: DefaultWidth, Height: DefaultHeight}

var DefaultHorizontalPadding float32 = 10
var DefaultVerticalPadding float32 = 2

func CreateInput(placeHolder, defaultValue string) *widget.Entry {
	input := widget.NewEntry()
	input.SetPlaceHolder(placeHolder)
	input.SetText(defaultValue)
	input.Resize(DefaultSize)
	return input
}

func CreateButton(placeHolder string) *widget.Button {
	button := widget.NewButton(placeHolder, nil)
	button.Resize(DefaultSize)
	return button
}

func CreateLabel(placeHolder string) *canvas.Text {
	text := canvas.NewText(placeHolder, color.White)
	text.Alignment = fyne.TextAlignCenter
	text.Resize(DefaultSize)
	return text
}

func CreateRectangel(size fyne.Size, pos fyne.Position) *canvas.Rectangle {
	rect := canvas.NewRectangle(color.Opaque)
	rect.Resize(size)
	rect.Move(pos)
	return rect
}
