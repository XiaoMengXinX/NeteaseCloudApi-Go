package request

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	netUrl "net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/crypt"
)

func EapiRequest(options map[string]interface{}) (result map[string]interface{}) {
	data := SpliceStr(options["path"].(string), options["str"].(string))
	answer := CreateNewRequest(Format2Params(data), options["url"].(string), options)
	result = map[string]interface{}{
		"status": 500,
		"body":   map[string]interface{}{},
	}
	if answer["status"].(int) == 200 {
		var decrypted []byte
		if options["decrypt"] != 0 {
			decrypted = crypt.AesDecryptECB(answer["body"].([]byte))
		} else {
			decrypted = answer["body"].([]byte)
		}
		if _, ok := options["resultType"]; ok {
			if options["resultType"] == "json" {
				result["status"] = answer["status"].(int)
				result["body"] = string(decrypted)
			} else {
				bodyJson := map[string]interface{}{}
				if err := json.Unmarshal(decrypted, &bodyJson); err == nil {
					result["body"] = bodyJson
					if _, ok := bodyJson["code"]; ok {
						result["status"] = int(bodyJson["code"].(float64))
					} else {
						result["status"] = answer["status"].(int)
					}
				}
			}
		} else {
			bodyJson := map[string]interface{}{}
			if err := json.Unmarshal(decrypted, &bodyJson); err == nil {
				result["body"] = bodyJson
				if _, ok := bodyJson["code"]; ok {
					result["status"] = int(bodyJson["code"].(float64))
				} else {
					result["status"] = answer["status"].(int)
				}
			}
		}
	}
	return result
}

func SpliceStr(path string, data string) (result string) {
	text := fmt.Sprintf("nobody%suse%smd5forencrypt", path, data)
	md5 := md5.Sum([]byte(text))
	md5str := fmt.Sprintf("%x", md5)
	result = fmt.Sprintf("%s-36cd479b6b5-%s-36cd479b6b5-%s", path, data, md5str)
	return result
}

func Format2Params(str string) (data string) {
	data = fmt.Sprintf("params=%X", crypt.AesEncryptECB(str))
	return data
}

func ChooseUserAgent() string {
	userAgentList := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
		"NeteaseMusic/6.5.0.1575377963(164);Dalvik/2.1.0 (Linux; U; Android 9; MIX 2 MIUI/V12.0.1.0.PDECNXM)",
	}
	rand.Seed(time.Now().UnixNano())
	var index int
	index = rand.Intn(len(userAgentList))
	return userAgentList[index]
}

func encodeURIComponent(str string) string {
	r := netUrl.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func CreateNewRequest(data string, url string, options map[string]interface{}) map[string]interface{} {

	answer := map[string]interface{}{
		"status": 500,
		"body":   "string",
	}

	client := &http.Client{}
	reqBody := strings.NewReader(data)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		answer["status"] = 502
		answer["body"] = map[string]interface{}{
			"code": 502,
			"msg":  err.Error(),
		}
		return answer
	}

	value, ok := options["cookie"]
	if ok && reflect.ValueOf(value).Kind() == reflect.Map {
		Value, ok := options["cookie"].(map[string]interface{})
		cookie := map[string]interface{}{}
		if ok {
			cookie = Value
		}

		csrfValue, isok := cookie["__csrf"]
		csrfToken := ""
		if isok {
			csrfToken = fmt.Sprintf("%v", csrfValue)
		}
		header := make(map[string]interface{})
		keys := [...]string{"osver", "deviceId", "mobilename", "channel"}
		for _, val := range keys {
			value, ok := cookie[val]
			if ok {
				header[val] = value
			}
		}
		header["appver"] = func() string {
			val, ok := cookie["appver"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return "6.5.0"
		}()
		header["versioncode"] = func() string {
			val, ok := cookie["versioncode"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return "164"
		}()
		header["buildver"] = func() string {
			val, ok := cookie["buildver"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return strconv.FormatInt(time.Now().Unix(), 10)[0:10]
		}()
		header["resolution"] = func() string {
			val, ok := cookie["resolution"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return "1920x1080"
		}()
		header["os"] = func() string {
			val, ok := cookie["os"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return "android"
		}()
		header["__csrf"] = csrfToken
		cookieMusicU, ok := cookie["MUSIC_U"]
		if ok {
			header["MUSIC_U"] = cookieMusicU
		}
		cookieMusicA, ok := cookie["MUSIC_A"]
		if ok {
			header["MUSIC_A"] = cookieMusicA
		}

		cookies := ""
		for key, val := range header {
			cookies += encodeURIComponent(key) + "=" + encodeURIComponent(fmt.Sprintf("%v", val)) + "; "
		}
		req.Header.Set("Cookie", strings.TrimRight(cookies, "; "))
	} else if ok {
		req.Header.Set("Cookie", options["cookie"].(string))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ChooseUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		answer["status"] = 502
		answer["body"] = map[string]interface{}{
			"code": 502,
			"msg":  err.Error(),
		}
		return answer
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		answer["status"] = 502
		answer["body"] = map[string]interface{}{
			"code": 502,
			"msg":  err.Error(),
		}
		return answer
	}
	answer["body"] = []byte(body)
	answer["status"] = resp.StatusCode
	return answer
}
