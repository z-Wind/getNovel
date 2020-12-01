package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/z-Wind/concurrencyengine"
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
	flag.StringVar(&urlNovel, "url", "", "小說目錄網址")
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
		concurrencyengine.ELog.Start("engine.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		concurrencyengine.ELog.SetFlags(0)
		defer concurrencyengine.ELog.Stop()

		novel, err := chooseNoveler(urlNovel)
		if err != nil {
			concurrencyengine.ELog.Fatalf("%v\n", err)
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
func chooseNoveler(URLNovel string) (noveler.Noveler, error) {
	u, err := url.Parse(URLNovel)
	if err != nil {
		return nil, errors.Wrap(err, "url.Parse")
	}

	switch u.Host {
	case "m.wanbentxt.com": // 完本神站 mobile
		u.Host = "www.wanbentxt.com"
		fmt.Printf("%s => %s\n", URLNovel, u.String())
		URLNovel = u.String()
		fallthrough
	case "www.wanbentxt.com": // 完本神站
		concurrencyengine.ELog.LPrintf("Noveler Choose 完本神站\n")
		return noveler.NewWanbentxtNoveler(URLNovel), nil
	case "czbooks.net": // 小說狂人
		concurrencyengine.ELog.LPrintf("Noveler Choose 小說狂人\n")
		return noveler.NewCzbooksNoveler(URLNovel), nil
	case "www.hjwzw.com": // 黃金屋 簡體
		u.Host = "tw.hjwzw.com"
		fmt.Printf("%s => %s\n", URLNovel, u.String())
		URLNovel = u.String()
		fallthrough
	case "tw.hjwzw.com": // 黃金屋
		concurrencyengine.ELog.LPrintf("Noveler Choose 黃金屋\n")
		return noveler.NewHjwzwNoveler(URLNovel), nil
	case "tw.uukanshu.com": // UU看書網 tw
		u.Host = "www.uukanshu.com"
		fmt.Printf("%s => %s\n", URLNovel, u.String())
		URLNovel = u.String()
		fallthrough
	case "www.uukanshu.com": // UU看書網
		concurrencyengine.ELog.LPrintf("Noveler Choose UU看書網\n")
		return noveler.NewUUkanshuNoveler(URLNovel), nil
	case "www.ptwxz.com": // 飄天文學
		concurrencyengine.ELog.LPrintf("Noveler Choose 飄天文學\n")
		return noveler.NewPtwxzNoveler(URLNovel), nil
	default:
		return nil, fmt.Errorf("%s No useful interface", u.Host)
	}
}

// getNovel 取得小說內容
func getNovel(novel noveler.Noveler) error {
	tmpPath := "temp"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}

	resultPath := "finish"
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		os.MkdirAll(resultPath, os.ModePerm)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqToKey := func(req concurrencyengine.Request) interface{} { return req.Item.(noveler.NovelChapter) }
	e := concurrencyengine.New(ctx, 500, reqToKey)
	e.Recorder.JsonRWSetup(jsonUnmarshal, jsonMarshal)

	err := novel.GetInfo()
	if err != nil {
		return errors.Wrap(err, "novel.GetInfo")
	}

	// 讀取記錄
	fileNames := []string{}
	filePath := path.Join(tmpPath, fmt.Sprintf("%s-record.dat", novel.GetName()))
	m, err := e.Recorder.Load(filePath)
	if err != nil {
		return errors.Wrap(err, "e.Recorder.Load")
	}
	var novelPagesRecord []noveler.NovelChapter
	for k := range m {
		novelPagesRecord = append(novelPagesRecord, k.(noveler.NovelChapter))
	}

	var requests []concurrencyengine.Request
	for i := range novelPagesRecord {
		req := concurrencyengine.Request{Item: novelPagesRecord[i], ParseFunc: novel.GetParseResult}
		if !e.Recorder.IsDone(req) {
			requests = append(requests, req)
		} else {
			fileNames = append(fileNames, novelPagesRecord[i].Order)
		}
	}

	// 取得章節網址
	novelPages, err := novel.GetChapterURLs()
	if err != nil {
		concurrencyengine.ELog.Printf("novel.GetChapterURLs Fail: %s\n", err)
		return errors.Wrap(err, "novel.GetChapterURLs")
	}

	for i := range novelPages {
		req := concurrencyengine.Request{Item: novelPages[i], ParseFunc: novel.GetParseResult}
		if !e.Recorder.IsProcessed(req) {
			requests = append(requests, req)
		}
	}

	// 網址傳進 engine 抓取 HTML，並將小說內容存檔
	dataChan := e.Run(requests...)
	for data := range dataChan {
		// 不加 .txt 以免檔名排序錯誤，導致合併出錯
		fileName := fmt.Sprintf("%s", data.(noveler.NovelChapterHTML).Order)
		fileNames = append(fileNames, fileName)
		filePath := path.Join(tmpPath, fileName)

		if _, err := os.Stat(filePath); os.IsExist(err) {
			continue
		}

		text := data.(noveler.NovelChapterHTML).Text
		err = ioutil.WriteFile(filePath, []byte(text), os.ModePerm)
		if err != nil {
			concurrencyengine.ELog.Printf("ioutil.WriteFile Fail: %s\n", err)
			return errors.Wrap(err, "ioutil.WriteFile")
		}
		concurrencyengine.ELog.LPrintf("Chapter Content Write to %s\n", fileName)
		e.Recorder.Done(data.(noveler.NovelChapterHTML).NovelChapter)

		err = e.Recorder.Save(path.Join(tmpPath, fmt.Sprintf("%s-record.dat", novel.GetName())))
		if err != nil {
			return errors.Wrap(err, "e.Recorder.Save")
		}
	}

	// 合併暫存檔
	sort.Strings(fileNames)
	if err := novel.MergeContent(fileNames, tmpPath, resultPath); err != nil {
		return err
	}

	// 移除暫存檔
	for _, fileName := range fileNames {
		filePath := path.Join(tmpPath, fileName)
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	if err := os.Remove(path.Join(tmpPath, fmt.Sprintf("%s-record.dat", novel.GetName()))); err != nil {
		return err
	}

	return nil
}
