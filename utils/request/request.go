package main

import (
	"fmt"
	"time"
	"math/rand"
	"net/http"
	"io/ioutil"
	"strings"
	"reflect"
	netUrl "net/url"
	"strconv"
	"../crypt"
)

func main() {
	var url, data string
	url = "https://music.163.com/eapi/v3/song/detail"
	data = "params=11675D69CF25E0559750EF4BF81ECC680607C372872A3914F346ED123220A487D1D3764D2399095F2190002208159F14DA0AD23A941FD54017114FF3AB4EAE6B4EAD2EE19EA6C185932E6873F9AB06FFAC684A83B1169F65F7C61A8CFAE6442A2EBF137BF14B07E875E3E7E6F1163A2239A4AEABD25750A7DDCA71039E215F331F576F421BE7F88C98CCE22DD3ABB6D6"
	var options map[string]interface{}
    options = make(map[string]interface{})
    var cookies map[string]interface{}
    cookies = make(map[string]interface{})
    cookies["buildver"] = "1575377963"
    cookies["resolution"] = "2030x1080"
    cookies["appver"] = "6.5.0"
    cookies["MUSIC_U"] = "984e8c072dc9c670f40d019a3699f326b07414de7fe3522d93c1dd4fdb7286b833a649814e309366"
    options["cookie"] = cookies
	answer := CreateNewRequest(data,url,options)
	fmt.Println(answer["status"].(int))
	fmt.Println(string(crypt.AesDecryptECB(answer["body"].([]byte))))
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
		"body": "string",
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
			return "6.1.1"
		}()
		header["versioncode"] = func() string {
			val, ok := cookie["versioncode"]
			if ok {
				return fmt.Sprintf("%v", val)
			}
			return "140"
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
		fmt.Println(strings.TrimRight(cookies, "; "))
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
