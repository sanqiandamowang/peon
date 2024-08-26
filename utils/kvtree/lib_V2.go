package kvtree

import (
	"errors"
	"fmt"
	"peon/utils/jsonutils"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
)
const (
	TreeSignDash     = "─"
	TreeSignVertical = "│"
	TreeSignUpMiddle = "├"
	TreeSignUpEnding = "└"
)
type KV_Node struct {
	Key      string
	Value    interface{}
	IsExpand bool
	Parent   *KV_Node
	Next     *KV_Node
	Child    *KV_Node
	No       int
}
type KV_Tree struct {
	FileName string
	Source   interface{}
	NodeList *KV_Node
	// DisNodeList []*KV_Node
}

// 枚举
const (
	TYPE_ROOT = iota
	TYPE_MAP
	TYPE_ARRAY
	TYPE_ERR
)

// 0 root 1 map 2 []interface{}
func (tree *KV_Tree) getPraentValueType(node *KV_Node) int {

	if node.Parent == nil {
		if node.Key == "root" {
			return TYPE_ROOT
		}
		return TYPE_ERR
	}
	switch node.Parent.Value.(type) {

	case map[string]interface{}:
		return TYPE_MAP
	case []interface{}:
		return TYPE_ARRAY
	default:
		return TYPE_ERR
	}
}
func (tree *KV_Tree) updateNodeWithChild(node *KV_Node) *KV_Node {
	parent := node.Parent
	if parent == nil {
		//root
		tree.Source = tree.NodeList.Value
		tree.NodeList = tree.SourceToKVNode(tree.Source, "root", nil)
		return nil
	}
	var newNode *KV_Node
	// var lastNode *KV_Node
	switch parent.Value.(type) {
	case map[string]interface{}:
		newNode = tree.SourceToKVNode(parent.Value.(map[string]interface{})[node.Key], node.Key, parent)
	case []interface{}:
		newNode = tree.SourceToKVNode(parent.Value.([]interface{})[node.No], fmt.Sprintf("[%d]", node.No), parent)
	}
	lastNode := parent.Child
	if lastNode.Key == node.Key {
		parent.Child = newNode
		if lastNode.Next != nil {
			newNode.Next = lastNode.Next
		}
	} else {
		for lastNode.Next != nil {
			if lastNode.Next.Key == node.Key {
				break
			}
			lastNode = lastNode.Next
		}
		if lastNode.Next.Next != nil {
			newNode.Next = lastNode.Next.Next
		}
		lastNode.Next = newNode
	}
	return newNode
}
func (tree *KV_Tree) getNodeValue(node *KV_Node) interface{} {
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
	return sourceValue
}
func (tree *KV_Tree) UpdateNode(node *KV_Node, updateBuf string) (err error,isTrerChange bool) {
	sourceValue:= tree.getNodeValue(node)
	buf := strings.ReplaceAll(updateBuf, " ", "")
	buf = strings.ReplaceAll(buf, "\r", "")
	buf = strings.ReplaceAll(buf, "\n", "")
	var treeChangeFlag = false
	parentValueType := tree.getPraentValueType(node)
	switch nodeValueType := sourceValue.(type){
	case string:
		buf = strings.ReplaceAll(buf, "\"", "")
		if parentValueType == TYPE_MAP {
			node.Value.(map[string]interface{})[node.Key] = buf
		} else if parentValueType == TYPE_ARRAY {
			node.Value.([]interface{})[node.No] = buf
		}
	case float64:
		newFloatValue, err := cast.ToFloat64E(buf)
		if err != nil {
			// pageError(g, "float value error: "+buff)
			return err,false
		}
		if parentValueType == TYPE_MAP {
			node.Value.(map[string]interface{})[node.Key] = newFloatValue
		} else if parentValueType == TYPE_ARRAY {
			node.Value.([]interface{})[node.No] = newFloatValue
		}
	case map[string]interface{}:
		var newMapData map[string]interface{}
		err := sonic.Unmarshal([]byte(buf), &newMapData)
		if err != nil {
			// pageError(g, "error decoding JSON: "+err.Error())
			return err,false
		}
		if parentValueType == TYPE_MAP {
			node.Parent.Value.(map[string]interface{})[node.Key] = newMapData
		} else if parentValueType == TYPE_ARRAY {
			node.Parent.Value.([]interface{})[node.No] = newMapData
		} else if parentValueType == TYPE_ROOT {
			tree.NodeList.Value = newMapData 
		}
		treeChangeFlag = true
	case []interface{}:
		var newArrayData []interface{}
		err := sonic.Unmarshal([]byte(buf), &newArrayData)
		if err != nil {
			// pageError(g, "error decoding JSON: "+err.Error())
			return err,false
		}
		if parentValueType == TYPE_MAP {
			node.Parent.Value.(map[string]interface{})[node.Key]= newArrayData
		} else if parentValueType == TYPE_ARRAY {
			node.Parent.Value.([]interface{})[node.No] = newArrayData
		}
		treeChangeFlag = true
	case nil:
		var newData interface{}
		var newDataArray []interface{}
		err := sonic.Unmarshal([]byte(buf), &newData)
		if err != nil {
			err := sonic.Unmarshal([]byte(buf), &newDataArray)
			{
				if err != nil {
					err := sonic.Unmarshal([]byte(buf), &newDataArray)
					//pageError(g, "error decoding JSON: "+err.Error())
					return err,false
				} else {
					node.Value.(map[string]interface{})[node.Key] = newDataArray
				}
			}
		}else
		{
			node.Value.(map[string]interface{})[node.Key] = newData
		}
		treeChangeFlag = true
	default:
		// pageError(g, "unsupported type: "+fmt.Sprintf("%T", value))
		return errors.New("unsupported type: "+fmt.Sprintf("%T", nodeValueType)),false
	}
	if treeChangeFlag {
		tree.updateNodeWithChild(node)
		// v, _ := g.View(fileTreeView)
		// v.Clear()
		
		// updateTtreeIndexIsExpandBuf()
		// buildDisNodeList(KV_tree.NodeList)
		// printKVTree(v, KV_tree.NodeList, "", true)
	}
	return nil,treeChangeFlag
}
func (tree *KV_Tree) SourceToKVNode(data interface{}, key string, parent *KV_Node) *KV_Node {
	var node *KV_Node
	switch dataType := data.(type) {
	case map[string]interface{}:
		node = &KV_Node{
			Key:      key,
			Value:    data,
			IsExpand: key == "root",
			Parent:   parent,
		}
		var lastNode *KV_Node
		keys := make([]string, 0, len(node.Value.(map[string]interface{})))
		for k := range dataType {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			newNode := tree.SourceToKVNode(node.Value.(map[string]interface{})[k], k, node)
			if node.Child == nil {
				node.Child = newNode
			} else {
				lastNode.Next = newNode
			}
			lastNode = newNode
		}
	case []interface{}:
		node = &KV_Node{
			Key:      key,
			Value:    data,
			IsExpand: key == "root",
			Parent:   parent,
		}
		var lastNode *KV_Node
		for i := range dataType {
			newNode := tree.SourceToKVNode(node.Value.([]interface{})[i], fmt.Sprintf("[%d]", i), node)
			newNode.No = i
			if node.Child == nil {
				node.Child = newNode
			} else {
				lastNode.Next = newNode
			}
			lastNode = newNode
		}
	default:
		if parent != nil {
			node = &KV_Node{
				Key:      key,
				Value:    parent.Value,
				IsExpand: key == "root",
				Parent:   parent,
			}
		}
	}
	return node
}
func (tree *KV_Tree) printKVNode(node *KV_Node, indent string, isLast bool) {
	if node == nil {
		return
	}
	prefix := TreeSignUpMiddle
	if isLast {
		prefix = TreeSignUpEnding
	}
	switch node.Value.(type) {
	case map[string]interface{}:
		if node.Key == "root" {
			fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value)
		} else {
			if node.Child != nil {
				fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value)
			} else {
				fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value.(map[string]interface{})[node.Key])
			}
		}
	case []interface{}:
		if node.Child != nil {
			fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value)
		} else {
			fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value.([]interface{})[node.No])
		}

	}
	// var DisNodeList []*KV_Node
	// DisNodeList = append(DisNodeList, node)
	newIndent := indent
	if !isLast {
		newIndent += TreeSignVertical + " "
	} else {
		newIndent += " "
	}
	if node.IsExpand {
		if node.Child != nil {
			tree.printKVNode(node.Child, newIndent, node.Child.Next == nil)
		}
	}
	if node.Next != nil {
		tree.printKVNode(node.Next, indent, node.Next.Next == nil)
	}
}
func (tree *KV_Tree) Load(fileName string, source map[string]interface{}) error {
	tree.FileName = fileName
	tree.Source = &source
	tree.NodeList = tree.SourceToKVNode(source, "root", nil)
	// tree.DisNodeList = make([]*KV_Node, 0)
	return nil
}
func (tree *KV_Tree) Save() error {
	err := jsonutils.Write(tree.FileName, tree.NodeList.Value)
	return err
}
