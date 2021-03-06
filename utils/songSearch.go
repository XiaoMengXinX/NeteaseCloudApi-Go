package utils

import (
	"fmt"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/request"
)

func SearchSong(keywords string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v1/search/song/get"
	options["url"] = "https://music.163.com/eapi/v1/search/song/get"
	if _, ok := options["disable_https"].(bool); ok {
		if options["disable_https"].(bool) {
			options["url"] = "http://music.163.com/eapi/v1/search/song/get"
		}
	}
	limit := 10
	offset := 0
	if _, ok := options["limit"].(int); ok {
		limit = options["limit"].(int)
	}
	if _, ok := options["offset"].(int); ok {
		offset = options["offset"].(int)
	}
	options["str"] = fmt.Sprintf("{\"sub\":\"false\",\"s\":%q,\"offset\":\"%v\",\"limit\":\"%v\",\"queryCorrect\":\"true\",\"strategy\":\"5\",\"header\":\"{}\",\"e_r\":\"true\"}", keywords, offset, limit)
	result = request.EapiRequest(options)
	return result
}
