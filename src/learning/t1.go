package main

import (
	_ "fmt"
	_ "net/http"
	_ "net/url"
	"fmt"
	"encoding/json"
)

type Parameters struct {
	Parameter []ParameterFormat `json:"parameter"`
}

type ParameterFormat struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var paras []ParameterFormat

func JenkinsConstructParas(args map[string]string) {
	for k, v := range args {
		para := ParameterFormat{
			Name:  k,
			Value: v,
		}
		paras = append(paras, para)
	}
}

func GetPostParas(args map[string]string) (parameters Parameters) {
	JenkinsConstructParas(args)
	parameters = Parameters{
		Parameter: paras,
	}
	return
}

func main() {
	args := map[string]string{
		"name": "chenlin",
		"age":  "18",
		"sex":  "male",
	}
	data, _ := json.Marshal(GetPostParas(args))
	fmt.Println(string(data))
}
