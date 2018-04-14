package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Info struct {
	Pid  int
	Size int64
	Comm string
}

var (
	nullBytes  = []byte{0x0}
	emptyBytes = []byte(" ")
)

func GetInfo(pid int) (info Info, err error) {
	info.Pid = pid
	var bs []byte
	bs, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return
	}
	if bytes.HasSuffix(bs, nullBytes) {
		bs = bs[:len(bs)-1]
	}

	info.Comm = string(bytes.Replace(bs, nullBytes, emptyBytes, -1))
	bs, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return
	}
	var total int64
	for _, line := range bytes.Split(bs, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("Swap:")) {
			start := bytes.IndexAny(line, "0123456789")
			end := bytes.Index(line[start:], []byte(" "))
			size, err := strconv.ParseInt(string(line[start:start+end]), 10, 0)
			if err != nil {
				continue
			}
			total += size
		}
	}
	info.Size = total * 1024
	return
}

func main() {

}
