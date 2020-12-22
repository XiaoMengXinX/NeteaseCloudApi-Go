package utils

import "fmt"

func DownloadPlaylist(id string, options map[string]interface{}) {
	result := GetPlaylistDetail(id, options)
	for _, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
		var mid string
		mid = fmt.Sprintf("%v",int(v.(map[string]interface{})["id"].(float64)))
		//fmt.Println(mid)
		DownloadSong(mid, options)
	}
}
