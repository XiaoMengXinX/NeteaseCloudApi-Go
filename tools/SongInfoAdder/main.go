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
	picPath   = "./pic/"
	musicPath = "./temp/"
)

func main() {
	SongDownloader.CheckPathExists(picPath)
	SongDownloader.CheckPathExists(musicPath)
	var filename string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "music file name",
				Destination: &filename,
			},
		},
		Action: func(c *cli.Context) error {
			var id string = strings.Replace(filepath.Base(path.Base(filename)), path.Ext(filename), "", -1)
			SongInfoAdder(id, filename)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func SongInfoAdder(id, filename string) {
	var options map[string]interface{}
	options = make(map[string]interface{})
	options["savePath"] = "./"
	result := utils.GetSongDetail(id, options)

	if len(result["body"].(map[string]interface{})["songs"].([]interface{})) > 0 {
		for i := 0; i < len(result["body"].(map[string]interface{})["songs"].([]interface{})); i++ {
			var i int = 0
			artist, artistMap := SongDownloader.ParseArtist(id, i, result)
			name := SongDownloader.ParseName(id, i, result)
			album, albumId, albumPic, albumPicDocId := SongDownloader.ParseAlbum(id, i, result)
			musicMarker := SongDownloader.MusicMarker(id, filename, name, album, albumId, albumPic, albumPicDocId, i, options, result, artistMap)
			//fmt.Println(marker)
			picName := SongDownloader.DownloadPic(fmt.Sprintf("%v", int(result["body"].(map[string]interface{})["songs"].([]interface{})[i].(map[string]interface{})["id"].(float64))), i, result)

			format := strings.Replace(path.Ext(filename), ".", "", -1)
			switch format {
			case "mp3":
				SongDownloader.AddMp3Id3v2(filename, name, artist, album, picName, musicMarker)
			case "flac":
				SongDownloader.AddFlacId3v2(filename, name, artist, album, picName, musicMarker)
			}
		}
	}
}
