package kvtree

import (
	"fmt"
	"peon/utils/jsonutils"
	"sort"
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
}

type KV_Tree struct {
	FileName    string
	Source      *map[string]interface{}
	NodeList    *KV_Node
	DisNodeList []*KV_Node
}

func (tree *KV_Tree) jsonToKVNode(data interface{}, key string, parent *KV_Node) *KV_Node {
	node := &KV_Node{
		Key:      key,
		Value:    data,
		IsExpand: false,
		Parent:   parent,
	}

	if key == "root" {
		node.IsExpand = true
	}

	switch value := data.(type) {
	case map[string]interface{}:
		var lastChild *KV_Node

		// Sort keys to ensure consistent order
		keys := make([]string, 0, len(value))
		for k := range value {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			child := tree.jsonToKVNode(value[k], k, node)
			if node.Child == nil {
				node.Child = child
			} else {
				lastChild.Next = child
			}
			lastChild = child
		}
	case []interface{}:
		var lastChild *KV_Node
		for i, v := range value {
			child := tree.jsonToKVNode(v, fmt.Sprintf("[%d]", i), node)
			if node.Child == nil {
				node.Child = child
			} else {
				lastChild.Next = child
			}
			lastChild = child
		}
	}

	return node
}

// 更新父节点
func (tree *KV_Tree) UpdateParentNodes(node *KV_Node) {
	if node == nil || node.Parent == nil {
		return
	}
	parent := node.Parent
	switch parent.Value.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for child := parent.Child; child != nil; child = child.Next {
			newMap[child.Key] = child.Value
		}
		parent.Value = newMap
	case []interface{}:
		newSlice := make([]interface{}, 0)
		for child := parent.Child; child != nil; child = child.Next {
			newSlice = append(newSlice, child.Value)
		}
		parent.Value = newSlice
	}
	tree.UpdateParentNodes(parent)
}

// 更新子节点
func (tree *KV_Tree) UpdateChildNodes(node *KV_Node) *KV_Node {
	if node == nil {
		return nil
	}

	// 清空原有的子节点
	node.Child = nil

	// 根据节点的值类型来构建新的子节点链表
	switch value := node.Value.(type) {
	case map[string]interface{}:
		var lastChild *KV_Node
		for key, v := range value {
			childNode := &KV_Node{
				Key:   key,
				Value: v,
			}
			if node.Child == nil {
				node.Child = childNode
			} else {
				lastChild.Next = childNode
			}
			lastChild = childNode
		}
	case []interface{}:
		var lastChild *KV_Node
		for i, v := range value {
			childNode := &KV_Node{
				Key:   fmt.Sprintf("[%d]", i),
				Value: v,
			}
			if node.Child == nil {
				node.Child = childNode
			} else {
				lastChild.Next = childNode
			}
			lastChild = childNode
		}
	}

	return node
}
func (tree *KV_Tree) Load(fileName string, source map[string]interface{}) error {
	tree.FileName = fileName
	tree.Source = &source
	tree.NodeList = tree.jsonToKVNode(source, "root", nil)
	tree.DisNodeList = make([]*KV_Node, 0)
	return nil
}

func (tree *KV_Tree) printKVNode(node *KV_Node, indent string, isLast bool) {
	if node == nil {
		return
	}

	prefix := TreeSignUpMiddle
	if isLast {
		prefix = TreeSignUpEnding
	}

	key := node.Key
	fmt.Printf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)
	tree.DisNodeList = append(tree.DisNodeList, node)

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
func (tree *KV_Tree)Save() error {
	err:=jsonutils.Write(tree.FileName, tree.NodeList.Value)
	return err
}
