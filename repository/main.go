package repository

import (
	"flag"
	"fmt"
)

var dir = flag.String("dir", ".", "Directory to publish")


func Run() {
	flag.Parse()
	fmt.Println(*dir)
}