package gui

import (
	"peon/model/config"

	"github.com/jroimartin/gocui"
)

var BaseGui *gocui.Gui = nil
var cmdConfig config.Config