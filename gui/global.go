package gui

import (
	"fmt"
	"peon/model/config"
	"peon/utils/jsonutils"

	"github.com/CaoYnag/gocui"
)

const (
	version      = "0.1.0" //手动添加的版本号，用于后续版本更新
	baseViwe     = "base"
	infoView     = "info"
	cmdView      = "cmd"
	debugView    = "debug"
	errorView    = "error"
	fileListView = "fileList"
	fileTreeView = "fileTree"
	fileEditView = "fileEdit"
)

var PeonGui *gocui.Gui
var cmdConfig = &config.Config{
	Version:   version,
	ConfigDir: new(string),
	Plugins: []config.Plugin{
		//todo 待实现功能
		{
			Name: "plugin1",
			Cmd:  "ls",
		},
		{
			Name: "plugin2",
			Cmd:  "clear",
		},
	},
}

// var cursor_len = 0
var isDisPageDebug = true
var previousView []*string

func LoadConfig() error {
	err := jsonutils.Read(configDir, &cmdConfig)
	if err != nil {
		fmt.Println(err)
		fmt.Println("加载配置文件失败，使用默认配置")
		return nil
	}
	return nil
}

func updatePreviousView() {
	if v := PeonGui.CurrentView(); v != nil {
		currentViewName := v.Name()
		previousView = append(previousView, &currentViewName)
	}
	peonDebug("update previous view " + *previousView[len(previousView)-1])
}

func returnPreviousView(g *gocui.Gui, v *gocui.View) error {

	if v.Name() == cmdView {
		return nil
	}
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(*previousView[len(previousView)-1]); err != nil {
		return err
	}
	previousView = previousView[:len(previousView)-1]
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		l, err := v.Line(cy + 1)
		if err != nil {
			return err
		}
		if len(l) == 0 {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			line, _ := v.Line(cy)
			lx := len(line) - 3
			maxX, _ := v.Size()
			if lx > maxX {
				ox = lx - maxX
			}

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
			line, _ := v.Line(cy)
			lx := len(line) - 3
			maxX, _ := v.Size()
			if lx > maxX {
				ox = lx - maxX
			}

			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
