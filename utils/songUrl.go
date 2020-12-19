package utils

import (
	"fmt"
	"strings"

	"./request"git 
)

func GetSongUrl(id string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/song/enhance/player/url/v1"
	options["url"] = "https://music.163.com/eapi/song/enhance/player/url/v1"
	encodeType := "mp3"
	level := "lossless"
	if _, ok := options["encodeType"]; ok {
		encodeType = options["encodeType"].(string)
	}
	if _, ok := options["level"]; ok {
		level = options["level"].(string)
	}
	ids := strings.Split(id, ",")
	var data string
	if len(ids) > 0 {
		data = fmt.Sprintf("\\\"%v\\\"", ids[0])
		for i := 1; i < len(ids); i++ {
			data = fmt.Sprintf("%v,\\\"%v\\\"", data, ids[i])
		}
	}
	options["str"] = fmt.Sprintf("{\"e_r\":\"true\",\"encodeType\":\"%v\",\"header\":\"{}\",\"ids\":\"[%v]\",\"level\":\"%v\"}", encodeType, data, level)
	result = request.EapiRequest(options)
	return result
}
