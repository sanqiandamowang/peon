package jsonutils

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
)

func Read(filePath string, v interface{}) error {
	// 读取文件
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取%s文件时发生错误: %w", filePath, err)
	}

	// 解码JSON数据到结构体
	err = sonic.Unmarshal(jsonData, v)
	if err != nil {
		return fmt.Errorf("解码%s文件时发生错误: %w", filePath, err)
	}

	return nil
}
