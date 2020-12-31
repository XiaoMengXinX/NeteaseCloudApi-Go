package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"io"
    "net/http"
	"os"
	"../utils"
	"github.com/bogem/id3v2"
)

func mainm() {
	tag, err := id3v2.Open("file.flac", id3v2.Options{Parse: true})
	if err != nil {
 		log.Fatal("Error while opening mp3 file: ", err)
 	}
	defer tag.Close()

	// Read tags.
	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

	// Set tags.
	tag.SetArtist("作曲家")
	tag.SetTitle("题目")

	artwork, err := ioutil.ReadFile("artwork.jpg")
	if err != nil {
		log.Fatal("Error while reading artwork file", err)
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
		Description: "My opinion",
		Text:        "I like this song!",
	}
	tag.AddCommentFrame(comment)

	// Write tag to file.mp3.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

func main() {
	id := os.Args[1]
	var options, cookies map[string]interface{}
	options = make(map[string]interface{})
	cookies = make(map[string]interface{})
	cookies["MUSIC_U"] = "007FAE83E56DA25DEA67D9A43C259E0DD55472DCBCA3C23B437FE930C561430891232A84D785314719AD0B0B936E24166705E05F091742F24052CCFE4937AFB3C9C1F2542B6D67F58C93D0CBEF45AF3FF1B2DAEC833D5A33FA576AA25616726FC3C1E02E9E512240844F13C60C5B3C2EF77182821D294CCFCC7359065F1DEA9D03F7DD3657CCD0950782EBE415DC51C88B72D61C898DF9E3B4951BCD941949E177B6842EFB795075CC8ED43678E06B6E219D0992F75258A8F7B86B088B030EEAEF320E6A7223925B7AAEF112CCFCF25123898935AAEBB8107669CBC35DF54A2F5BC26353E0A238CA3A8BA1FC29A82326041BCB5542367768B567912DBE1AF0C8EDAAC3FA94BB74BE94060C51FDF572B3940E4DC23D32009869B892E45E67D9827CC48AF8EBA7C74C16012CF836181BBF83FA2600F3FCE31C954223DCEBC9FA410A81D7B7A16810C06F6643E3CD519FE082FE742A52FCC87E1902899CD39B77FAC86C54BC7DE8637AD6841A7FDEC5F23DF4CB646A165D47D6ED2AF551D18D834E60"
	options["cookie"] = cookies
	result := utils.GetSongDetail(id, options)

	a := 0
	artist := ParseArtist(id,a,result)
	name := ParseName(id,a,result)
	fmt.Println(name,artist)
	//DownloadPic(id,result)
}

func DownloadPic(id string, result map[string]interface{}) {
	picpath := "./pic/"
	picname := id + ".jpg"
	picurl := fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[0].(map[string]interface{})["al"].(map[string]interface{})["picUrl"])
	resp, err := http.Get(picurl)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    out, err := os.Create(picpath+picname)
    if err != nil {
        panic(err)
    }
    defer out.Close()
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        panic(err)
    }
}

func ParseArtist(id string, s int,result map[string]interface{}) (artist string) {
	if _, ok := result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{}); ok {
		if len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})) > 0 {
		artist = fmt.Sprintf("%v", result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[0].(map[string]interface{})["name"])
		for i:=1; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})); i++ {
			artist = fmt.Sprintf("%v/%v", artist, result["body"].(map[string]interface{})["songs"].([]interface{})[s].(map[string]interface{})["ar"].([]interface{})[i].(map[string]interface{})["name"])
			}
		}
	}
	return artist
}

func ParseName(id string, i int, result map[string]interface{}) (name string) {
	name = fmt.Sprintf("%s", result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["name"])
	return name
}