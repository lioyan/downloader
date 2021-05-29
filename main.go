package main

import (
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

type Resource struct {
	FileName string
	Url      string
}
type Downloader struct {
	wg         *sync.WaitGroup
	pool       chan *Resource
	Concurrent int
	TargetDir  string
	HttpClient http.Client
	Resources  []Resource
}

func NewDownloader(targetDir string) *Downloader {
	concurrent := runtime.NumCPU()
	return &Downloader{
		TargetDir:  targetDir,
		Concurrent: concurrent,
		wg:         &sync.WaitGroup{},
	}
}
func (d *Downloader) AppendResource(filename, url string) {
	d.Resources = append(d.Resources, Resource{
		FileName: filename,
		Url:      url,
	})
}
func (d *Downloader) Download(resource Resource, progress *mpb.Progress) error {
	defer d.wg.Done()
	d.pool <- &resource
	filePath := d.TargetDir + "/" + resource.FileName
	// 创建临时文件
	if err:= os.MkdirAll(filepath.Dir(filePath), 0770); err != nil { // 创建目录
		fmt.Println(err)
		return err
	}
	target, err := os.Create(filePath + ".tmp")
	if err != nil {
		fmt.Println(err)
		return err
	}
	//开始下载
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
	// 创建进度条
	bar := progress.AddBar(
		int64(fileSize),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"), // 已下载数量
			decor.Percentage(decor.WCSyncSpace)), // 进度百分比
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60)),
	)
	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()
	// 将下载的文件拷贝到临时文件
	if _, err := io.Copy(target, reader); err != nil {
		target.Close()
		return err
	}
	// 关闭临时文件，并修改临时文件为最终文件
	target.Close()
	if err := os.Rename(filePath+".tmp", filePath); err != nil {
		return err
	}
	<-d.pool
	return nil
}
func (d *Downloader) Start() error {
	d.pool = make(chan *Resource, d.Concurrent)
	fmt.Println("开始下载， 当前并发： ", d.Concurrent)
	p := mpb.New(mpb.WithWaitGroup(d.wg))
	for _, resource := range d.Resources {
		d.wg.Add(1)
		go d.Download(resource, p)
	}
	p.Wait()
	d.wg.Wait()
	return nil
}
func main() {
	var targetDir string
	fmt.Println("请输入下载位置，并按回车结束：")
	fmt.Scanln(&targetDir)
	downloader := NewDownloader(targetDir)
	downloader.AppendResource("download1/1.pdf", "https://test-xjy-file.obs.cn-east-2.myhuaweicloud.com/202004%2Faa18f6a5-644e-4519-804b-5dcd82cf8828.bmp?response-content-disposition=attachment%3Bfilename%3D20.bmp%3Bfilename%2A%3Dutf-8%27%2720.bmp&response-content-type=binary%2Foctet-stream&AWSAccessKeyId=R7WWGCCN8NPIWHPLFOLT&Expires=1653797183&Signature=xqEEZkNAOjwvoTz5MO7URuYy1iQ%3D")
	downloader.AppendResource("download2/2.pdf", "https://test-xjy-file.obs.cn-east-2.myhuaweicloud.com/202004%2Faa18f6a5-644e-4519-804b-5dcd82cf8828.bmp?response-content-disposition=attachment%3Bfilename%3D20.bmp%3Bfilename%2A%3Dutf-8%27%2720.bmp&response-content-type=binary%2Foctet-stream&AWSAccessKeyId=R7WWGCCN8NPIWHPLFOLT&Expires=1653797183&Signature=xqEEZkNAOjwvoTz5MO7URuYy1iQ%3D")
	err := downloader.Start()
	if err != nil {
		panic(err)
	}
}

