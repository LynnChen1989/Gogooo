package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

var showCompletion bool
var listHost string
var hostListCacheFilePrefix = "/tmp/.goo.hostcache_"
var cacheExpireSeconds int64 = 86400
var cleanCache bool

type Hosts struct {
	CloudName  string `json:"cloud_name"`
	Hostname   string `json:"hostname"`
	Region     string `json:"region"`
	IpPrivate  string `json:"ip_private"`
	IpPublic   string `json:"ip_public"`
	Sn         string `json:"sn"`
	ServerType string `json:"server_type"`
	ServerAttr string `json:"server_attr"`
}

func writeCache(outputFile string, data string) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	err := ioutil.WriteFile(outputFile, []byte(data), 0777) // oct, not hex
	if err != nil {
		panic(err.Error())
	}
}

func readCache(inputFile string) (data string, err error) {
	oldMask := syscall.Umask(0) // set new mask, return old mask
	defer syscall.Umask(oldMask)
	f, err := os.OpenFile(inputFile, os.O_RDWR, 0777)
	if err != nil {
		return
	}
	cacheInfo, err := f.Stat()
	if err != nil {
		return
	}

	if time.Now().Unix()-cacheInfo.ModTime().Unix() > cacheExpireSeconds {
		err = errors.New("cache expired")
		return
	}

	cacheData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return
	}
	return string(cacheData), err
}

func printCompletion() {
	fmt.Println(`_goo()
		{
			case $COMP_CWORD in
				0)
					: ;;
				1)
					local cur="${COMP_WORDS[COMP_CWORD]}"
					LIST_CMD="goo"
					local SERVER_LIST=$(${LIST_CMD} -listhost=${COMP_WORDS[1]})
					COMPREPLY=( $(compgen -W "${SERVER_LIST}" -- ${cur}) )
					;;
			esac
        }
		complete -o filenames -o nospace -o bashdefault -F _goo goo`)
}

func httpGet(url string, header map[string]string) (content string) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fatal Error:", err.Error())
		os.Exit(0)
	}
	for key, value := range header {
		req.Header.Add(key, value)
	}

	response, err := client.Do(req)
	defer response.Body.Close()

	for _, cookie := range response.Cookies() {
		fmt.Println("cookie:", cookie)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Get response error", err)
		os.Exit(0)
	}
	content = string(body)
	return
}

func usage() {
	fmt.Println(`usage: goo <host>
  	options:
	--help			show this page
	--listhost <pattern>	list all the hosts match the pattern
	--completion		show bash completion script
	--clean			clean cache file under /tmp dir`)

}

func getHost() (hosts []Hosts) {
	header := map[string]string{
		"Authorization": "token 92dd208254c7692690a582179bb3bdb4bdcb3f47",
	}
	url := "http://cmdb.xwfintech.com/api/v1/hosts/"
	var data = httpGet(url, header)
	json.Unmarshal([]byte(data), &hosts)
	return
}

func getHostWithIp() (hostIp string) {
	hostIp, err := readCache(hostListCacheFilePrefix)
	if err != nil {
		hosts := getHost()
		for _, v := range hosts {
			hostIp += v.Hostname + "-" + getIp(v.IpPrivate) + "\n"
		}
	}
	fmt.Println(hostIp)
	writeCache(hostListCacheFilePrefix, hostIp)
	return
}

func getIp(hostname string) (ip string) {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	ipPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock
	reg := regexp.MustCompile(ipPattern)
	ip = reg.FindString(hostname)
	return
}

func main() {
	flag.BoolVar(&showCompletion, "completion", false, "")
	flag.StringVar(&listHost, "listhost", "", "")
	flag.BoolVar(&cleanCache, "clean", false, "")
	flag.Parse()

	if cleanCache {
		files, err := filepath.Glob("/tmp/.goo.hostcache_")
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}
		return
	}
	if showCompletion {
		printCompletion()
		return
	}
	if listHost != "" {
		fmt.Println(getHostWithIp())
		return
	}
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		panic(lookErr)
	}
	sshUser := os.Getenv("DEFAULT_GOTO_USER")
	if sshUser == "" {
		sshUser = "root"
	}
	port := "22"
	args := []string{"ssh", "-p", port, sshUser + "@" + getIp(os.Args[1])}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
