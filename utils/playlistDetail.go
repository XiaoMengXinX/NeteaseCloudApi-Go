package utils

import (
	"fmt"

	"./request"
)

func GetPlaylistDetail(id string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v6/playlist/detail"
	options["url"] = "https://interface3.music.163.com/eapi/v6/playlist/detail"
	s := 8
	if _, ok := options["s"].(int); ok {
		s = options["s"].(int)
	}
	options["str"] = fmt.Sprintf("{\"id\":\"%v\",\"t\":\"0\",\"n\":\"1000\",\"s\":\"%v\",\"shareUserId\":\"0\",\"header\":\"{}\",\"e_r\":\"true\"}", id, s)
	result = request.EapiRequest(options)
	return result
}
