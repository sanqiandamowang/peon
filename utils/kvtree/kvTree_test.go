package kvtree

import (
	"peon/utils/jsonutils"
	"testing"
)

func TestLoad(t *testing.T) {
	var source map[string]interface{}
	err := jsonutils.Read("/home/sanqian/peon/peon/testValue/dev_set.json", &source)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	// t.Logf("Read() data = %v", jsonData)

	// sortedMap := json.SortJSON(source)
	// sortedJSON, err := sonic.ConfigDefault.MarshalIndent(sortedMap, "", "  ")
	// if err != nil {
	// 	t.Errorf("MarshalIndent() error = %v", err)
	// }
	// t.Log(string(sortedJSON))
	// index := 0
	var tree KV_Tree
	
	tree.Load("/home/sanqian/peon/peon/testValue/dev_set.json",source)
	
	// Print KV_Node structure
	tree.printKVNode(tree.NodeList, "", false)
}
