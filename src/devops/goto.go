package src

import (
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/user"
	"os/exec"
	"regexp"
	"syscall"
	"io/ioutil"
	"strings"
	"time"
	"errors"
	"path/filepath"
)

type HostInfo struct {
	Host     []string          `json:"hosts"`
	Children []string          `json:"children"`
	Vars     map[string]string `json:"vars"`
}

type Appinfo map[string]interface{}

var cacheExpireSeconds int64 = 7200
var auth_token = "Token  " + os.Getenv("CMDB_API_TOKEN")
var listApp bool
var listHost string
var clean_cache bool
var show_completion bool
var current_user, _ = user.Current()
var appCacheFile = "/tmp/.goto." + current_user.Username + ".appcache"
var hostListCacheFilePrefix = "/tmp/.goto." + current_user.Username + ".hostcache_"

func main() {
	flag.BoolVar(&listApp, "listapp", false, "")
	flag.BoolVar(&clean_cache, "clean", false, "")
	flag.BoolVar(&show_completion, "completion", false, "")
	flag.StringVar(&listHost, "listhost", "", "")
	flag.Parse()

	if clean_cache {
		files, err := filepath.Glob("/tmp/.goto." + current_user.Username + ".*")
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

	if show_completion {
		PrintCompletion()
		return
	}

	if listApp {
		fmt.Println(GetAppList())
		return
	}

	if listHost != "" {
		fmt.Println(GetHostList(listHost))
		return
	}

	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		panic(lookErr)
	}

	ssh_user := os.Getenv("DEFAULT_GOTO_USER")
	if ssh_user == "" {
		ssh_user = current_user.Username
	}
	if len(os.Args) > 3 {
		ssh_user = os.Args[3]
	}
	port := "22"
	if len(os.Args) > 4 {
		port = os.Args[4]
	}
	args := []string{"ssh", "-p", port, ssh_user + "@" + getIp(os.Args[2])}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func usage() {
	//fmt.Println("usage: goto <app_code> <host> [user] [port]")
	fmt.Println(`usage: goto <app_code> <host> [user] [port]
  options:
    --help			show this page
    --listapp			list all your apps' code
    --listhost <pattern>	list all the hosts match the pattern
    --clean			clean cache file under /tmp dir
    --completion		show bash completion script`)
}

func PrintCompletion() {
	fmt.Println(`_goto()
{
case $COMP_CWORD in
    0)
        : ;;
    1) 
        local cur="${COMP_WORDS[COMP_CWORD]}"
        LIST_CMD="goto"
        local APP_LIST=$(${LIST_CMD} -listapp)
        COMPREPLY=( $(compgen -W "${APP_LIST}" -- ${cur}) )
        ;;
    2)
        local cur="${COMP_WORDS[COMP_CWORD]}"
        LIST_CMD="goto"
        local SERVER_LIST=$(${LIST_CMD} -listhost=${COMP_WORDS[1]})
        COMPREPLY=( $(compgen -W "${SERVER_LIST}" -- ${cur}) )
        ;;
    esac
}
complete -o filenames -o nospace -o bashdefault -F _goto goto`)
}

func getIp(hostname string) (ip string) {
	reg := regexp.MustCompile(`[\d]+\.[\d]+\.[\d]+\.[\d]+$`)
	ip = reg.FindString(hostname)
	//fmt.Println(ip)
	return
}

func GetAppList() (appCodeList string) {
	var _codeArray []string
	appCodeList, err := ReadCache(appCacheFile)
	if err != nil {
		appInfoList := []Appinfo{}
		request := gorequest.New()
		request.Get("http://cmdb-api.ops.dragonest.com/api/apps").Set("Authorization", auth_token).EndStruct(&appInfoList)
		for _, v := range appInfoList {
			_codeArray = append(_codeArray, v["code"].(string))
		}
		//fmt.Println(resp.Body)
		appCodeList = strings.Join(_codeArray, "\n")
		WriteCache(appCacheFile, appCodeList)
	}
	return
}

func GetHostList(app_code string) (hostList string) {
	var _hostArray []string
	hostList, err := ReadCache(hostListCacheFilePrefix + app_code)
	if err != nil {
		result := map[string]HostInfo{}
		request := gorequest.New()
		request.Get("http://cmdb-api.ops.dragonest.com/ops/inventory/" + app_code + "/*/?use_node_code_in_host=true&use_public_ip=true").Set("Authorization", auth_token).EndStruct(&result)
		for k, v := range result {
			if k != "_meta" {
				_hostArray = append(_hostArray, v.Host...)
			}
		}
		hostList = strings.Join(_hostArray, "\n")
		WriteCache(hostListCacheFilePrefix+app_code, hostList)
	}
	return
}

func WriteCache(outputFile string, data string) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	err := ioutil.WriteFile(outputFile, []byte(data), 0600) // oct, not hex
	if err != nil {
		panic(err.Error())
	}
}

func ReadCache(inputFile string) (data string, err error) {
	old_mask := syscall.Umask(0) // set new mask, return old mask
	defer syscall.Umask(old_mask)
	f, err := os.OpenFile(inputFile, os.O_RDWR, 0600)
	if err != nil {
		return
	}
	cache_info, err := f.Stat()
	if err != nil {
		return
	}

	if time.Now().Unix()-cache_info.ModTime().Unix() > cacheExpireSeconds {
		err = errors.New("cache expired")
		return
	}

	cache_data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return
	}
	return string(cache_data), err
}
