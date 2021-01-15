package main

import (
	"flag"
	"fmt"

	SongDownloader "github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils"
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
	options["picPath"] = picPath
	options["fileNameStyle"] = fileNameStyle

	SongDownloader.CheckPathExists(picPath)
	SongDownloader.CheckPathExists(musicPath)

	var musicid = flag.String("m", "", "歌曲id")
	var playlistid = flag.String("p", "", "歌单id")
	var playlistoffset = flag.Int("s", 0, "歌单偏移量")
	var encrypted = flag.String("enc", "", "Only for test")

	flag.Parse()
	if *musicid != "" {
		SongDownloader.DownloadSongWithMetadata(*musicid, options)
	}
	if *playlistid != "" {
		if *playlistoffset != 0 {
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, *playlistoffset, options)
		} else {
			SongDownloader.DownloadPLaylistWithMetadata(*playlistid, 0, options)
		}
	}
	if *encrypted != "" {
		fmt.Println(SongDownloader.Decrypt163key(*encrypted))
	}
}
