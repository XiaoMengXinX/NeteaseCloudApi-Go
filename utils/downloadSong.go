package utils

import (
	"time"
	"fmt"
	"log"
	"path"

	"./downloader"
)

func DownloadSong(id string, options map[string]interface{}) {
	var url, filename string
	startTime := time.Now()
	result := GetSongUrl(id, options)
	threads := 4
	if _, ok := options["threads"].(int); ok {
		threads = options["threads"].(int)
	}
	for _, v := range result["body"].(map[string]interface{})["data"].([]interface{}) {
		if v.(map[string]interface{})["url"] != nil {
			url = v.(map[string]interface{})["url"].(string)
			filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),path.Ext(path.Base(url)))
			downloader := downloader.NewFileDownloader(url, filename, "", threads, "")
			if err := downloader.Run(); err != nil {
			log.Fatal(err)
			}
			fmt.Printf("%v 下载完成 耗时: %f second\n\n", filename, time.Now().Sub(startTime).Seconds())
		}
	}
}
