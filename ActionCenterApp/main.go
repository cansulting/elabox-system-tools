package actioncenterapp

import (
	"fmt"
	"time"
)

var running bool = true

func main() {
	fmt.Println("Welcome")
	go RunServer()
	for running {
		time.Sleep(time.Second * 1)
	}
}
