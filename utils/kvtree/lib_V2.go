package kvtree

import (
	"fmt"
	"peon/utils/jsonutils"
	"regexp"
	"sort"
	"strconv"
)

type KV_Node_V2 struct {
	Key      string
	Value    interface{}
	IsExpand bool
	Parent   *KV_Node_V2
	Next     *KV_Node_V2
	Child    *KV_Node_V2
}
type KV_Tree_V2 struct {
	FileName    string
	Source      interface{}
	NodeList    *KV_Node_V2
	DisNodeList []*KV_Node_V2
}

func (tree *KV_Tree_V2) sourceToKVNode(data interface{}, key string, parent *KV_Node_V2) *KV_Node_V2 {

	var node *KV_Node_V2
	switch dataType := data.(type) {
	case map[string]interface{}:
		node = &KV_Node_V2{
			Key:   key,
			Value: data,
			IsExpand: key == "root",
			// IsExpand: true,
			Parent:   parent,
		}
		var lastNode *KV_Node_V2
		keys := make([]string, 0, len(node.Value.(map[string]interface{})))
		for k := range dataType {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			newNode := tree.sourceToKVNode(node.Value.(map[string]interface{})[k], k, node)
			if node.Child == nil {
				node.Child = newNode
			} else {
				lastNode.Next = newNode
			}
			lastNode = newNode
		}
	case []interface{}:
		node = &KV_Node_V2{
			Key:   key,
			Value: data,
			IsExpand: key == "root",
			// IsExpand: true,
			Parent:   parent,
		}
		var lastNode *KV_Node_V2
		for i := range dataType {
			newNode := tree.sourceToKVNode(node.Value.([]interface{})[i], fmt.Sprintf("[%d]", i), node)
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
				// IsExpand:  true,
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
			re := regexp.MustCompile(`\[(\d+)\]`)
			matches := re.FindStringSubmatch(node.Key)
			no, _ := strconv.Atoi(matches[1])
			fmt.Printf("%s%s%s %s %v\n", indent, prefix, TreeSignDash, node.Key, node.Value.([]interface{})[no])
		}

	}
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
func (tree *KV_Tree_V2) Load(fileName string, source map[string]interface{}) error {
	tree.FileName = fileName
	tree.Source = &source
	tree.NodeList = tree.sourceToKVNode(source, "root", nil)
	tree.DisNodeList = make([]*KV_Node_V2, 0)
	return nil
}
func (tree *KV_Tree_V2) Save() error {
	err := jsonutils.Write(tree.FileName, tree.NodeList.Value)
	return err
}
