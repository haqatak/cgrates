package main

import (
    "fmt"
    "github.com/cgrates/cgrates/utils"
)

func main() {
    a := make(utils.StringMap)
    a["test"] = true
    fmt.Printf("%v\n", a)
}
