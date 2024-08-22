package kvtree

import (
	"fmt"
	"peon/utils/jsonutils"
	"sort"
)

type KV_Node_V2 struct {
	Key      string
	Value    interface{}
	IsExpand bool
	Parent   *KV_Node_V2
	Next     *KV_Node_V2
	Child    *KV_Node_V2
	No       int
}
type KV_Tree_V2 struct {
	FileName    string
	Source      interface{}
	NodeList    *KV_Node_V2
	// DisNodeList []*KV_Node_V2
}
// 枚举
const (
	TYPE_ROOT = iota
	TYPE_MAP
	TYPE_ARRAY
	TYPE_ERR
)

// 0 root 1 map 2 []interface{}
func (tree *KV_Tree_V2) GetPraentValueType(node *KV_Node_V2) int {

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

func (tree *KV_Tree_V2) UpdateNodeWithChild(node *KV_Node_V2) *KV_Node_V2 {
	parent := node.Parent
	if parent == nil {
		//root
		tree.Source = tree.NodeList.Value
		tree.NodeList = tree.SourceToKVNode(tree.Source, "root", nil)
		return nil
	}
	var newNode *KV_Node_V2
	// var lastNode *KV_Node_V2
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
func (tree *KV_Tree_V2) SourceToKVNode(data interface{}, key string, parent *KV_Node_V2) *KV_Node_V2 {
	var node *KV_Node_V2
	switch dataType := data.(type) {
	case map[string]interface{}:
		node = &KV_Node_V2{
			Key:      key,
			Value:    data,
			IsExpand: key == "root",
			Parent: parent,
		}
		var lastNode *KV_Node_V2
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
		node = &KV_Node_V2{
			Key:      key,
			Value:    data,
			IsExpand: key == "root",
			Parent: parent,
		}
		var lastNode *KV_Node_V2
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
			node = &KV_Node_V2{
				Key:      key,
				Value:    parent.Value,
				IsExpand: key == "root",
				Parent: parent,
			}
		}
	}
	return node
}
func (tree *KV_Tree_V2) printKVNode(node *KV_Node_V2, indent string, isLast bool) {
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
	// var DisNodeList []*KV_Node_V2
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
func (tree *KV_Tree_V2) Load(fileName string, source map[string]interface{}) error {
	tree.FileName = fileName
	tree.Source = &source
	tree.NodeList = tree.SourceToKVNode(source, "root", nil)
	// tree.DisNodeList = make([]*KV_Node_V2, 0)
	return nil
}
func (tree *KV_Tree_V2) Save() error {
	err := jsonutils.Write(tree.FileName, tree.NodeList.Value)
	return err
}
