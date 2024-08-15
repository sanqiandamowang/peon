package kvtree

import (
	"peon/utils/jsonutils"
	"testing"
)

func TestLoad_v2(t *testing.T) {
	var source map[string]interface{}
	err := jsonutils.Read("/home/sanqian/peon/peon/testValue/test.json", &source)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	tree:= &KV_Tree_V2{}
	tree.Source = &source
	tree.NodeList = tree.SourceToKVNode(source, "root", nil)
	tree.NodeList.Value.(map[string]interface{})["step1"].(map[string]interface{})["step2"] =2
	tree.NodeList.Value.(map[string]interface{})["test"].([]interface{})[0] ="wa"
	tree.NodeList.Value.(map[string]interface{})["test2"].([]interface{})[1].(map[string]interface{})["b"] = "sa"
	tree.NodeList.Value.(map[string]interface{})["test2"].([]interface{})[1].(map[string]interface{})["c"] = "sa"
	tree.printKVNode(tree.NodeList, "", false)
}
