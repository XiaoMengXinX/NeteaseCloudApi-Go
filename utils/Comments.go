package utils

import (
	"fmt"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/request"
)

func GetComments(id string, commentType int, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v2/resource/comments"
	options["url"] = "https://music.163.com/eapi/v2/resource/comments"
	var cursor, pageNo, pageSize, sortType int = 0, 1, 20, 0
	var resourceType string
	if _, ok := options["cursor"].(int); ok {
		cursor = options["cursor"].(int)
	}
	if _, ok := options["pageNo"].(int); ok {
		pageNo = options["pageNo"].(int)
	}
	if _, ok := options["pageSize"].(int); ok {
		pageSize = options["pageSize"].(int)
	}
	if _, ok := options["sortType"].(int); ok {
		sortType = options["sortType"].(int)
	}
	switch commentType {
	case 0:
		resourceType = "R_SO_4" //歌曲
	case 1:
		resourceType = "R_MV_5" //MV
	case 2:
		resourceType = "A_PL_0" //歌单
	case 3:
		resourceType = "R_AL_3" //专辑
	case 4:
		resourceType = "A_DJ_1" //电台
	case 5:
		resourceType = "R_VI_62" //视频
	case 6:
		resourceType = "R_MLOG_1001" //Mlog
	}
	options["decrypt"] = 0
	options["str"] = fmt.Sprintf("{\"cursor\":\"%v\",\"pageNo\":%v,\"pageSize\":%v,\"showInner\":false,\"sortType\":%v,\"threadId\":\"%v_%v\"}", cursor, pageNo, pageSize, sortType, resourceType, id)
	result = request.EapiRequest(options)
	return result
}

func GetSongComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 0, options)
}

func GetMVComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 1, options)
}

func GetPlaylistComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 2, options)
}

func GetAlbumComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 3, options)
}

func GetDJComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 4, options)
}

func GetVideoComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 5, options)
}

func GetMlogComments(id string, options map[string]interface{}) (result map[string]interface{}) {
	return GetComments(id, 6, options)
}

// 我是看不懂网易云这破烂api到底啥玩意了
func GetEventComments(id, userid string, options map[string]interface{}) (result map[string]interface{}) {
	options["path"] = "/api/v1/resource/comments/A_EV_2_" + id + "_" + userid
	options["url"] = "https://music.163.com/eapi/v1/resource/comments/A_EV_2_" + id
	var pageNo, pageSize int = 1, 20
	if _, ok := options["pageNo"].(int); ok {
		pageNo = options["pageNo"].(int)
	}
	if _, ok := options["pageSize"].(int); ok {
		pageSize = options["pageSize"].(int)
	}
	options["decrypt"] = 0
	options["str"] = fmt.Sprintf("{\"resourceId\":\"%v\",\"resourceType\":\"2\",\"limit\":\"%v\",\"beforeTime\":\"0\",\"compareUserLocation\":\"true\",\"showInner\":false,\"pageNum\":%v}", id, pageSize, pageNo)
	result = request.EapiRequest(options)
	return result
}