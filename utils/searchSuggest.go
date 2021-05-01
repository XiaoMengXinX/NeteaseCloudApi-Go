package utils

import (
	"fmt"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/request"
)

func GetSearchSuggest(s string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/search/suggest/keyword"
	options["url"] = "https://music.163.com/eapi/search/suggest/keyword"
	if _, ok := options["disable_https"].(bool); ok {
		if options["disable_https"].(bool) {
			options["url"] = "http://music.163.com/eapi/search/suggest/keyword"
		}
	}
	var lastKeyword, lastKeywordForm string
	lastTime := 0
	Type := 1018
	limit := 10
	if _, ok := options["lastKeyword"].(string); ok {
		lastKeyword = options["lastKeyword"].(string)
	}
	if _, ok := options["lastTime"].(int); ok {
		lastTime = options["lastTime"].(int)
	}
	if _, ok := options["Type"].(int); ok {
		Type = options["Type"].(int)
	}
	if _, ok := options["limit"].(int); ok {
		limit = options["limit"].(int)
	}
	if lastKeyword != "" {
		lastKeywordForm = fmt.Sprintf("\"lastKeyword\":\"%v\",", lastKeyword)
	}
	options["str"] = fmt.Sprintf("{\"lastTime\":\"%v\",\"s\":%q,\"type\":\"%v\",%v\"limit\":\"%v\",\"header\":\"{}\",\"e_r\":\"true\"}", lastTime, s, Type, lastKeywordForm, limit)
	result = request.EapiRequest(options)
	return result
}
