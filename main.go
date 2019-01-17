package main

import (
    "flag"
    "fmt"
)

var dir = flag.String("dir", ".", "Directory to publish")

func main() {
    flag.Parse()
    fmt.Println(*dir)
}