package gui

import (
	"peon/utils/json"
	"fmt"

	"github.com/jroimartin/gocui"
)

const configDir = "config/config.json"



func loadConfig() error {
	err := json.Read(configDir, &cmdConfig)
	if err != nil {
		fmt.Println("load config failed:", err)
		return err
	}
	// fmt.Println("load config success:", cmdConfig)
	return nil
}
func disMainpage() error {
	maxX, maxY := BaseGui.Size()
	if v, err := BaseGui.SetView("info", 0, 0, maxX/2, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		v.Wrap = true
		fmt.Fprintln(v, "version:"+cmdConfig.Version+ "   json_dir:  "+cmdConfig.ConfigDir)
	}
	if v, err := BaseGui.SetView("cmd", 0, 3, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "cmd"
		v.Wrap = true
		fmt.Fprintln(v, "Edit JSON")
		for _,value := range cmdConfig.Plugins {
			fmt.Fprintln(v, value.Name+" :  "+value.Cmd)
		}
	}
	if _, err := BaseGui.SetCurrentView("cmd"); err != nil {
		return err
	}
	return nil
}
func pageMain() error {
	err:=loadConfig()
	if err!= nil {
		return err
	}
	err=disMainpage()
	if err!= nil {
		return err
	}
	return nil
}
