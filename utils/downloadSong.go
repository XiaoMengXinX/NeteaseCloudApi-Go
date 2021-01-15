package utils

import (
	"time"
	"fmt"
	"log"
	"path"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/downloader"
)

func DownloadSong(id string, options map[string]interface{}) (fileName []string) {
	var url, filename, savePath string
	savePath = "./"
	startTime := time.Now()
	result := GetSongUrl(id, options)
	threads := 4
	if _, ok := options["threads"].(int); ok {
		threads = options["threads"].(int)
	}
	if _, ok := options["savePath"].(string); ok {
		savePath = options["savePath"].(string)
	}
	for _, v := range result["body"].(map[string]interface{})["data"].([]interface{}) {
		if v.(map[string]interface{})["url"] != nil {
			url = v.(map[string]interface{})["url"].(string)
			switch path.Ext(path.Base(url)){
    		case ".mp3":
        		filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),".mp3")
			case ".flac":
        		filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),".flac")
        	case ".m4a":
        		filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),".m4a")
    		default:
        		filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),".mp3")
			}
			fileName = append(fileName, filename)
			//filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),path.Ext(path.Base(url)))
			downloader := downloader.NewFileDownloader(url, filename, savePath, threads, "")
			if err := downloader.Run(); err != nil {
			log.Fatal(err)
			}
			fmt.Printf("%v 下载完成 耗时: %f second\n\n", filename, time.Now().Sub(startTime).Seconds())
		} else {
			fileName = append(fileName, "null")
		}
	}
	return fileName
}
