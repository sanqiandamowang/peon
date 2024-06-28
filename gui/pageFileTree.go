package gui

import (
	"errors"
	"fmt"
	"peon/utils/jsonutils"
	"peon/utils/kvtree"
	// "strconv"
	"strings"

	"github.com/CaoYnag/gocui"
	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
)

const (
	TreeSignDash     = "─"
	TreeSignVertical = "│"
	TreeSignUpMiddle = "├"
	TreeSignUpEnding = "└"
)

var KVTree kvtree.KV_Tree
var treeIndex int = 0
var treeIndexMax int = 0

func fileTreeReturnPreviousView(g *gocui.Gui, v *gocui.View) error {

	if v.Name() == cmdView {
		return nil
	}
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	if err := g.DeleteView(fileEditView); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(*previousView[len(previousView)-1]); err != nil {
		return err
	}
	previousView = previousView[:len(previousView)-1]
	return nil
}

func fileTreecursorDown(g *gocui.Gui, v *gocui.View) error {
	err := cursorDown(g, v)
	if err != nil {
		return err
	}
	if treeIndex < treeIndexMax-1 {
		treeIndex += 1
		updatefileEditView(g)
	}
	return nil
}

func fileTreecursorUp(g *gocui.Gui, v *gocui.View) error {

	err := cursorUp(g, v)

	if err != nil {
		return err
	}
	if treeIndex > 0 {
		treeIndex -= 1
		updatefileEditView(g)
	}
	return nil
}

func updatefileTree(_ *gocui.Gui, v *gocui.View) error {
	var err error = nil
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	KVTree.DisNodeList[treeIndex].IsExpand = !KVTree.DisNodeList[treeIndex].IsExpand
	v.Clear()
	KVTree.DisNodeList = nil
	treeIndexMax = 0
	printKVNode(v, KVTree.NodeList, "", true)
	if err := v.SetCursor(cx, cy); err != nil && treeIndex > 0 {
		if err := v.SetOrigin(ox, oy); err != nil {
			return err
		}
	}
	return err
}

func printKVNode(v *gocui.View, node *kvtree.KV_Node, indent string, isLast bool) {
	if node == nil {
		return
	}

	// Choose the appropriate prefix for the current node
	prefix := TreeSignUpMiddle
	if isLast {
		prefix = TreeSignUpEnding
	}

	key := ""
	if node.Source != nil {
		for k := range *node.Source {
			key = k
			break
		}
	}
	if node.Child != nil {
		key += "->"
	}
	buf := fmt.Sprintf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)
	fmt.Fprint(v, buf)
	treeIndexMax += 1
	peonDebug(fmt.Sprint(&node.Source))
	KVTree.DisNodeList = append(KVTree.DisNodeList, node)
	peonDebug(fmt.Sprint(&KVTree.DisNodeList[len(KVTree.DisNodeList)-1].Source))
	newIndent := indent
	if !isLast {
		newIndent += TreeSignVertical + " "
	} else {
		newIndent += " "
	}
	if node.IsExpand {
		if node.Child != nil {
			printKVNode(v, node.Child, newIndent, node.Child.Next == nil)
		}
	}

	if node.Next != nil {
		printKVNode(v, node.Next, indent, node.Next.Next == nil)
	}
}

func updatefileEditView(g *gocui.Gui) error {

	v, err := g.View(fileEditView)
	if err != nil {

		return nil
	}
	v.Clear()
	if KVTree.DisNodeList[treeIndex].Source != nil {
		for _, value := range *KVTree.DisNodeList[treeIndex].Source {
			switch value.(type) {
			case map[string]interface{}:
				buf, err := sonic.ConfigDefault.MarshalIndent(value, "", "  ")
				if err != nil {
					return err
				}
				_buf, err := jsonutils.SortJSON(string(buf))
				if err != nil {
					return err
				}
				fmt.Fprint(v, string(_buf))
			case []interface{}:
				buf, err := sonic.ConfigDefault.MarshalIndent(value, "", "  ")
				if err != nil {
					return err
				}
				_buf, err := jsonutils.SortJSON(string(buf))
				if err != nil {
					return err
				}
				fmt.Fprint(v, string(_buf))
			default:
				fmt.Fprint(v, value)
			}
		}
	}
	return nil
}

func changeView2FileEditView(g *gocui.Gui, _ *gocui.View) error {
	if _, err := g.SetCurrentView(fileEditView); err != nil {
		return err
	}
	v, err := g.View(fileEditView)
	if err != nil {
		return err
	}
	if err := v.SetCursor(0, 0); err != nil {
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
	}
	return nil
}

func changeView2FileTreeView(g *gocui.Gui, _ *gocui.View) error {

	if err := saveEditFile(g); err != nil {
		return nil
	}
	if _, err := g.SetCurrentView(fileTreeView); err != nil {
		return err
	}
	return nil
}

func saveEditFile(g *gocui.Gui) error {
	// ghostDisNodeList := KVTree.DisNodeList
	source := KVTree.DisNodeList[treeIndex].Source
	if source == nil {
		return errors.New("nil value")
	}
	var value interface{}
	var key = ""
	for k, v := range *source {
		value = v
		key = k
		break
	}
	view, err := g.View(fileEditView)
	if err != nil {
		return err
	}
	buff := strings.ReplaceAll(view.Buffer(), " ", "")
	// 去除回车符 (\r)
	buff = strings.ReplaceAll(buff, "\r", "")
	// 去除换行符 (\n)
	buff = strings.ReplaceAll(buff, "\n", "")
	switch value.(type) {
	case string:
		(*KVTree.DisNodeList[treeIndex].Source)[key] = buff
	case float64:
		floatValue, err := cast.ToFloat64E(buff)
		if err != nil {
			pageError(g, "float value error: "+buff)
			return err
		}
		(*KVTree.DisNodeList[treeIndex].Source)[key] = floatValue
		peonDebug(fmt.Sprint(&(KVTree.DisNodeList[treeIndex].Source)))
	case map[string]interface{}:
		peonDebug(" map[string]interface{}")
	case []interface{}:
		peonDebug(" []interface{}")
	default:
		peonDebug("default" + fmt.Sprintf("%T", value))

	}
	updatefileEditView(g)
	// jsonutils.Write(KVTree.FileName, KVTree.Source)
	return nil
}
func disfileTreeView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(fileTreeView, 0, 3, maxX/4, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "file tree"
		v.Wrap = true
		treeIndexMax = 0
		treeIndex = 0
		printKVNode(v, KVTree.NodeList, "", true)
		if _, err := g.SetCurrentView(fileTreeView); err != nil {
			return err
		}
	}
	if v, err := g.SetView(fileEditView, maxX/4, 3, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Editable = true
		v.Title = "file edit"
		v.Wrap = true
		updatefileEditView(g)
	}

	return nil
}

func pageFileTree(g *gocui.Gui, fileName string) error {
	//load 文件
	var fileData map[string]interface{}
	err := jsonutils.Read(fileName, &fileData)
	if err != nil {
		peonDebug("Error loading file: " + fileName + " " + err.Error())
		return err
	}
	KVTree.Load(fileName, fileData)
	KVTree.DisNodeList = KVTree.DisNodeList[:0]
	peonDebug("load file success " + fileName)

	updatePreviousView()
	err = disfileTreeView(g)
	if err != nil {
		return err
	}

	return nil
}
