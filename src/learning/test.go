package main

import "os"

type Movie struct {
	Title   string
	Year    int    `json:"released"`
	Color   string `json:"color, omitempty"`
	Authors []string
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

func main() {
	os.Stdout.Write([]byte("hello"))
}
