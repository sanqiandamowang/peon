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

// sortMapKeys 非递归地对 map 的键进行排序
func sortMapKeys(data interface{}) interface{} {
	type queueItem struct {
		key    string
		value  interface{}
		parent interface{}
		index  int
	}

	// 初始化队列，将根元素入队
	queue := []queueItem{{key: "", value: data, parent: nil, index: -1}}
	var result interface{}

	for len(queue) > 0 {
		// 取出队列的第一个元素
		current := queue[0]
		queue = queue[1:]

		switch v := current.value.(type) {
		case map[string]interface{}:
			// 对 map 的键进行排序
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			sortedMap := make(map[string]interface{})
			for _, k := range keys {
				sortedMap[k] = v[k]
				// 将 map 中的子元素加入队列
				queue = append(queue, queueItem{key: k, value: v[k], parent: sortedMap})
			}

			if current.parent == nil {
				// 如果没有父级，这是最顶层的 map，保存为结果
				result = sortedMap
			} else if pMap, ok := current.parent.(map[string]interface{}); ok {
				// 将排序后的 map 赋值回父 map
				pMap[current.key] = sortedMap
			} else if pSlice, ok := current.parent.([]interface{}); ok {
				// 将排序后的 map 赋值回父 slice
				pSlice[current.index] = sortedMap
			}

		case []interface{}:
			// 对 slice 进行非递归处理，将子元素加入队列
			for i := 0; i < len(v); i++ { // 正序入队，保持顺序
				queue = append(queue, queueItem{value: v[i], parent: v, index: i})
			}

			if current.parent == nil {
				// 如果没有父级，这是最顶层的 slice，保存为结果
				result = v
			}

		default:
			if current.parent == nil {
				// 如果没有父级，这是最顶层的值，保存为结果
				result = v
			} else if pMap, ok := current.parent.(map[string]interface{}); ok {
				// 将值赋回父 map
				pMap[current.key] = v
			} else if pSlice, ok := current.parent.([]interface{}); ok {
				// 将值赋回父 slice
				pSlice[current.index] = v
			}
		}
	}

	return result
}

//func sortMapKeys(data interface{}) interface{} {
//		switch v := data.(type) {
//		case map[string]interface{}:
//			// 对 map 的键进行排序
//			keys := make([]string, 0, len(v))
//			for k := range v {
//				keys = append(keys, k)
//			}
//			sort.Strings(keys)
//
//			sortedMap := make(map[string]interface{})
//			for _, k := range keys {
//				sortedMap[k] = sortMapKeys(v[k])
//			}
//			return sortedMap
//		case []interface{}:
//			// 对 slice 进行递归处理
//			for i, item := range v {
//				v[i] = sortMapKeys(item)
//			}
//			return v
//		default:
//			// 其他类型直接返回
//			return v
//		}
//	}
//}
