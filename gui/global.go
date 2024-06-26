package gui

import (
	"peon/model/config"
	"peon/utils/json"

	"github.com/CaoYnag/gocui"
)

const (
	baseViwe = "base"
	infoView = "info"
	cmdView  = "cmd"
	debugView = "debug"
	errorView = "error"
	fileListView = "fileList"
	fileTreeView = "fileTree"
	fileEditView = "fileEdit"
)


var PeonGui *gocui.Gui
var cmdConfig config.Config
// var cursor_len = 0
var isDisPageDebug = true
var previousView []*string

func updatePreviousView(){
	if v := PeonGui.CurrentView(); v != nil {
		currentViewName :=v.Name()
		previousView=append(previousView, &currentViewName)
	}
	peonDebug("update previous view "+ *previousView[len(previousView)-1])
}

func LoadConfig() error {
	err := json.Read(configDir, &cmdConfig)
	if err != nil {
		return err
	}
	return nil
}
