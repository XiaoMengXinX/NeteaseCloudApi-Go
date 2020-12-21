package main

import (
	"log"
	"time"
	"fmt"

	"./downloader"
)

func main() {
	startTime := time.Now()
	url := "http://m7.music.126.net/20201222001732/4dec707db0ecb1e93af82b050f9833ae/ymusic/025e/560c/550f/549b64866caf78457292bec989ffaa95.flac"
	downloader := NewFileDownloader(url, "", "", 8, "")
	if err := downloader.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n 文件下载完成耗时: %f second\n", time.Now().Sub(startTime).Seconds())
}