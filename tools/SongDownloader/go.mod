module github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader

go 1.15

require (
	github.com/XiaoMengXinX/NeteaseCloudApi-Go v0.0.0
	github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils v0.0.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
)

replace github.com/XiaoMengXinX/NeteaseCloudApi-Go => ../../

replace github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader/utils => ./utils
