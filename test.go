package main

import (
    "fmt"
    "github.com/cgrates/cgrates/utils"
)

func main() {
    id := utils.ConcatenatedKey("tenant", "subject")
    splt := utils.SplitConcatenatedKey(id)
    fmt.Printf("%v\n", splt)
    if len(splt) > 1 {
        fmt.Printf("%v\n", splt[1])
    }
}
