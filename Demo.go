// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/jroimartin/gocui"
// )

// var (
// 	viewArr = []string{"view1", "view2", "view3"}
// 	currentViewIndex = 0
// )

// func main() {
// 	g, err := gocui.NewGui(gocui.OutputNormal)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer g.Close()

// 	g.Highlight = true
// 	g.SelFgColor = gocui.ColorBlack
// 	g.SelBgColor = gocui.ColorYellow

// 	g.SetManagerFunc(layout)

// 	if err := keybindings(g); err != nil {
// 		log.Panicln(err)
// 	}

// 	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
// 		log.Panicln(err)
// 	}
// }

// func layout(g *gocui.Gui) error {
// 	maxX, maxY := g.Size()

// 	if v, err := g.SetView("view1", 0, 0, maxX/2-1, maxY/2-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "View 1"
// 		v.Wrap = true
// 		fmt.Fprintln(v, "This is view 1. Press 'i' to edit, Tab to switch views.")
// 	}

// 	if v, err := g.SetView("view2", maxX/2, 0, maxX-1, maxY/2-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "View 2"
// 		v.Wrap = true
// 		fmt.Fprintln(v, "This is view 2. Press 'i' to edit, Tab to switch views.")
// 	}

// 	if v, err := g.SetView("view3", 0, maxY/2, maxX-1, maxY-1); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Title = "View 3"
// 		v.Wrap = true
// 		fmt.Fprintln(v, "This is view 3. Press 'i' to edit, Tab to switch views.")
// 	}

// 	if _, err := g.SetCurrentView(viewArr[currentViewIndex]); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func keybindings(g *gocui.Gui) error {
// 	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {	
// 		return err
// 	}
// 	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("", 'i', gocui.ModNone, enterEditMode); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, exitEditMode); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func nextView(g *gocui.Gui, v *gocui.View) error {
// 	currentViewIndex = (currentViewIndex + 1) % len(viewArr)
// 	_, err := g.SetCurrentView(viewArr[currentViewIndex])
// 	return err
// }

// func enterEditMode(g *gocui.Gui, v *gocui.View) error {
// 	v.Editable = true
// 	g.Cursor = true
// 	return nil
// }

// func exitEditMode(g *gocui.Gui, v *gocui.View) error {
// 	v.Editable = false
// 	g.Cursor = false
// 	return nil
// }

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
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

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
		return err
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	vb := v.ViewBuffer()
	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, 20); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprint(v, "\rWill be")
		fmt.Fprint(v, "deleted\rItem 4\nItem 5")
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func Demo() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}