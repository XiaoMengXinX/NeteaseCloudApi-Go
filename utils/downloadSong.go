package utils

import (
	"fmt"
	"path"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/downloader"
	log "github.com/sirupsen/logrus"
)

func DownloadSong(id string, options map[string]interface{}) (fileName []string) {
	var url, filename, savePath string
	savePath = "./"
	//startTime := time.Now()
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
			switch path.Ext(path.Base(url)) {
			case ".mp3":
				filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".mp3")
			case ".flac":
				filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".flac")
			case ".m4a":
				filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".m4a")
			default:
				filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".mp3")
			}
			fileName = append(fileName, filename)
			//filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),path.Ext(path.Base(url)))
			//downloader := downloader.NewFileDownloader(url, filename, savePath, threads, "")
			//if err := downloader.Run(); err != nil {
			//	log.Fatal(err)
			//}
			downloader := downloader.NewDownloader(savePath)
			downloader.AppendResource(filename, url)
			downloader.Concurrent = threads
			err := downloader.Start()
			if err != nil {
				log.Error(err)
			}
			//log.Printf("%v 下载完成 耗时: %f second\n", filename, time.Now().Sub(startTime).Seconds())
		} else {
			fileName = append(fileName, "null")
		}
	}
	return fileName
}

func MultiDownloadSong(ids []string, options map[string]interface{}) (fileName, validIds []string) {
	var id, url, filename, savePath string
	savePath = "./"
	threads := 4
	if _, ok := options["threads"].(int); ok {
		threads = options["threads"].(int)
	}
	if _, ok := options["savePath"].(string); ok {
		savePath = options["savePath"].(string)
	}
	downloader := downloader.NewDownloader(savePath)
	for i := 0; i < len(ids); i++ {
		id = ids[i]
		result := GetSongUrl(id, options)
		for _, v := range result["body"].(map[string]interface{})["data"].([]interface{}) {
			if v.(map[string]interface{})["url"] != nil {
				url = v.(map[string]interface{})["url"].(string)
				switch path.Ext(path.Base(url)) {
				case ".mp3":
					filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".mp3")
				case ".flac":
					filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".flac")
				case ".m4a":
					filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".m4a")
				default:
					filename = fmt.Sprintf("%v%v", int(v.(map[string]interface{})["id"].(float64)), ".mp3")
				}
				fileName = append(fileName, filename)
				validIds = append(validIds, fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64))))
				//filename = fmt.Sprintf("%v%v",int(v.(map[string]interface{})["id"].(float64)),path.Ext(path.Base(url)))
				//downloader := downloader.NewFileDownloader(url, filename, savePath, threads, "")
				//if err := downloader.Run(); err != nil {
				//	log.Fatal(err)
				//}
				downloader.AppendResource(filename, url)
				//log.Printf("%v 下载完成 耗时: %f second\n", filename, time.Now().Sub(startTime).Seconds())
			} else {
				fileName = append(fileName, "null")
			}
		}
	}
	downloader.Concurrent = threads
	err := downloader.Start()
	if err != nil {
		log.Error(err)
	}
	return fileName, validIds
}
