package utils

import "fmt"

func DownloadPlaylist(id string, options map[string]interface{}) {
	result := GetPlaylistDetail(id, options)
	if _, ok := options["s"].(int); ok {
		var mid string
		var i int = 0
		for t, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
			if i < options["s"].(int) {
				if i == 0 {
					mid = fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64)))
				} else {
					mid = fmt.Sprintf("%v,%v", mid, int(v.(map[string]interface{})["id"].(float64)))
				}
				if i == options["s"].(int)-1 {
					i = 0
					DownloadSong(mid, make(map[string]interface{}), options)
				} else {
					if len(result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}))-t == 1 {
						DownloadSong(mid, make(map[string]interface{}), options)
					} else {
						i++
					}
				}
			}
		}
	} else {
		for _, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
			var mid string = fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64)))
			//fmt.Println(mid)
			DownloadSong(mid, make(map[string]interface{}), options)
		}
	}
}
