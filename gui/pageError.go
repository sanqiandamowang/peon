package gui

import (
	"fmt"

	"github.com/CaoYnag/gocui"
)

func closeErrorPopup(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	if len(previousView) > 1 {
		previousViewName := previousView[len(previousView)-1]
		g.SetCurrentView(*previousViewName)
		previousView = previousView[:len(previousView)-1]
	}
	return nil
}
func pageError(g *gocui.Gui, errMsg string) error {

	updatePreviousView()
	maxX, maxY := g.Size()
	if v, err := g.SetView(errorView, 0, 3, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Error"
		v.Wrap = true
		fmt.Fprintln(v, errMsg)
		fmt.Fprintln(v, "Press Enter to continue")

		if _, err := g.SetCurrentView(errorView); err != nil {
			return err
		}
	}

	return nil
}
