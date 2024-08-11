package gui

import (
	"fmt"
	"github.com/CaoYnag/gocui"
)

const configDir = "config/config.json"

func enterChooseItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	if cy == 0 {
		err = pageJsonChoose(g)
	}
	return err
}

func disMainpage(g *gocui.Gui) error {
	pageDebug(g)

	maxX, maxY := g.Size()
	if v, err := g.SetView(infoView, 0, 0, maxX/2, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		v.Wrap = true
		configDir := ""
		if cmdConfig.ConfigDir != nil {
			configDir = *cmdConfig.ConfigDir
		}
		fmt.Fprintln(v, "version:"+cmdConfig.Version+"  json_dir:"+configDir)
	}

	if v, err := g.SetView(cmdView, 0, 3, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "commands"
		v.Wrap = true
		fmt.Fprintln(v, "Edit JSON")
		for _, value := range cmdConfig.Plugins {
			fmt.Fprintln(v, value.Name+" :  "+value.Cmd)
		}
		if _, err := g.SetCurrentView(cmdView); err != nil {
			return err
		}
		peonDebug("system strat")
	}
	return nil
}

func pageMain(g *gocui.Gui) error {

	err := disMainpage(g)
	if err != nil {
		return err
	}

	return nil
}
