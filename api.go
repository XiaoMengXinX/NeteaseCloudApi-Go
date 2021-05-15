package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	utils "github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils"
	log "github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetFormatter(new(LogFormatter))
	log.SetLevel(log.InfoLevel)
}

func main() {
	arg := os.Args[1]
	var options map[string]interface{}
	options = make(map[string]interface{})
	var cookies map[string]interface{}
	cookies = make(map[string]interface{})
	config, err := utils.ReadConfig("./config.ini")
	if err != nil {
		log.Fatal(err)
	}
	cookies["MUSIC_U"] = config["MUSIC_U"]
	options["cookie"] = cookies
	options["s"] = 5
	//options["savePath"] = "./download"
	//options["limit"] = 2
	//options["resultType"] = "json"
	//result := utils.DownloadSong(arg, options)
	//utils.DownloadSong(arg, options)
	//result := utils.GetSongDetail(arg, options)
	result := utils.GetPlaylistDetail(arg, options)
	//status := result["status"].(int)
	data := result["body"]
	//data := result["body"].(map[string]interface{})["songs"].([]interface{})[0].(map[string]interface{})["al"].(map[string]interface{})["picUrl"]
	data, _ = json.Marshal(data)
	fmt.Printf("%s\n", data)

	//fmt.Printf("%d\n", status)
	//var i int = 0
	//for _, v := range result["body"].(map[string]interface{})["data"].(map[string]interface{})["comments"].([]interface{}) {
	//	i++
	//	fmt.Println(i, v.(map[string]interface{})["content"], "\n")
	//fmt.Println(int(v.(map[string]interface{})["id"].(float64)))
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
