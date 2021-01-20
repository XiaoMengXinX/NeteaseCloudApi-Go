package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	SongDownloader "github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils"
	"github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils"
	"github.com/urfave/cli/v2"
)

const (
	picPath = "./pic/"
)

func main() {
	SongDownloader.CheckPathExists(picPath)
	var file string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "music file name",
				Destination: &file,
			},
		},
		Action: func(c *cli.Context) error {
			var fileName, filePath string
			var id string = strings.Replace(filepath.Base(path.Base(file)), path.Ext(file), "", -1)
			fileName = filepath.Base(file)
			filePath = filepath.Dir(file)
			SongInfoAdder(id, fileName, filePath)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func SongInfoAdder(id, fileName, filePath string) {
	var options map[string]interface{}
	options = make(map[string]interface{})
	options["savePath"] = filePath
	options["picPath"] = picPath
	result := utils.GetSongDetail(id, options)

	if len(result["body"].(map[string]interface{})["songs"].([]interface{})) > 0 {
		for i := 0; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})); i++ {
			var i int = 0
			artist, artistMap := SongDownloader.ParseArtist(id, i, result)
			name := SongDownloader.ParseName(id, i, result)
			album, albumId, albumPic, albumPicDocId := SongDownloader.ParseAlbum(id, i, result)
			musicMarker := SongDownloader.MusicMarker(id, fileName, name, album, albumId, albumPic, albumPicDocId, i, options, result, artistMap)
			//fmt.Println(marker)
			picName := SongDownloader.DownloadPic(fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["id"].(float64))), i, result, options)

			format := strings.Replace(path.Ext(fileName), ".", "", -1)
			switch format {
			case "mp3":
				SongDownloader.AddMp3Id3v2(fileName, name, artist, album, picName, musicMarker, options)
			case "flac":
				SongDownloader.AddFlacId3v2(fileName, name, artist, album, picName, musicMarker, options)
			}
		}
	}
}
