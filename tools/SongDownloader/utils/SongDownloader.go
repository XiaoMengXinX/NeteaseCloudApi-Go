package songdownloader

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils"
	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils/crypt"
	"github.com/bogem/id3v2"
	"github.com/goulash/audio/flac"
	log "github.com/sirupsen/logrus"
	"github.com/tcolgate/mp3"
	"github.com/yoki123/ncmdump/tag"
)

var musicPath, picPath string = "./pic", "./music"
var fileNameStyle int = 1

func DownloadSongWithMetadata(ids []string, resultCache, options map[string]interface{}) error {
	startTime := time.Now()
	var fileName, validIds []string
	if len(ids) == 1 && resultCache["SongUrl"] != nil {
		if _, ok := resultCache["SongUrl"].(map[string]interface{}); ok {
			if resultCache["SongUrl"].(map[string]interface{})["body"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["url"] != nil {
				fileName = utils.DownloadSong(ids[0], resultCache, options)
				validIds = []string{ids[0]}
			} else {
				return fmt.Errorf("获取 musicid : %s 下载链接失败", ids[0])
			}
		} else {
			return fmt.Errorf("获取 musicid : %s 下载链接失败", ids[0])
		}
	} else {
		fileName, validIds = utils.MultiDownloadSong(ids, options)
	}

	for m := 0; m < len(validIds); m++ {
		id := string(validIds[m])
		var result map[string]interface{}
		if len(ids) == 1 && resultCache["SongDetail"] != nil {
			if _, ok := resultCache["SongDetail"].(map[string]interface{}); ok {
				result = resultCache["SongDetail"].(map[string]interface{})
			} else {
				result = utils.GetSongDetail(id, options)
			}
		} else {
			result = utils.GetSongDetail(id, options)
		}

		if _, ok := options["savePath"].(string); ok {
			musicPath = options["savePath"].(string)
		}
		if _, ok := options["picPath"].(string); ok {
			picPath = options["picPath"].(string)
		}
		if _, ok := options["fileNameStyle"].(int); ok {
			fileNameStyle = options["fileNameStyle"].(int)
		}

		if len(result["body"].(map[string]interface{})["songs"].([]interface{})) > 0 {
			for i := 0; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})); i++ {
				artist, artistMap := ParseArtist(i, result)
				name := ParseName(id, i, result)
				album, albumId, albumPic, albumPicDocId := ParseAlbum(id, i, result)
				//fmt.Println(artistMap)
				//fmt.Println(name, artist, album)
				filename := fileName[m]
				if filename == "null" {
					continue
				}
				musicMarker := MusicMarker(id, filename, name, album, albumId, albumPic, albumPicDocId, i, options, result, artistMap)
				//fmt.Println(musicMarker)
				picName := DownloadPic(fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["id"].(float64))), i, result, options)

				format := strings.Replace(path.Ext(filename), ".", "", -1)
				switch format {
				case "mp3":
					AddMp3Id3v2(filename, name, artist, album, picName, musicMarker, options)
				case "flac":
					AddFlacId3v2(filename, name, artist, album, picName, musicMarker, options)
				}
				optionsJson, _ := json.Marshal(options)
				log.Debugf("\n\tfilename: %v\n\tname: %v\n\tartist: %v\n\talbum: %v\n\tpicName: %v\n\tmusicMarker: %v\n\toptions: %v", filename, name, artist, album, picName, musicMarker, string(optionsJson))

				//var replacer = strings.NewReplacer("/", " ")
				//sysType := runtime.GOOS
				//if sysType == "windows" {
				//var replacer = strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ", "|", " ", "\\", " ", "<", " ", ">", " ")
				//}
				var replacer = strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ", "|", " ", "\\", " ", "<", " ", ">", " ", "\"", " ")

				var newFilename string
				switch fileNameStyle {
				case 1:
					newFilename = replacer.Replace(fmt.Sprintf("%v - %v%v", strings.Replace(artist, "/", ",", -1), name, path.Ext(filename)))
				case 2:
					newFilename = replacer.Replace(fmt.Sprintf("%v - %v%v", name, strings.Replace(artist, "/", ",", -1), path.Ext(filename)))
				case 3:
					newFilename = replacer.Replace(fmt.Sprintf("%v%v", strings.Replace(name, "/", " ", -1), path.Ext(filename)))
				}
				if fileNameStyle != 0 {
					err := os.Rename(musicPath+"/"+filename, musicPath+"/"+newFilename)
					log.Printf("%s 下载完成 耗时: %f second\n", fmt.Sprintf("%v - %v", artist, name), time.Since(startTime).Seconds())
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	return nil
}

func DownloadPLaylistWithMetadata(id string, offset int, options map[string]interface{}) {
	result := utils.GetPlaylistDetail(id, options)
	if _, ok := result["body"].(map[string]interface{})["playlist"]; ok {
		var ids []string
		if _, ok := options["s"].(int); ok {
			var i int = 0
			for t, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["trackIds"].([]interface{}) {
				if t >= offset {
					if i < options["s"].(int) {
						if _, ok := v.(map[string]interface{})["id"].(float64); ok {
							ids = append(ids, fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64))))
							if i == options["s"].(int)-1 {
								i = 0
								DownloadSongWithMetadata(ids, make(map[string]interface{}), options)
								ids = ids[0:0]
							} else {
								if len(result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["trackIds"].([]interface{}))-t == 1 {
									DownloadSongWithMetadata(ids, make(map[string]interface{}), options)
									ids = ids[0:0]
								} else {
									i++
								}
							}
						}
					}
				}
			}
		} else {
			for t, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["trackIds"].([]interface{}) {
				if t >= offset {
					ids = append(ids, fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64))))
					//fmt.Println(mid)
					DownloadSongWithMetadata(ids, make(map[string]interface{}), options)
					ids = ids[0:0]
				}
			}
		}
	}
}

func DownloadPic(id string, i int, result, options map[string]interface{}) (picName string) {
	client := &http.Client{}
	if _, ok := options["proxy"].(string); ok {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	picName = id + ".jpg"
	if _, ok := options["savePath"].(string); ok {
		musicPath = options["savePath"].(string)
	}
	if _, ok := options["picPath"].(string); ok {
		picPath = options["picPath"].(string)
	}
	picurl := fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["al"].(map[string]interface{})["picUrl"])
	resp, err := client.Get(picurl)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(picPath + "/" + picName)
	if err != nil {
		log.Error(err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
	}
	return picName
}

func MusicMarker(id, filename, name, album, albumId, albumPic, albumPicDocId string, s int, options, result, artistMap map[string]interface{}) (marker string) {
	var data map[string]interface{} = make(map[string]interface{})
	format := path.Ext(filename)
	data["format"] = strings.Replace(format, ".", "", -1)
	data["musicId"], _ = strconv.Atoi(id)
	data["musicName"] = name
	data["artist"] = artistMap["artist"]
	if album != "" {
		data["album"] = album
		data["albumId"], _ = strconv.Atoi(albumId)
	} else {
		data["album"] = ""
		data["albumId"] = 0
	}
	//data["albumPic"] = strings.Replace(albumPic, "/", "\\/", -1)
	data["albumPic"] = albumPic
	data["albumPicDocId"], _ = strconv.Atoi(albumPicDocId)
	if _, ok := result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["mv"].(float64); ok {
		data["mvId"] = int(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["mv"].(float64))
	} else {
		data["mvId"] = 0
	}
	data["flag"] = 0

	var bitRate, duration int
	switch data["format"].(string) {
	case "mp3":
		bitRate, duration = GetMp3Info(filename, options)
	case "flac":
		bitRate, duration = GetFlacInfo(filename, options)
	}
	data["bitRate"] = bitRate
	data["duration"] = duration

	data["alias"] = result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["alia"]
	jsonStruct := struct {
		Format        string        `json:"format"`
		MusicId       int           `json:"musicId"`
		MusicName     string        `json:"musicName"`
		Artist        []interface{} `json:"artist"`
		Album         string        `json:"album"`
		AlbumId       int           `json:"albumId"`
		AlbumPicDocId int           `json:"albumPicDocId"`
		AlbumPic      string        `json:"albumPic"`
		MvId          int           `json:"mvId"`
		Flag          int           `json:"flag"`
		Bitrate       int           `json:"bitrate"`
		Duration      int           `json:"duration"`
		Alias         []interface{} `json:"alias"`
	}{data["format"].(string), data["musicId"].(int), data["musicName"].(string), data["artist"].([]interface{}), data["album"].(string), data["albumId"].(int), data["albumPicDocId"].(int), data["albumPic"].(string), data["mvId"].(int), data["flag"].(int), data["bitRate"].(int), data["duration"].(int), data["alias"].([]interface{})}
	jsonData, _ := json.Marshal(jsonStruct)
	//jsonData := strings.Replace(string(jsondata), "/", "\\/", -1)
	marker = fmt.Sprintf("163 key(Don't modify):%v", string(base64.StdEncoding.EncodeToString(crypt.MarkerAesEncryptECB("music:"+string(jsonData)))))
	//fmt.Println(string(jsonData))
	return marker
}

func ParseArtist(s int, result map[string]interface{}) (artist string, artistMap map[string]interface{}) {
	if _, ok := result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{}); ok {
		if len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})) > 0 {
			var id string
			artist = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["name"])
			id = fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["id"].(float64)))
			var ar string
			ar = fmt.Sprintf("[\"%v\",%v]", artist, id)
			for i := 1; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})); i++ {
				artist = fmt.Sprintf("%v/%v", artist, result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["name"])
				id = fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["id"].(float64)))
				var Artist string = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["name"])
				ar = fmt.Sprintf("%v,[\"%v\",%v]", ar, Artist, id)
			}
			jsonStr := []byte(fmt.Sprintf("{\"artist\":[%v]}", ar))
			json.Unmarshal(jsonStr, &artistMap)
		}
	}
	return artist, artistMap
}

func ParseName(id string, i int, result map[string]interface{}) (name string) {
	name = fmt.Sprintf("%s", result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["name"])
	return name
}

func ParseAlbum(id string, i int, result map[string]interface{}) (album, albumId, albumPic, albumPicDocId string) {
	album = fmt.Sprintf("%s", result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["al"].(map[string]interface{})["name"])
	albumId = fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["al"].(map[string]interface{})["id"].(float64)))
	albumPic = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[0].(map[string]interface{})["al"].(map[string]interface{})["picUrl"])
	albumPicDocId = strings.Replace(filepath.Base(path.Base(albumPic)), path.Ext(albumPic), "", -1)
	return album, albumId, albumPic, albumPicDocId
}

func GetMp3Info(filename string, options map[string]interface{}) (bitRate, duration int) {
	t := 0.0
	r, err := os.Open(options["savePath"].(string) + "/" + filename)
	if err != nil {
		log.Error(err)
		return
	}
	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
			return
		}
		t = t + f.Duration().Seconds()
	}
	bitRate = int(f.Header().BitRate())
	duration = int(math.Floor(t * 1000))
	err = r.Close()
	if err != nil {
		log.Error(err)
	}

	return bitRate, duration
}

func GetFlacInfo(filename string, options map[string]interface{}) (bitRate, duration int) {
	file, _ := os.Stat(options["savePath"].(string) + "/" + filename)
	data, _ := flac.ReadFileMetadata(options["savePath"].(string) + "/" + filename)
	length := data.Length()
	duration = int(length / 1000000)
	bitRate = (int(file.Size()) * 8) / (duration / 1000)

	return bitRate, duration
}

func AddMp3Id3v2(filename, name, artist, album, picName, MusicMarker string, options map[string]interface{}) {
	if _, ok := options["savePath"].(string); ok {
		musicPath = options["savePath"].(string)
	}
	if _, ok := options["picPath"].(string); ok {
		picPath = options["picPath"].(string)
	}

	tag, _ := id3v2.Open(musicPath+"/"+filename, id3v2.Options{Parse: false})
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(name)
	tag.SetArtist(artist)

	if album != "" {
		tag.SetAlbum(album)
	}

	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "",
		Text:        MusicMarker,
	}
	tag.AddCommentFrame(comment)

	artwork, err := ioutil.ReadFile(picPath + "/" + picName)
	if err != nil {
		log.Error("Error while reading AlbumPic", err)
	}
	var mime string
	fileCode := bytesToHexString(artwork)
	if strings.HasPrefix(fileCode, "ffd8ffe000104a464946") {
		mime = "image/jpeg"
	}
	if strings.HasPrefix(fileCode, "89504e470d0a1a0a0000") {
		mime = "image/png"
	}
	if mime != "" {
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    mime,
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     artwork,
		}
		tag.AddAttachedPicture(pic)
	}

	if err = tag.Save(); err != nil {
		log.Error("Error: ", err)
	}
}

func AddFlacId3v2(filename, name, artist, album, picName, MusicMarker string, options map[string]interface{}) {
	if _, ok := options["savePath"].(string); ok {
		musicPath = options["savePath"].(string)
	}
	if _, ok := options["picPath"].(string); ok {
		picPath = options["picPath"].(string)
	}

	tag, err := tag.NewFlacTagger(musicPath + "/" + filename)
	if err != nil {
		log.Error(err)
	}
	tag.SetTitle(name)

	artists := make([]string, 0)
	artists = append(artists, artist)
	tag.SetArtist(artists)

	if album != "" {
		tag.SetAlbum(album)
	}

	artwork, err := ioutil.ReadFile(picPath + "/" + picName)
	if err != nil {
		log.Error(err)
	}

	var mime string
	fileCode := bytesToHexString(artwork)
	if strings.HasPrefix(fileCode, "ffd8ffe000104a464946") {
		mime = "image/jpeg"
	}
	if strings.HasPrefix(fileCode, "89504e470d0a1a0a0000") {
		mime = "image/png"
	}

	if mime != "" {
		err = tag.SetCover(artwork, mime)
		if err != nil {
			log.Error(err)
		}
	}

	tag.SetComment(MusicMarker)
	err = tag.Save()
	if err != nil {
		log.Error("Error: ", err)
	}
}

func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Errorf("mkdir %v failed: %v\n", path, err)
		}
		return false
	}
	log.Errorf("Error: %v\n", err)
	return false
}

func Decrypt163key(encrypted string) (decrypted string) {
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	return string(crypt.MarkerAesDecryptECB(data))
}

func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
