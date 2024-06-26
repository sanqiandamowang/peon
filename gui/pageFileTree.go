package gui

import (
	"fmt"
	"peon/utils/json"
	"peon/utils/kvtree"

	"github.com/CaoYnag/gocui"
)

const (
	TreeSignDash     = "─"
	TreeSignVertical = "│"
	TreeSignUpMiddle = "├"
	TreeSignUpEnding = "└"
)

var KVTree kvtree.KV_Tree
var treeIndex int=0
var treeIndexMax int=0
func fileTreecursorDown(g *gocui.Gui, v *gocui.View) error {
	err := cursorDown(g, v)
	if err !=nil {
		return err
	}
	if treeIndex<treeIndexMax-1 {
		treeIndex+=1
	}
	return nil
}

func fileTreecursorUp(g *gocui.Gui, v *gocui.View) error {
	
	err := cursorUp(g, v)
	
	if err !=nil {
		return err
	}
	if treeIndex>0 {
		treeIndex-=1
	}
	return nil
}

func fileTreeUpdate(_ *gocui.Gui, v *gocui.View) error {
	var err error = nil
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	KVTree.DisNodeList[treeIndex].IsExpand = !KVTree.DisNodeList[treeIndex].IsExpand
	v.Clear()
	KVTree.DisNodeList = nil
	treeIndexMax =0
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
	if node.Child != nil{
		key +="->"
	}
	buf := fmt.Sprintf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)
	fmt.Fprint(v, buf)
	treeIndexMax +=1
	KVTree.DisNodeList = append(KVTree.DisNodeList, node)
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
		treeIndexMax =0
		printKVNode(v, KVTree.NodeList, "", true)
		if _, err := g.SetCurrentView(fileTreeView); err != nil {
			return err
		}
	}
	return nil
}
func pageFileTree(g *gocui.Gui, fileName string) error {
	//load 文件
	var fileData any
	err := json.Read(fileName, &fileData)
	if err != nil {
		peonDebug("Error loading file: " + fileName + " " + err.Error())
		return err
	}
	KVTree.Load(fileName, fileData)
	peonDebug("load file success" + fileName)
	
	updatePreviousView()
	err = disfileTreeView(g)
	treeIndex =0
	
	if err != nil {
		return err
	}

	return nil
}
