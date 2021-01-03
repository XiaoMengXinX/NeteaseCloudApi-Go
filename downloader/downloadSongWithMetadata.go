package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"./tools"
	"../utils"
	"../utils/crypt"
	"github.com/bogem/id3v2"
	"github.com/goulash/audio/flac"
	"github.com/tcolgate/mp3"
	"github.com/yoki123/ncmdump/tag"
)

const (
	picPath       = "./pic/"
	musicPath     = "./music/"
	fileNameStyle = 1
)

func main() {
	var options, cookies map[string]interface{}
	options = make(map[string]interface{})
	cookies = make(map[string]interface{})
	cookies["MUSIC_U"] = "007FAE83E56DA25DEA67D9A43C259E0DD55472DCBCA3C23B437FE930C561430891232A84D785314719AD0B0B936E24166705E05F091742F24052CCFE4937AFB3C9C1F2542B6D67F58C93D0CBEF45AF3FF1B2DAEC833D5A33FA576AA25616726FC3C1E02E9E512240844F13C60C5B3C2EF77182821D294CCFCC7359065F1DEA9D03F7DD3657CCD0950782EBE415DC51C88B72D61C898DF9E3B4951BCD941949E177B6842EFB795075CC8ED43678E06B6E219D0992F75258A8F7B86B088B030EEAEF320E6A7223925B7AAEF112CCFCF25123898935AAEBB8107669CBC35DF54A2F5BC26353E0A238CA3A8BA1FC29A82326041BCB5542367768B567912DBE1AF0C8EDAAC3FA94BB74BE94060C51FDF572B3940E4DC23D32009869B892E45E67D9827CC48AF8EBA7C74C16012CF836181BBF83FA2600F3FCE31C954223DCEBC9FA410A81D7B7A16810C06F6643E3CD519FE082FE742A52FCC87E1902899CD39B77FAC86C54BC7DE8637AD6841A7FDEC5F23DF4CB646A165D47D6ED2AF551D18D834E60"
	options["cookie"] = cookies
	//options["s"] = 5
	options["savePath"] = musicPath
	CheckPathExists(picPath)
	CheckPathExists(musicPath)

	var musicid = flag.String("m", "", "歌曲id")
	var playlistid = flag.String("p", "", "歌单id")
	var playlistoffset = flag.Int("s", 0, "歌单偏移量")
	var encrypted = flag.String("enc", "", "Only for test")

	flag.Parse()
	if *musicid != "" {
		DownloadSongWithMetadata(*musicid, options)
	}
	if *playlistid != "" {
		if *playlistoffset != 0 {
			DownloadPLaylistWithMetadata(*playlistid, *playlistoffset, options)
		} else {
			DownloadPLaylistWithMetadata(*playlistid, 0, options)
		}
	}
	if *encrypted != "" {
		fmt.Println(Decrypt163key(*encrypted))
	}
}

func DownloadSongWithMetadata(id string, options map[string]interface{}) {
	result := utils.GetSongDetail(id, options)
	fileName := utils.DownloadSong(id, options)

	if len(result["body"].(map[string]interface{})["songs"].([]interface{})) > 0 {
		for i := 0; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})); i++ {
			artist, artistMap := ParseArtist(id, i, result)
			name := ParseName(id, i, result)
			album, albumId, albumPic, albumPicDocId := ParseAlbum(id, i, result)
			//fmt.Println(artistMap)
			//fmt.Println(name, artist, album)
			filename := fileName[i]
			if filename == "null" {
				continue
			}
			musicMarker := MusicMarker(id, filename, name, album, albumId, albumPic, albumPicDocId, i, options, result, artistMap)
			//fmt.Println(marker)
			picName := DownloadPic(fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["id"].(float64))), i, result)

			format := strings.Replace(path.Ext(filename), ".", "", -1)
			switch format {
			case "mp3":
				AddMp3Id3v2(filename, name, artist, album, picName, musicMarker)
			case "flac":
				AddFlacId3v2(filename, name, artist, album, picName, musicMarker)
			}

			//var replacer = strings.NewReplacer("/", " ")
			//sysType := runtime.GOOS
			//if sysType == "windows" {
			//	replacer = strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ", "|", " ", "\\", " ", "<", " ", ">", " ")
			//}
			var replacer = strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ", "|", " ", "\\", " ", "<", " ", ">", " ")

			var newFilename string
			switch fileNameStyle {
			case 1:
				newFilename = replacer.Replace(fmt.Sprintf("%v - %v%v", strings.Replace(artist, "/", ",", -1), name, path.Ext(filename)))
			case 2:
				newFilename = replacer.Replace(fmt.Sprintf("%v - %v%v", name, strings.Replace(artist, "/", ",", -1), path.Ext(filename)))
			case 3:
				newFilename = replacer.Replace(fmt.Sprintf("%v%v", strings.Replace(name, "/", " ", -1), path.Ext(filename)))
			}
			err := os.Rename(musicPath+filename, musicPath+newFilename)
			fmt.Println(newFilename + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
}

func DownloadPLaylistWithMetadata(id string, offset int, options map[string]interface{}) {
	result := utils.GetPlaylistDetail(id, options)
	if _, ok := result["body"].(map[string]interface{})["playlist"]; ok {
		if _, ok := options["s"].(int); ok {
			var mid string
			var i int = 0
			for t, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
				if t > offset {
					if i < options["s"].(int) {
						if _, ok := v.(map[string]interface{})["id"].(float64); ok {
							if i == 0 {
								mid = fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64)))
							} else {
								mid = fmt.Sprintf("%v,%v", mid, int(v.(map[string]interface{})["id"].(float64)))
							}
							if i == options["s"].(int)-1 {
								i = 0
								DownloadSongWithMetadata(mid, options)
							} else {
								if len(result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}))-t == 1 {
									DownloadSongWithMetadata(mid, options)
								} else {
									i++
								}
							}
						}
					}
				}
			}
		} else {
			for t, v := range result["body"].(map[string]interface{})["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
				if t > offset {
					var mid string
					mid = fmt.Sprintf("%v", int(v.(map[string]interface{})["id"].(float64)))
					//fmt.Println(mid)
					DownloadSongWithMetadata(mid, options)
				}
			}
		}
	}
}

func DownloadPic(id string, i int, result map[string]interface{}) (picName string) {
	picName = id + ".jpg"
	picurl := fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["al"].(map[string]interface{})["picUrl"])
	resp, err := http.Get(picurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(picPath + picName)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	return picName
}

func MusicMarker(id, filename, name, album, albumId, albumPic, albumPicDocId string, s int, options, result, artistMap map[string]interface{}) (marker string) {
	var data map[string]interface{}
	data = make(map[string]interface{})
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
	if _, ok := result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["mv"].(float64) ; ok {
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
		Flag		  int			`json:"flag"`
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

func ParseArtist(id string, s int, result map[string]interface{}) (artist string, artistMap map[string]interface{}) {
	if _, ok := result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{}); ok {
		if len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})) > 0 {
			artist = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["name"])
			id = fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["id"].(float64)))
			var ar string
			ar = fmt.Sprintf("[\"%v\",%v]", artist, id)
			for i := 1; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})); i++ {
				artist = fmt.Sprintf("%v/%v", artist, result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["name"])
				id = fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["id"].(float64)))
				var Artist string
				Artist = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["name"])
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
	r, err := os.Open(options["savePath"].(string) + filename)
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			return
		}
		t = t + f.Duration().Seconds()
	}
	bitRate = int(f.Header().BitRate())
	duration = int(math.Floor(t * 1000))
	err = r.Close()
	if err != nil {
		fmt.Println(err)
	}

	return bitRate, duration
}

func GetFlacInfo(filename string, options map[string]interface{}) (bitRate, duration int) {
	file, _ := os.Stat(options["savePath"].(string) + filename)
	data, _ := flac.ReadFileMetadata(options["savePath"].(string) + filename)
	length := data.Length()
	duration = int(length / 1000000)
	bitRate = (int(file.Size()) * 8) / (duration / 1000)

	return bitRate, duration
}

func AddMp3Id3v2(filename, name, artist, album, picName, MusicMarker string) {
	tag, _ := id3v2.Open(musicPath+filename, id3v2.Options{Parse: false})
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(name)
	tag.SetArtist(artist)

	if album != "" {
		tag.SetAlbum(album)
	}

	artwork, err := ioutil.ReadFile(picPath + picName)
	if err != nil {
		log.Fatal("Error while reading AlbumPic", err)
	}
	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     artwork,
	}
	tag.AddAttachedPicture(pic)

	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "",
		Text:        MusicMarker,
	}
	tag.AddCommentFrame(comment)

	if err = tag.Save(); err != nil {
		log.Fatal("Error: ", err)
	}
}

func AddFlacId3v2(filename, name, artist, album, picName, MusicMarker string) {
	tag, err := tag.NewFlacTagger(musicPath+filename)
	if err != nil {
		log.Fatal(err)
	}
	tag.SetTitle(name)
	
	artists := make([]string, 0)
	artists = append(artists, artist)
	tag.SetArtist(artists)

	if album != "" {
		tag.SetAlbum(album)
	}

	artwork, err := ioutil.ReadFile(picPath + picName)
	fileType := tools.GetFileType(artwork)
	var mime string = "image/png"
	switch fileType {
	case "jpg":
		mime = "image/jpeg"
	case "png":
		mime = "image/png"
	}
	err = tag.SetCover(artwork, mime)
	if err != nil {
		log.Fatal(err)
	}

	tag.SetComment(MusicMarker)
	err = tag.Save()
	if err != nil {
		log.Fatal("Error: ", err)
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
			fmt.Printf("mkdir %v failed: %v\n", path, err)
		}
		return false
	}
	fmt.Printf("Error: %v\n", err)
	return false
}

func Decrypt163key(encrypted string) (decrypted string) {
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	return string(crypt.MarkerAesDecryptECB(data))
}