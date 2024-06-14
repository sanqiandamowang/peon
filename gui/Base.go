package gui

import (
	"log"

	"github.com/jroimartin/gocui"
)



func Dis_base() {
	var err error
	BaseGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer BaseGui.Close()

	BaseGui.Highlight = true
	// BaseGui.SelFgColor = gocui.ColorBlack
	// BaseGui.SelBgColor = gocui.ColorYellow

	BaseGui.SetManagerFunc(layout)

	if err := keybindings(BaseGui); err != nil {
		log.Panicln(err)
	}

	if err := BaseGui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("base", 0, 0, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Base"
		v.Wrap = true
	}

	if _, err := g.SetCurrentView("base"); err != nil {
		return err
	}
	err := pageMain()
	if err != nil {
		return err
	}
	return nil
}


func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
