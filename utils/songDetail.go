package utils

import (
	"fmt"
	"strings"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/request"
)

func GetSongDetail(id string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v3/song/detail"
	options["url"] = "https://music.163.com/eapi/v3/song/detail"
	ids := strings.Split(id, ",")
	var data string
	if len(ids) > 0 {
		data = fmt.Sprintf("{\\\"id\\\":%v,\\\"v\\\":0}", ids[0])
		for i:=1; i < len(ids); i++ {
			data = fmt.Sprintf("%v,{\\\"id\\\":%v,\\\"v\\\":0}", data, ids[i])
		}
	}
	options["str"] = fmt.Sprintf("{\"c\":\"[%v]\",\"e_r\":\"true\",\"header\":\"{}\"}", data)
	result = request.EapiRequest(options)
	return result
}
