package downloader

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

type Resource struct {
	Filename string
	Url      string
}

type Downloader struct {
	wg         *sync.WaitGroup
	pool       chan *Resource
	Concurrent int
	HttpClient http.Client
	TargetDir  string
	Resources  []Resource
}

func NewDownloader(targetDir string) *Downloader {
	concurrent := runtime.NumCPU()
	return &Downloader{
		wg:         &sync.WaitGroup{},
		TargetDir:  targetDir,
		Concurrent: concurrent,
	}
}

func (d *Downloader) AppendResource(filename, url string) {
	d.Resources = append(d.Resources, Resource{
		Filename: filename,
		Url:      url,
	})
}

func (d *Downloader) Download(resource Resource, progress *mpb.Progress) error {
	defer d.wg.Done()
	d.pool <- &resource
	finalPath := d.TargetDir + "/" + resource.Filename
	// 创建临时文件
	target, err := os.Create(finalPath + ".tmp")
	if err != nil {
		return err
	}

	// 开始下载
	req, err := http.NewRequest(http.MethodGet, resource.Url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		target.Close()
		return err
	}
	defer resp.Body.Close()
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	// 创建一个进度条
	bar := progress.AddBar(
		int64(fileSize),
		// 进度条前的修饰
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"), // 已下载数量
			decor.Percentage(decor.WCSyncSpace),     // 进度百分比
		),
		// 进度条后的修饰
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" | "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()
	// 将下载的文件流拷贝到临时文件
	if _, err := io.Copy(target, reader); err != nil {
		target.Close()
		return err
	}

	// 关闭临时并修改临时文件为最终文件
	target.Close()
	if err := os.Rename(finalPath+".tmp", finalPath); err != nil {
		return err
	}
	<-d.pool
	return nil
}

func (d *Downloader) Start() error {
	d.pool = make(chan *Resource, d.Concurrent)
	if len(d.Resources) == 1 {
		log.Printf("开始下载: %v，线程数：%v", d.Resources[0].Filename, d.Concurrent)
	} else {
		var Filename string
		for i := 0; i < len(d.Resources); i++ {
			Filename = Filename + " , " + d.Resources[i].Filename
		}
		log.Printf("开始下载: %v，线程数：%v", Filename, d.Concurrent)
	}
	p := mpb.New(mpb.WithWaitGroup(d.wg))
	for _, resource := range d.Resources {
		d.wg.Add(1)
		go d.Download(resource, p)
	}
	p.Wait()
	d.wg.Wait()
	return nil
}
