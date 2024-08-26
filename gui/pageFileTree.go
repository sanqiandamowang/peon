package gui

import (
	"errors"
	"fmt"
	"peon/utils/jsonutils"
	"peon/utils/kvtree"


	"github.com/CaoYnag/gocui"
	"github.com/bytedance/sonic"
)

const (
	TreeSignDash     = "─"
	TreeSignVertical = "│"
	TreeSignUpMiddle = "├"
	TreeSignUpEnding = "└"
)

var KV_tree kvtree.KV_Tree_V2
var treeIndex int = 0
var treeIndexMax int = 0
var treeIndexIsExpandBuf *bool
// var treeExPandList []bool
var disNodeList []*kvtree.KV_Node_V2
func updateTtreeIndexIsExpandBuf(){
	if treeIndexIsExpandBuf==nil{
		treeIndexIsExpandBuf = new(bool)
	}
	*treeIndexIsExpandBuf = disNodeList[treeIndex].IsExpand
	disNodeList = nil
	treeIndexMax = 0
}
func fileTreeSave(g *gocui.Gui, v *gocui.View) error {
	err := KV_tree.Save()
	if err != nil {
		peonDebug("保存失败")
		return err
	}
	peonDebug("保存成功")
	return nil

}
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
	disNodeList[treeIndex].IsExpand = !disNodeList[treeIndex].IsExpand
	v.Clear()
	updateTtreeIndexIsExpandBuf()
	
	buildDisNodeList(KV_tree.NodeList)
	printKVTree(v, KV_tree.NodeList, "", true)
	
	if err := v.SetCursor(cx, cy); err != nil && treeIndex > 0 {
		if err := v.SetOrigin(ox, oy); err != nil {
			return err
		}
	}
	return err
}
func buildDisNodeList(node *kvtree.KV_Node_V2) {

	disNodeList = append(disNodeList, node)
	if treeIndex == len(disNodeList)-1 &&  treeIndexIsExpandBuf!=nil {
		disNodeList[treeIndex].IsExpand = *treeIndexIsExpandBuf
	}
	if node.IsExpand {
		if node.Child != nil {
			buildDisNodeList(node.Child)
		}
	}
	if node.Next != nil {
		buildDisNodeList(node.Next)
	}
}
func printKVTree(v *gocui.View, node *kvtree.KV_Node_V2, indent string, isLast bool) {
	if node == nil {
		return
	}

	prefix := TreeSignUpMiddle
	if isLast {
		prefix = TreeSignUpEnding
	}

	key := node.Key
	if node.Child != nil {
		key += "->"
	}
	buf := fmt.Sprintf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)
	fmt.Fprint(v, buf)
	treeIndexMax += 1
	// disNodeList = append(disNodeList, node)
	newIndent := indent
	if !isLast {
		newIndent += TreeSignVertical + " "
	} else {
		newIndent += " "
	}
	if node.IsExpand {
		if node.Child != nil {
			printKVTree(v, node.Child, newIndent, node.Child.Next == nil)
		}
	}

	if node.Next != nil {
		printKVTree(v, node.Next, indent, node.Next.Next == nil)
	}
}

func updatefileEditView(g *gocui.Gui) error {
	var node *kvtree.KV_Node_V2 = disNodeList[treeIndex]
	var disValue interface{}
	v, err := g.View(fileEditView)
	if err != nil {
		return nil
	}
	v.Clear()
	if node.Value != nil {
		switch value := node.Value.(type) {
		case map[string]interface{}:
			if node.Key == "root" {
				disValue = node.Value
			} else {
				if disNodeList[treeIndex].Child != nil {
					disValue = node.Value
				} else {
					disValue = node.Value.(map[string]interface{})[node.Key]
				}
			}
			buf, err := sonic.ConfigDefault.MarshalIndent(disValue, "", "  ")
			if err != nil {
				return err
			}
			_buf, err := jsonutils.SortJSON(string(buf))
			if err != nil {
				return err
			}
			fmt.Fprint(v, string(_buf))
		case []interface{}:

			if node.Child != nil {
				disValue = node.Value
			} else {
				disValue = node.Value.([]interface{})[node.No]
			}

			buf, err := sonic.ConfigDefault.MarshalIndent(disValue, "", "  ")
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
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
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
	node := disNodeList[treeIndex]
	if node == nil {
		return errors.New("nil value")
	}
	view, err := g.View(fileEditView)
	if err != nil {
		return err
	}
	buf := view.Buffer()
	err,inTreeChange := KV_tree.UpdateNode(node , buf)
	if err != nil{
		pageError(g , err.Error())
		return err
	}
	
	if inTreeChange {
		v, _ := g.View(fileTreeView)
		v.Clear()
		updateTtreeIndexIsExpandBuf()
		buildDisNodeList(KV_tree.NodeList)
		printKVTree(v, KV_tree.NodeList, "", true)
	}
	updatefileEditView(g)
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
		buildDisNodeList(KV_tree.NodeList)
		printKVTree(v, KV_tree.NodeList, "", true)
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
	var fileData map[string]interface{}
	err := jsonutils.Read(fileName, &fileData)
	if err != nil {
		peonDebug("Error loading file: " + fileName + " " + err.Error())
		return err
	}
	KV_tree.Load(fileName, fileData)
	disNodeList = disNodeList[:0]
	peonDebug("load file success " + fileName)
	treeIndexIsExpandBuf =nil
	updatePreviousView()
	err = disfileTreeView(g)
	if err != nil {
		return err
	}
	return nil
}
