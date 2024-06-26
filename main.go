package main

import (

	"peon/gui"
	
)

func main() {
	err := gui.LoadConfig()
	if err != nil {
		panic (err)
	}
	gui.DisBase()
}
