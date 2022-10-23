package tool

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func SingleDownloadFile(file_link, file_path string) error {
	res, err := http.Get(file_link)
	if err != nil {
		return err
	}
	f, err := os.Create(file_path)
	if err != nil {
		return err
	}
	io.Copy(f, res.Body)
	return nil
}

func GetFileContentType(out *bufio.Reader) (string, string, error) {
	out = bufio.NewReader(out)
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	if contentType == "" {
		contentType = "none"
	}
	return contentType, strings.Split(contentType, "/")[0], nil
}

func Download_file(link, file_path string) {
	// go run . https://az764295.vo.msecnd.net/stable/74b1f979648cc44d385a2286793c226e611f59e7/VSCodeUserSetup-x64-1.71.2.exe
	concurrencyN := runtime.NumCPU() // 默认并发数
	fmt.Println("默认并发数" + strconv.Itoa(concurrencyN))
	NewDownloader(concurrencyN).Download(link, file_path)
}

type Downloader struct {
	concurrency int
}

func NewDownloader(concurrency int) *Downloader {
	return &Downloader{concurrency: concurrency}
}

func (d *Downloader) Download(strURL, filename string) error {
	if filename == "" {
		filename = path.Base(strURL)
	}
	resp, err := http.Head(strURL)
	fmt.Print("文件大小为")
	fmt.Print(float32(resp.ContentLength) / 1024 / 1024)
	fmt.Println("Mb")
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.multiDownload(strURL, filename, int(resp.ContentLength))
	} else {
		fmt.Println("Download不能并发下载,只能singleDownload")
		return SingleDownloadFile(strURL, filename)
	}
}

func (d *Downloader) downloadPartial(strURL, filename string, rangeStart, rangeEnd, i int) {
	if rangeStart >= rangeEnd {
		return
	}
	part_file_name := d.getPartFilename(filename, i)
	fmt.Printf("%s start download\n", part_file_name)
	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(part_file_name, flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(partFile, resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("%s download success\n", part_file_name)
}

// getPartDir 部分文件存放的目录
func (d *Downloader) getFileName(filename string) string {
	arr := strings.Split(filename, "/")
	return arr[len(arr)-1]
}

// getPartDir 部分文件存放的目录
func (d *Downloader) getPartDir(filename string) string {
	return strings.SplitN(filename, ".", 2)[0]
}

// getPartFilename 构造部分文件的名字
func (d *Downloader) getPartFilename(filename string, partNum int) string {
	file_name := d.getFileName(filename)
	return fmt.Sprintf("%s-slice/%s-%d", filename, file_name, partNum)
}

func (d *Downloader) multiDownload(strURL, filename string, contentLen int) error {
	partSize := contentLen / d.concurrency
	// 创建部分文件的存放目录
	partDir := d.getPartDir(filename + "-slice")
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.concurrency)

	rangeStart := 0

	for i := 0; i < d.concurrency; i++ {
		// 并发请求
		go func(i, rangeStart int) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if i == d.concurrency-1 {
				rangeEnd = contentLen
			}

			d.downloadPartial(strURL, filename, rangeStart, rangeEnd, i)

		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	wg.Wait()

	// 合并文件
	d.merge(filename)
	fmt.Printf("%s download and merge success\n", filename)
	return nil
}

func (d *Downloader) merge(filename string) error {
	destFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.concurrency; i++ {
		partFileName := d.getPartFilename(filename, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		fmt.Printf("[%d] %s merge to %s success\n", i, partFileName, filename)
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}
