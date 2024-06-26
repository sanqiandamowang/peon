package json

import "testing"

type test struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func TestJsonFileRead(t *testing.T) {
	var jsonData any
	err := Read("test.json", &jsonData)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	t.Logf("Read() data = %v", jsonData)
}

func TestJsonFileWrite(t *testing.T) {
	var data test
	data.Name = "Alice"
	data.Age = 25
	data.City = "Beijing"
	err := Write("test1.json", data)
	if err != nil {
		t.Errorf("Write() error = %v", err)
		return
	}	
}