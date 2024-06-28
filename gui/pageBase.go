package gui

import (
	"log"

	"github.com/CaoYnag/gocui"
)

func DisBase() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	PeonGui = g
	defer g.Close()

	g.Cursor = true
	g.Highlight = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(baseViwe, 0, 0, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Base"
		v.Wrap = true
		if _, err := g.SetCurrentView(baseViwe); err != nil {
			return err
		}
	}

	err := pageMain(g)
	if err != nil {
		return err
	}

	return nil
}

