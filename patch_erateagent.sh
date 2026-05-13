cat << 'EOF2' > patch_bon_voyage_event.go
package main

import (
	"io/ioutil"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("agents/erateagent.go")
	if err != nil {
		panic(err)
	}
	s := string(b)

	s = strings.Replace(s, "// TODO: Create engine event here if required. For now, just return success.", "", 1)

	ioutil.WriteFile("agents/erateagent.go", []byte(s), 0644)
}
EOF2
go run patch_bon_voyage_event.go
rm patch_bon_voyage_event.go
