package gui

import "github.com/CaoYnag/gocui"

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(cmdView, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(cmdView, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(cmdView, gocui.KeyEnter, gocui.ModNone, enterChooseItem); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileListView, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileListView, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileListView, gocui.KeyArrowLeft, gocui.ModNone, returnPreviousView); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileListView, gocui.KeyEnter, gocui.ModNone, enterChooseFile); err != nil {
		return err
	}
	if err := g.SetKeybinding(errorView, gocui.KeyEnter, gocui.ModNone, closeErrorPopup); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeyArrowDown, gocui.ModNone, fileTreecursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeyArrowUp, gocui.ModNone, fileTreecursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeySpace, gocui.ModNone, updatefileTree); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeyArrowLeft, gocui.ModNone, fileTreeReturnPreviousView); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeyEnter, gocui.ModNone, changeView2FileEditView); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileTreeView, gocui.KeyCtrlS, gocui.ModNone, fileTreeSave); err != nil {
		return err
	}
	if err := g.SetKeybinding(fileEditView, gocui.KeyCtrlS, gocui.ModNone, changeView2FileTreeView); err != nil {
		return err
	}
	return nil
}
