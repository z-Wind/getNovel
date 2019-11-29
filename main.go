package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/noveler"
)

var (
	buildstamp = ""
	githash    = ""
	goversion  = ""

	urlNovel string
	version  bool
)

func init() {
	flag.StringVar(&urlNovel, "url", "", "小說網址")
	flag.BoolVar(&version, "version", false, "程式版本")
}

func main() {
	flag.Parse()

	switch {
	case version:
		fmt.Printf("Git Commit Hash: %s\n", githash)
		fmt.Printf("Build Time : %s\n", buildstamp)
		fmt.Printf("Golang Version : %s\n", goversion)
	case urlNovel != "":
		novel, err := chooseNoveler(urlNovel)
		if err != nil {
			log.Fatal(err)
		}

		if err := getNovel(novel); err != nil {
			fmt.Printf("m: %s\ntype: %s", err, errors.Cause(err))
			fmt.Println("\n===========stack================")
			fmt.Printf("%+v", err)
			fmt.Println("\n================================")
		}
	default:
		flag.PrintDefaults()
	}

	// 偵測是否有未關閉的 goroutine
	// 可能會因對方未斷線，而導致 goroutine 未關閉
	// time.Sleep(time.Second * 10)
	// debug.SetTraceback("all")
	// panic(1)
}

// chooseNoveler 選擇合適的 noveler
func chooseNoveler(URLNovel string) (Noveler, error) {
	u, err := url.Parse(URLNovel)
	if err != nil {
		return nil, errors.Wrap(err, "url.Parse")
	}

	switch u.Host {
	case "www.wanbentxt.com": // 完本神站
		return &noveler.WanbentxtNoveler{URL: URLNovel}, nil
	case "czbooks.net": // 小說狂人
		return &noveler.CzbooksNoveler{URL: URLNovel}, nil
	default:
		return nil, fmt.Errorf("%s No useful interface", u.Host)
	}
}

// getNovel 取得小說內容
func getNovel(novel Noveler) error {
	tmpPath := "temp"
	resultPath := "finish"

	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		os.MkdirAll(resultPath, os.ModePerm)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := ConcurrentEngine{
		Scheduler:   &QueueScheduler{ctx: ctx},
		WorkerCount: 10,
		ctx:         ctx,
	}

	// 取得章節網址
	novelPages, err := novel.GetChapterURLs()
	if err != nil {
		fmt.Printf("novel.GetChapterURLs Fail: %s\n", err)
		return errors.Wrap(err, "novel.GetChapterURLs")
	}

	var requests []Request
	for i := range novelPages {
		fileName := fmt.Sprintf("temp/%d.txt", novelPages[i].Order)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			requests = append(requests, Request{Order: novelPages[i].Order, URL: novelPages[i].URL})
		}
	}

	// 網址傳進 engine 抓取 HTML，並將小說內容存檔
	dataChan := e.Run(requests...)
	for e.numTasks != 0 {
		data := <-dataChan
		fileName := fmt.Sprintf("%d.txt", data.Order)
		filePath := path.Join(tmpPath, fileName)

		text, err := novel.GetText(data.HTML)
		if err != nil {
			fmt.Printf("novel.GetText Fail: %s\n", err)
		}

		err = ioutil.WriteFile(filePath, []byte(text), 0666)
		if err != nil {
			fmt.Printf("ioutil.WriteFile Fail: %s\n", err)
			return errors.Wrap(err, "ioutil.WriteFile")
		}
	}

	// 合併暫存檔
	if err := novel.MergeContent(tmpPath, resultPath); err != nil {
		return err
	}

	// 移除暫存檔
	if err := os.RemoveAll(tmpPath); err != nil {
		return err
	}

	return nil
}
