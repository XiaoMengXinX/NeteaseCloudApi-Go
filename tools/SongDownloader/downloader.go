package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	SongDownloader "github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils"
	utils "github.com/XiaoMengXinX/NeteaseCloudApi-Go/utils"
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
	if !utils.FileExists("./config.ini") {
		var input string
		log.Println("配置文件不存在，是否重新配置token？[y/n]")
		fmt.Scanln(&input)
		if input == "Y" || input == "y" {
			log.Println("请输入你的 MUSIC_U 并回车")
			fmt.Scanln(&input)
			if input != "" {
				file, err := os.OpenFile("./config.ini", os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal("文件创建失败", err)
				}
				writeConfig := func() {
					defer file.Close()
					write := bufio.NewWriter(file)
					write.WriteString("MUSIC_U=" + input)
					write.Flush()
				}
				writeConfig()
			}
		}
	}
	config, err := utils.ReadConfig("./config.ini")
	if err != nil {
		log.Fatal(err)
	}
	cookies["MUSIC_U"] = config["MUSIC_U"]
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
	var threads = flag.Int("t", 4, "下载线程数")
	var concurrency = flag.Int("c", 1, "并发量")
	var proxy = flag.String("proxy", "", "代理设置")
	//var encrypted = flag.String("enc", "", "Only for debug")

	flag.Parse()

	if *proxy != "" {
		os.Setenv("HTTP_PROXY", *proxy)
		//os.Setenv("HTTPS_PROXY", *proxy)
		options["proxy"] = *proxy
		options["disable_https"] = true
	}
	if *threads != 4 {
		//os.Setenv("HTTPS_PROXY", *proxy)
		options["threads"] = *threads
	}
	if *concurrency != 1 {
		options["s"] = *concurrency
	}
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
	}
	if *musicid != "" {
		var ids []string
		ids = append(ids, *musicid)
		SongDownloader.DownloadSongWithMetadata(ids, options)
	}
	if *playlistid != "" {
		if *playlistoffset != 0 {
			*playlistoffset = *playlistoffset - 1
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, *playlistoffset, options)
		} else {
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, 0, options)
		}
	}
	//if *encrypted != "" {
	//	fmt.Println(SongDownloader.Decrypt163key(*encrypted))
	//}
}
