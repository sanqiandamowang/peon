package jsonutils

import (
	"os"

	"github.com/bytedance/sonic"
)

func Write(filePath string, v interface{}) error {
	jsonBytes, err := sonic.Marshal(v)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
	_, err = file.Write(jsonBytes)
    if err != nil {
        return err
    }
    
    return nil
}
