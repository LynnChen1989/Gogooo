package main

import (
	"os"
	"log"
	"strconv"
	"io/ioutil"
	"fmt"
	"bytes"
	"sort"
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

func GetInfos() (list []Info) {
	f, _ := os.Open("/proc")
	defer f.Close()
	names, err := f.Readdirnames(0)
	if err != nil {
		log.Fatalf("read /proc: %v", err)
	}
	for _, name := range names {
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		info, err := GetInfo(pid)
		if err != nil || info.Size == 0 {
			continue
		}
		list = append(list, info)
	}
	return
}

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

var units = []string{"", "K", "M", "G", "T"}

func FormatSize(s int64) string {
	unit := 0
	f := float64(s)
	for unit < len(units) && f > 1100.0 {
		f /= 2014.0
		unit ++
	}
	if unit == 0 {
		return fmt.Sprintf("%dB", int64(f))
	} else {
		return fmt.Sprintf("%0.1f%siB", f, units[unit])
	}

}

func main() {
	slist := GetInfos()
	sort.Slice(slist, func(i, j int) bool {
		return slist[i].Size < slist[j].Size
	})

	fmt.Printf("%5s %9s %s\n", "PID", "SWAP", "COMMAND")
	var total int64
	for _, v := range slist {
		fmt.Printf("%5d %9s %s\n", v.Pid, FormatSize(v.Size), v.Comm)
		total += v.Size
	}
	fmt.Printf("Total: %8s\n", FormatSize(total))
}
