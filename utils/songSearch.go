package utils

import (
	"fmt"

	"./request"
)

func SearchSong(keywords string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v1/search/song/get"
	options["url"] = "https://music.163.com/eapi/v1/search/song/get"
	limit := 10
	offset := 0
	if _, ok := options["limit"]; ok {
		limit = options["limit"].(int)
	}
	if _, ok := options["offset"]; ok {
		limit = options["limit"].(int)
	}
	options["str"] = fmt.Sprintf("{\"sub\":\"false\",\"s\":\"%v\",\"offset\":\"%v\",\"limit\":\"%v\",\"queryCorrect\":\"true\",\"strategy\":\"5\",\"header\":\"{}\",\"e_r\":\"true\"}", keywords, offset, limit)
	result = request.EapiRequest(options)
	return result
}
