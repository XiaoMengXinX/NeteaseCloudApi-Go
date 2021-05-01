package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	SongDownloader "github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils"
	log "github.com/sirupsen/logrus"
)

const (
	picPath       = "./pic"
	musicPath     = "./music"
	fileNameStyle = 1
)

//Custom log format definition
type LogFormatter struct{}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func init() {
	timeNow := time.Now()
	timeStamp := timeNow.Unix()
	logFile := fmt.Sprintf("log-%v.log", timeStamp)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
	}
	output := io.MultiWriter(os.Stdout, file)
	log.SetOutput(output)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetFormatter(new(LogFormatter))
	log.SetLevel(log.InfoLevel)
}

func main() {
	var options, cookies map[string]interface{}
	options = make(map[string]interface{})
	cookies = make(map[string]interface{})
	cookies["MUSIC_U"] = "aa31a213cb2dff5ba39ff7623ef3308ba14dd75bdf4434e22f9ccddbcf6aa43f33a649814e309366"
	options["cookie"] = cookies
	//options["s"] = 5

	options["savePath"] = musicPath
	options["picPath"] = picPath
	options["fileNameStyle"] = fileNameStyle

	SongDownloader.CheckPathExists(picPath)
	SongDownloader.CheckPathExists(musicPath)

	var musicid = flag.String("m", "", "歌曲id")
	var playlistid = flag.String("p", "", "歌单id")
	var playlistoffset = flag.Int("s", 0, "歌单偏移量")
	var loglevel = flag.Int("l", 4, "日志等级")
	var encrypted = flag.String("enc", "", "Only for debug")

	*playlistid = "6736706492"
	flag.Parse()
	if *loglevel != 4 {
		switch {
		case *loglevel == 0:
			log.SetLevel(log.PanicLevel)
		case *loglevel == 1:
			log.SetLevel(log.FatalLevel)
		case *loglevel == 2:
			log.SetLevel(log.ErrorLevel)
		case *loglevel == 3:
			log.SetLevel(log.WarnLevel)
		case *loglevel == 4:
			log.SetLevel(log.InfoLevel)
		case *loglevel == 5:
			log.SetLevel(log.DebugLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
		fmt.Println(SongDownloader.Decrypt163key(*encrypted))
	}
	if *musicid != "" {
		SongDownloader.DownloadSongWithMetadata(*musicid, options)
	}
	if *playlistid != "" {
		if *playlistoffset != 0 {
			*playlistoffset = *playlistoffset - 1
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, *playlistoffset, options)
		} else {
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, 0, options)
		}
	}
	if *encrypted != "" {
		fmt.Println(SongDownloader.Decrypt163key(*encrypted))
	}

}
