module github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongInfoAdder

go 1.15

require (
	github.com/XiaoMengXinX/NeteaseCloudApi-Go v0.0.0
	github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader v0.0.0
	github.com/urfave/cli/v2 v2.3.0
)

replace (
	github.com/XiaoMengXinX/NeteaseCloudApi-Go => ../../
	github.com/XiaoMengXinX/NeteaseCloudApi-Go/tools/SongDownloader => ../SongDownloader
)
