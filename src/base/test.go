package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	meta := make(map[string]interface{})
	hostVar := make(map[string]interface{})
	hostVar["hostvars"] = make(map[string]interface{})
	meta["_meta"] = hostVar
	b, _ := json.Marshal(meta)
	fmt.Println(string(b))
}
