package utils

import (
	"fmt"

	"./request"
)

func GetSearchSuggest(s string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/search/suggest/keyword"
	options["url"] = "https://music.163.com/eapi/search/suggest/keyword"
	var lastKeyword, lastKeywordForm string
	lastTime := 0
	Type := 1018
	limit := 10
	if _, ok := options["lastKeyword"]; ok {
		lastKeyword = options["lastKeyword"].(string)
	}
	if _, ok := options["lastTime"]; ok {
		lastTime = options["lastTime"].(int)
	}
	if _, ok := options["Type"]; ok {
		Type = options["Type"].(int)
	}
	if _, ok := options["limit"]; ok {
		limit = options["limit"].(int)
	}
	if lastKeyword != "" {
		lastKeywordForm = fmt.Sprintf("\"lastKeyword\":\"%v\",", lastKeyword)
	}
	options["str"] = fmt.Sprintf("{\"lastTime\":\"%v\",\"s\":\"%v\",\"type\":\"%v\",%v\"limit\":\"%v\",\"header\":\"{}\",\"e_r\":\"true\"}", lastTime, s, Type, lastKeywordForm, limit)
	result = request.EapiRequest(options)
	return result
}
