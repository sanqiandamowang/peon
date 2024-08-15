package gui

import (
	"errors"
	"fmt"
	"peon/utils/jsonutils"
	"peon/utils/kvtree"

	// "regexp"
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

var KVTree kvtree.KV_Tree_V2
var treeIndex int = 0
var treeIndexMax int = 0
var treeExPandList []bool
func updateTreeExPandList(){
	treeExPandList = nil
	for i:=0 ;i<treeIndex;i++{
		treeExPandList = append(treeExPandList, KVTree.DisNodeList[i].IsExpand)
	}
}
func fileTreeSave(g *gocui.Gui, v *gocui.View) error {
	err := KVTree.Save()
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
	KVTree.DisNodeList[treeIndex].IsExpand = !KVTree.DisNodeList[treeIndex].IsExpand
	v.Clear()
	KVTree.DisNodeList = nil
	treeIndexMax = 0
	buildDisNodeList(KVTree.NodeList)
	printKVTree(v, KVTree.NodeList, "", true)
	updateTreeExPandList()
	if err := v.SetCursor(cx, cy); err != nil && treeIndex > 0 {
		if err := v.SetOrigin(ox, oy); err != nil {
			return err
		}
	}
	return err
}

// 非递归
// func printKVNode(v *gocui.View, root *kvtree.KV_Node, indent string, isLast bool) {
// 	type nodeState struct {
// 		node   *kvtree.KV_Node
// 		indent string
// 		isLast bool
// 	}

// 	queue := []nodeState{{root, indent, isLast}}

// 	for len(queue) > 0 {
// 		// 弹出栈顶元素
// 		current := queue[len(queue)-1]
// 		queue = queue[:len(queue)-1]

// 		if current.node == nil {
// 			continue
// 		}

// 		// 构造前缀
// 		prefix := TreeSignUpMiddle
// 		if current.isLast {
// 			prefix = TreeSignUpEnding
// 		}

// 		// 处理当前节点
// 		key := current.node.Key
// 		if current.node.Child != nil {
// 			key += "->"
// 		}
// 		buf := fmt.Sprintf("%s%s%s %s\n", current.indent, prefix, TreeSignDash, key)
// 		fmt.Fprint(v, buf)
// 		treeIndexMax += 1
// 		KVTree.DisNodeList = append(KVTree.DisNodeList, current.node)

// 		// 更新缩进
// 		newIndent := current.indent
// 		if !current.isLast {
// 			newIndent += TreeSignVertical + " "
// 		} else {
// 			newIndent += " "
// 		}

// 		// 如果有兄弟节点，先把兄弟节点放入队列
// 		if current.node.Next != nil {
// 			queue = append(queue, nodeState{current.node.Next, current.indent, current.node.Next.Next == nil})
// 		}

//			// 如果当前节点是展开的，处理子节点，否则跳过
//			if current.node.IsExpand && current.node.Child != nil {
//				queue = append(queue, nodeState{current.node.Child, newIndent, current.node.Child.Next == nil})
//			}
//		}
//	}

func buildDisNodeList(node *kvtree.KV_Node_V2) {

	KVTree.DisNodeList = append(KVTree.DisNodeList, node)
	for i := 0; i < treeIndex; i++ {
		if i >= len(treeExPandList) {
			break;
		}
		if(i>= len(KVTree.DisNodeList)){
			break;
		}
		KVTree.DisNodeList[i].IsExpand = treeExPandList[i]
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
	// KVTree.DisNodeList = append(KVTree.DisNodeList, node)
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
	var node *kvtree.KV_Node_V2 = KVTree.DisNodeList[treeIndex]
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
				if KVTree.DisNodeList[treeIndex].Child != nil {
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
	node := KVTree.DisNodeList[treeIndex]
	if node == nil {
		return errors.New("nil value")
	}
	var sourceValue interface{}
	switch node.Value.(type) {
	case map[string]interface{}:
		if node.Key == "root" {
			sourceValue = node.Value
		} else {
			if node.Child != nil {
				sourceValue = node.Value
			} else {
				sourceValue = node.Value.(map[string]interface{})[node.Key]
			}
		}
	case []interface{}:
		if node.Child != nil {
			sourceValue = node.Value
		} else {
			sourceValue = node.Value.([]interface{})[node.No]
		}
	default:
		sourceValue = nil
	}

	view, err := g.View(fileEditView)
	if err != nil {
		return err
	}
	buff := strings.ReplaceAll(view.Buffer(), " ", "")
	buff = strings.ReplaceAll(buff, "\r", "")
	buff = strings.ReplaceAll(buff, "\n", "")
	var treeChangeFlag = false
	parentValueType := KVTree.GetPraentValueType(node)
	switch value := sourceValue.(type) {
	case string:
		buff = strings.ReplaceAll(buff, "\"", "")
		if parentValueType == kvtree.TYPE_MAP {
			node.Value.(map[string]interface{})[node.Key] = buff
		} else if parentValueType == kvtree.TYPE_ARRAY {
			node.Value.([]interface{})[node.No] = buff
		}
	case float64:
		floatValue, err := cast.ToFloat64E(buff)
		if err != nil {
			pageError(g, "float value error: "+buff)
			return err
		}
		if parentValueType == kvtree.TYPE_MAP {
			node.Value.(map[string]interface{})[node.Key] = floatValue
		} else if parentValueType == kvtree.TYPE_ARRAY {
			node.Value.([]interface{})[node.No] = floatValue
		}
	case map[string]interface{}:
		var newData map[string]interface{}
		err := sonic.Unmarshal([]byte(buff), &newData)
		if err != nil {
			pageError(g, "error decoding JSON: "+err.Error())
			return err
		}
		if parentValueType == kvtree.TYPE_MAP {
			node.Parent.Value.(map[string]interface{})[node.Key] = newData
		} else if parentValueType == kvtree.TYPE_ARRAY {
			node.Value.([]interface{})[node.No] = newData
		} else if parentValueType == kvtree.TYPE_ROOT {

			KVTree.NodeList.Value = newData
		}
		treeChangeFlag = true
	case []interface{}:
		var newData []interface{}
		err := sonic.Unmarshal([]byte(buff), &newData)
		if err != nil {
			pageError(g, "error decoding JSON: "+err.Error())
			return err
		}
		if parentValueType == kvtree.TYPE_MAP {
			node.Parent.Value.(map[string]interface{})[node.Key]= newData
		} else if parentValueType == kvtree.TYPE_ARRAY {
			node.Parent.Value.([]interface{})[node.No] = newData
		}
		treeChangeFlag = true
	case nil:
		var newData interface{}
		var newDataArray []interface{}
		err := sonic.Unmarshal([]byte(buff), &newData)
		if err != nil {
			err := sonic.Unmarshal([]byte(buff), &newDataArray)
			{
				if err != nil {
					err := sonic.Unmarshal([]byte(buff), &newDataArray)
					pageError(g, "error decoding JSON: "+err.Error())
					return err
				} else {
					node.Value.(map[string]interface{})[node.Key] = newDataArray
					return nil
				}
			}
		}
		node.Value.(map[string]interface{})[node.Key] = newData
		treeChangeFlag = true
	default:
		pageError(g, "unsupported type: "+fmt.Sprintf("%T", value))
		return errors.New("unsupported type")
	}
	if treeChangeFlag {
		//更新节点
		KVTree.Source = KVTree.NodeList.Value
		KVTree.NodeList = KVTree.SourceToKVNode(KVTree.Source, "root", nil)
		if node.Parent != nil {
			KVTree.SourceToKVNode(node.Parent.Value, node.Parent.Key, node.Parent.Parent)
		}
		v, _ := g.View(fileTreeView)
		v.Clear()
		KVTree.DisNodeList = nil
		treeIndexMax = 0
		buildDisNodeList(KVTree.NodeList)
		printKVTree(v, KVTree.NodeList, "", true)
		updateTreeExPandList()
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
		buildDisNodeList(KVTree.NodeList)
		printKVTree(v, KVTree.NodeList, "", true)
		updateTreeExPandList()
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
	KVTree.Load(fileName, fileData)
	KVTree.DisNodeList = KVTree.DisNodeList[:0]
	peonDebug("load file success " + fileName)
	treeExPandList =nil
	updatePreviousView()
	err = disfileTreeView(g)
	if err != nil {
		return err
	}
	return nil
}
