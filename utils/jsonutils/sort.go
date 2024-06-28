package jsonutils

import (
	"encoding/json"
	"sort"

	// "github.com/bytedance/sonic"
)

// / SortJSON 对传入的 map[string]interface{} 进行排序并返回排序后的 map
func SortJSON(jsonStr string) (string, error) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return "", err
	}

	sortedData := sortMapKeys(jsonData)

	sortedJSON, err := json.MarshalIndent(sortedData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(sortedJSON), nil
}

// sortMapKeys 递归地对 map 的键进行排序
func sortMapKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// 对 map 的键进行排序
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		sortedMap := make(map[string]interface{})
		for _, k := range keys {
			sortedMap[k] = sortMapKeys(v[k])
		}
		return sortedMap
	case []interface{}:
		// 对 slice 进行递归处理
		for i, item := range v {
			v[i] = sortMapKeys(item)
		}
		return v
	default:
		// 其他类型直接返回
		return v
	}
}
