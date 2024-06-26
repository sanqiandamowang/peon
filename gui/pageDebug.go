package gui

import (
	"fmt"

	"github.com/CaoYnag/gocui"
)

func pageDebug(g *gocui.Gui) error {
	if !isDisPageDebug {
		return nil
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView(debugView, maxX/2+2, 0, maxX/2+50, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Autoscroll = true
		v.Title = "debug"
		v.Wrap = true
	}
	return nil
}
func peonDebug(msg string) {
	if !isDisPageDebug {
		return
	}
	v, err := PeonGui.View(debugView)
	if err != nil {
		return
	}
	fmt.Fprintln(v, msg)
	cursorDown(PeonGui, v)
}
func PeonDebugClear(msg string) {
	if !isDisPageDebug {
		return
	}
	v, err := PeonGui.View(debugView)
	if err != nil {
		return
	}
	v.Clear()
	fmt.Fprintln(v, msg)
	cursorDown(PeonGui, v)
}