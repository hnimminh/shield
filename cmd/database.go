package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type Rule struct{
	SrcIP []string `json:"src_ip"`
    DstIP []string `json:"dst_ip"`
    SrcPort []int `json:"src_port"`
	DstPort []int `json:"dst_port"`
	Transport string `json:"transport"`
	Action string `json:"action"`
    Comment string `json:"comment"`
}

type Security struct{
    Builtin bool `json:"builtin"`
    Rules []Rule `json:"rules"`
}



func main(){
	jsonFile, err := os.Open("data.json")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("succes open file ", jsonFile)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	//var result map[string]interface{}
	var sec Security

    json.Unmarshal(byteValue, &sec)

	fmt.Println(sec)
}
