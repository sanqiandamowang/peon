package kvtree

import "fmt"

const (
	TreeSignDash     = "─"
	TreeSignVertical = "│"
	TreeSignUpMiddle = "├"
	TreeSignUpEnding = "└"
)

type KV_Node struct {
	Index    int
	IsExpand bool
	Source   *map[string]any
	Next     *KV_Node
	Child    *KV_Node
}

type KV_Tree struct {
	FileName    string
	Source      any
	NodeList    *KV_Node
	DisNodeList []*KV_Node
}

func (tree *KV_Tree) jsonToKVNode(data interface{}, key string, index *int) *KV_Node {
	node := &KV_Node{
		Index:    *index,
		IsExpand: false,
	}
	if node.Index == 0 {
		node.IsExpand = true
	}
	(*index)++
	switch value := data.(type) {
	case map[string]interface{}:
		source := map[string]any{key: value}
		node.Source = &source
		var lastChild *KV_Node
		for k, v := range value {
			child := tree.jsonToKVNode(v, k, index)
			if node.Child == nil {
				node.Child = child
			} else {
				lastChild.Next = child
			}
			lastChild = child
		}
	case []interface{}:
		source := map[string]any{key: value}
		node.Source = &source
		var lastChild *KV_Node
		for i, v := range value {
			child := tree.jsonToKVNode(v, fmt.Sprintf("[%d]", i), index)
			if node.Child == nil {
				node.Child = child
			} else {
				lastChild.Next = child
			}
			lastChild = child

		}
	default:
		source := map[string]any{key: value}
		node.Source = &source
	}

	return node
}

// Function to print KV_Node structure with border using tree symbols
func printKVNode(node *KV_Node, indent string, isLast bool) {
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

	fmt.Printf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)

	// Prepare the new indent for the child nodes
	newIndent := indent
	if !isLast {
		newIndent += TreeSignVertical + " "
	} else {
		newIndent += " "
	}

	// Print the child nodes
	if node.Child != nil {
		printKVNode(node.Child, newIndent, node.Child.Next == nil)
	}

	// Print the sibling nodes
	if node.Next != nil {
		printKVNode(node.Next, indent, node.Next.Next == nil)
	}
}

func (tree *KV_Tree) Load(fileName string, source any) error {
	tree.FileName = fileName
	tree.Source = source

	index := 0
	tree.NodeList = tree.jsonToKVNode(source, "root", &index)

	return nil
}
func (tree *KV_Tree) Update() {
	index := 0
	tree.NodeList = tree.jsonToKVNode(tree.Source, "root", &index)
}
