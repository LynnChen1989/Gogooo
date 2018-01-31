package main

import (
	"fmt"
	"strings"
)

func main() {

	fmt.Printf("[%q]", strings.Trim(" !!! Achtung hahh !!! ", "! ")) // ["Achtung"]
}
