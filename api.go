package main

import (
	"fmt"
	"os"

	"./utils"
)

func main() {
	arg := os.Args[1]
	var options map[string]interface{}
	options = make(map[string]interface{})
	var cookies map[string]interface{}
	cookies = make(map[string]interface{})
	cookies["MUSIC_U"] = ""
	options["cookie"] = cookies
	options["resultType"] = "json"
	//result := utils.GetSongDetail(data, options)
	result := utils.SearchSong(arg, options)
	status := result["status"].(int)
	data := result["body"]
	//data := result["body"].(map[string]interface{})["songs"].([]interface{})[0].(map[string]interface{})["name"]
	fmt.Printf("%d\n%s\n", status, data)

	//fmt.Printf("%d\n", status)

	//for _, v := range result["body"].(map[string]interface{})["songs"].([]interface{}) {
	//	fmt.Println(v.(map[string]interface{})["name"])
	//}

	//walk(result["body"])
}

func walk(v interface{}) {
	switch v := v.(type) {
	case []interface{}:
		for i, v := range v {
			fmt.Println(i)
			walk(v)
		}
	case map[string]interface{}:
		for k, v := range v {
			fmt.Println(k)
			walk(v)
		}
	default:
		fmt.Println(v)
	}
}
