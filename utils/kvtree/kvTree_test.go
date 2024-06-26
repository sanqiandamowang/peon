package kvtree

import (
	"peon/utils/json"
	"testing"
)

func TestLoad(t *testing.T) {
	var jsonData any
	err := json.Read("/home/sanqian/peon/peon/utils/json/dev_set.json", &jsonData)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	// t.Logf("Read() data = %v", jsonData)
	index := 0
	var tree KV_Tree
	root := tree.jsonToKVNode(jsonData, "root", &index)

	// Print KV_Node structure
	printKVNode(root, "",false)
	tree.jsonToKVNode(jsonData, "root", &index)
}
