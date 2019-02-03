package sangocui

import (
	"github.com/jroimartin/gocui"
)

// KeyStrToCode :  Maps the string representation of a key to its gocui code
func KeyStrToCode(keyStr string) gocui.Key {
	switch keyStr {
	case "ctrlC":
		return gocui.KeyCtrlC
	case "arrowUp":
		return gocui.KeyArrowUp
	case "arrowDown":
		return gocui.KeyArrowDown
	case "arrowRight":
		return gocui.KeyArrowRight
	case "arrowLeft":
		return gocui.KeyArrowLeft
	case "ctrlD":
		return gocui.KeyCtrlD
	case "ctrlP":
		return gocui.KeyCtrlP
	case "ctrlSpace":
		return gocui.KeyCtrlSpace
	case "ctrlF":
		return gocui.KeyCtrlF
	case "ctrlB":
		return gocui.KeyCtrlB
	case "ctrlA":
		return gocui.KeyCtrlA
	case "enter":
		return gocui.KeyEnter
	case "ctrlL":
		return gocui.KeyCtrlL
	default:
		return gocui.KeyCtrl2
	}
}

// KeyStrToCode :  Maps the string representation of a color to its gocui code
func ColorStrToCode(color string) gocui.Attribute {
	switch color {
	case "green":
		return gocui.ColorGreen
	case "black":
		return gocui.ColorBlack
	case "white":
		return gocui.ColorWhite
	case "default":
		return gocui.ColorDefault
	default:
		return gocui.ColorBlack
	}
}
