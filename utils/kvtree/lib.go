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

// 将json数据转换为KV树 非递归
func (tree *KV_Tree) jsonToKVNode(data interface{}, key string, parent *KV_Node) *KV_Node {
	node := &KV_Node{
		Key:      key,
		Value:    data,
		IsExpand: key == "root",
		Parent:   parent,
	}

	queue := []struct {
		data   interface{}
		node   *KV_Node
		parent *KV_Node
	}{
		{data, node, parent},
	}

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		currentNode := current.node
		currentData := current.data

		switch value := currentData.(type) {
		case map[string]interface{}:
			var lastChild *KV_Node

			// 对键进行排序以确保顺序一致
			keys := make([]string, 0, len(value))
			for k := range value {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				childNode := &KV_Node{
					Key:      k,
					Value:    value[k],
					IsExpand: false,
					Parent:   currentNode,
				}

				if currentNode.Child == nil {
					currentNode.Child = childNode
				} else {
					lastChild.Next = childNode
				}
				lastChild = childNode

				queue = append(queue, struct {
					data   interface{}
					node   *KV_Node
					parent *KV_Node
				}{value[k], childNode, currentNode})
			}

		case []interface{}:
			var lastChild *KV_Node
			for i, v := range value {
				childNode := &KV_Node{
					Key:      fmt.Sprintf("[%d]", i),
					Value:    v,
					IsExpand: false,
					Parent:   currentNode,
				}

				if currentNode.Child == nil {
					currentNode.Child = childNode
				} else {
					lastChild.Next = childNode
				}
				lastChild = childNode

				queue = append(queue, struct {
					data   interface{}
					node   *KV_Node
					parent *KV_Node
				}{v, childNode, currentNode})
			}
		}
	}

	return node
}

// 递归
//func (tree *KV_Tree) jsonToKVNode(data interface{}, key string, parent *KV_Node) *KV_Node {
//	node := &KV_Node{
//		Key:      key,
//		Value:    data,
//		IsExpand: false,
//		Parent:   parent,
//	}
//
//	if key == "root" {
//		node.IsExpand = true
//	}
//
//	switch value := data.(type) {
//	case map[string]interface{}:
//		var lastChild *KV_Node
//
//		// Sort keys to ensure consistent order
//		keys := make([]string, 0, len(value))
//		for k := range value {
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//
//		for _, k := range keys {
//			child := tree.jsonToKVNode(value[k], k, node)
//			if node.Child == nil {
//				node.Child = child
//			} else {
//				lastChild.Next = child
//			}
//			lastChild = child
//		}
//	case []interface{}:
//		var lastChild *KV_Node
//		for i, v := range value {
//			child := tree.jsonToKVNode(v, fmt.Sprintf("[%d]", i), node)
//			if node.Child == nil {
//				node.Child = child
//			} else {
//				lastChild.Next = child
//			}
//			lastChild = child
//		}
//	}
//
//	return node
//}

// 更新父节点 非递归
func (tree *KV_Tree) UpdateParentNodes(node *KV_Node) {
	for node != nil && node.Parent != nil {
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
		node = parent
	}
}

// 更新子节点 非递归
func (tree *KV_Tree) UpdateChildNodes(node *KV_Node) *KV_Node {
	if node == nil {
		return nil
	}

	queue := []*KV_Node{node}

	for len(queue) > 0 {
		currentNode := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// 清空当前节点的子节点
		currentNode.Child = nil

		// 根据节点的值类型来构建新的子节点链表
		switch value := currentNode.Value.(type) {
		case map[string]interface{}:
			var lastChild *KV_Node
			for key, v := range value {
				childNode := &KV_Node{
					Key:   key,
					Value: v,
				}
				if currentNode.Child == nil {
					currentNode.Child = childNode
				} else {
					lastChild.Next = childNode
				}
				lastChild = childNode

				queue = append(queue, childNode)
			}
		case []interface{}:
			var lastChild *KV_Node
			for i, v := range value {
				childNode := &KV_Node{
					Key:   fmt.Sprintf("[%d]", i),
					Value: v,
				}
				if currentNode.Child == nil {
					currentNode.Child = childNode
				} else {
					lastChild.Next = childNode
				}
				lastChild = childNode

				queue = append(queue, childNode)
			}
		}
	}

	return node
}

// // 更新父节点 递归
//
//	func (tree *KV_Tree) UpdateParentNodes(node *KV_Node) {
//		if node == nil || node.Parent == nil {
//			return
//		}
//		parent := node.Parent
//		switch parent.Value.(type) {
//		case map[string]interface{}:
//			newMap := make(map[string]interface{})
//			for child := parent.Child; child != nil; child = child.Next {
//				newMap[child.Key] = child.Value
//			}
//			parent.Value = newMap
//		case []interface{}:
//			newSlice := make([]interface{}, 0)
//			for child := parent.Child; child != nil; child = child.Next {
//				newSlice = append(newSlice, child.Value)
//			}
//			parent.Value = newSlice
//		}
//		tree.UpdateParentNodes(parent)
//	}
//
// // 更新子节点 递归
//
//	func (tree *KV_Tree) UpdateChildNodes(node *KV_Node) *KV_Node {
//		if node == nil {
//			return nil
//		}
//
//		// 清空原有的子节点
//		node.Child = nil
//
//			// 根据节点的值类型来构建新的子节点链表
//			switch value := node.Value.(type) {
//			case map[string]interface{}:
//				var lastChild *KV_Node
//				for key, v := range value {
//					childNode := &KV_Node{
//						Key:   key,
//						Value: v,
//					}
//					if node.Child == nil {
//						node.Child = childNode
//					} else {
//						lastChild.Next = childNode
//					}
//					lastChild = childNode
//
//					// 递归更新子节点的子节点
//					tree.UpdateChildNodes(childNode)
//				}
//			case []interface{}:
//				var lastChild *KV_Node
//				for i, v := range value {
//					childNode := &KV_Node{
//						Key:   fmt.Sprintf("[%d]", i),
//						Value: v,
//					}
//					if node.Child == nil {
//						node.Child = childNode
//					} else {
//						lastChild.Next = childNode
//					}
//					lastChild = childNode
//
//					// 递归更新子节点的子节点
//					tree.UpdateChildNodes(childNode)
//				}
//			}
//
//
//		return node
//	}
func (tree *KV_Tree) Load(fileName string, source map[string]interface{}) error {
	tree.FileName = fileName
	tree.Source = &source
	tree.NodeList = tree.jsonToKVNode(source, "root", nil)
	tree.DisNodeList = make([]*KV_Node, 0)
	return nil
}

// 打印KVNode 非递归
func (tree *KV_Tree) printKVNode(root *KV_Node, indent string, isLast bool) {
	type nodeState struct {
		node   *KV_Node
		indent string
		isLast bool
	}

	// 初始化栈，将根节点入队列
	queue := []nodeState{{root, indent, isLast}}

	for len(queue) > 0 {
		// 弹出队列
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if current.node == nil {
			continue
		}

		// 构建前缀
		prefix := TreeSignUpMiddle
		if current.isLast {
			prefix = TreeSignUpEnding
		}

		// 输出当前节点信息
		key := current.node.Key
		fmt.Printf("%s%s%s %s\n", current.indent, prefix, TreeSignDash, key)
		tree.DisNodeList = append(tree.DisNodeList, current.node)

		// 计算新的缩进
		newIndent := current.indent
		if !current.isLast {
			newIndent += TreeSignVertical + " "
		} else {
			newIndent += " "
		}

		// 先处理兄弟节点（Next），再处理子节点（Child）
		if current.node.Next != nil {
			queue = append(queue, nodeState{current.node.Next, current.indent, current.node.Next.Next == nil})
		}

		if current.node.IsExpand && current.node.Child != nil {
			queue = append(queue, nodeState{current.node.Child, newIndent, current.node.Child.Next == nil})
		}
	}
}

//	func (tree *KV_Tree) printKVNode(node *KV_Node, indent string, isLast bool) {
//		if node == nil {
//			return
//		}
//
//		prefix := TreeSignUpMiddle
//		if isLast {
//			prefix = TreeSignUpEnding
//		}
//
//		key := node.Key
//		fmt.Printf("%s%s%s %s\n", indent, prefix, TreeSignDash, key)
//		tree.DisNodeList = append(tree.DisNodeList, node)
//
//		newIndent := indent
//		if !isLast {
//			newIndent += TreeSignVertical + " "
//		} else {
//			newIndent += " "
//		}
//
//		if node.IsExpand {
//			if node.Child != nil {
//				tree.printKVNode(node.Child, newIndent, node.Child.Next == nil)
//			}
//		}
//
//		if node.Next != nil {
//			tree.printKVNode(node.Next, indent, node.Next.Next == nil)
//		}
//	}
func (tree *KV_Tree) Save() error {
	data := tree.nodeToMap()

	err := jsonutils.Write(tree.FileName, data)
	return err
}

func (tree *KV_Tree) nodeToMap() map[string]interface{} {
	// 如果根节点为空，返回 nil
	if tree.NodeList == nil {
		return nil
	}

	// 初始化根节点的映射表
	rootMap := make(map[string]interface{})

	// 使用队列来模拟树的广度优先遍历，初始时将根节点放入队列中
	queue := []*KV_Node{tree.NodeList}

	// nodeMap 用来存储每个节点对应的 map，便于构建嵌套结构
	nodeMap := make(map[*KV_Node]map[string]interface{})

	// visited 用来记录已经访问过的节点，防止重复处理
	visited := make(map[*KV_Node]bool)

	// 将根节点映射到 rootMap
	nodeMap[tree.NodeList] = rootMap

	// 开始遍历节点
	for len(queue) > 0 {
		// 从队列的前端取出节点进行处理
		currentNode := queue[0]
		queue = queue[1:]

		// 如果当前节点已经访问过，则跳过
		if visited[currentNode] {
			continue
		}
		// 标记当前节点为已访问
		visited[currentNode] = true

		// 获取或创建当前节点对应的映射表
		currentMap, exists := nodeMap[currentNode]
		if !exists {
			currentMap = make(map[string]interface{})
			nodeMap[currentNode] = currentMap
		}

		// 如果不是根节点，则将节点的 Key 和 Value 添加到映射表中
		if currentNode != tree.NodeList {
			currentMap[currentNode.Key] = currentNode.Value
		}

		// 如果当前节点有子节点，处理子节点
		if currentNode.Child != nil {
			childNode := currentNode.Child
			for childNode != nil {
				if !visited[childNode] {
					queue = append(queue, childNode) // 将子节点添加到队列末尾

					// 如果子节点的值不是 map，则直接添加到当前映射表中
					if _, ok := childNode.Value.(map[string]interface{}); !ok {
						currentMap[childNode.Key] = childNode.Value
					} else {
						// 否则，为子节点创建新的 map，并存储在当前映射表中
						childMap := make(map[string]interface{})
						currentMap[childNode.Key] = childMap
						nodeMap[childNode] = childMap
					}
				}
				// 继续处理下一个兄弟节点
				childNode = childNode.Next
			}
		}
	}

	// 返回最终构建的根节点映射表
	return rootMap
}
