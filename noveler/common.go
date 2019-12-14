package noveler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/util"
)

// getParseResult 獲得 章節的內容 & 下一頁的連結
func getParseResult(novel Noveler, req crawler.Request) (crawler.ParseResult, error) {
	parseResult := crawler.ParseResult{
		Item:     nil,
		Requests: []crawler.Request{},
		DoneN:    0,
	}

	url := req.Item.(NovelChapter).URL
	r, name, certain, err := util.URLHTMLToUTF8Encoding(url)
	if err != nil {
		fmt.Printf("GetParseResult: util.URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		parseResult.Requests = append(parseResult.Requests, req)
		parseResult.DoneN = 0
		return parseResult, errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		parseResult.Requests = append(parseResult.Requests, req)
		parseResult.DoneN = 0
		return parseResult, errors.Wrap(err, "ioutil.ReadAll")
	}

	requests, err := novel.getNextPage(bytes.NewReader(b), req)
	if err != nil {
		parseResult.Requests = append(parseResult.Requests, req)
		parseResult.DoneN = 0
		return parseResult, errors.Wrap(err, "GetNextPage")
	}
	parseResult.Requests = append(parseResult.Requests, requests...)

	text, err := novel.getText(bytes.NewReader(b))
	if err != nil {
		parseResult.Requests = append(parseResult.Requests, req)
		parseResult.DoneN = 0
		return parseResult, errors.Wrap(err, "GetText")
	}

	parseResult.Item = NovelChapterHTML{
		NovelChapter: NovelChapter{
			Order: req.Item.(NovelChapter).Order,
			URL:   req.Item.(NovelChapter).URL,
		},
		Text: text}
	parseResult.DoneN = 1
	return parseResult, nil
}

// mergeContent 合併章節
func mergeContent(novelName string, fileNames []string, fromPath, toPath string) error {
	savePath := path.Join(toPath, novelName)

	f, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("ioutil.ReadFile Fail: %s\n", err)
		return errors.Wrap(err, "ioutil.ReadFile")
	}

	defer f.Close()

	for _, fName := range fileNames {
		fPath := path.Join(fromPath, fName)

		b, err := ioutil.ReadFile(fPath)
		if err != nil {
			fmt.Printf("ioutil.ReadFile Fail: %s\n", err)
			return errors.Wrap(err, "ioutil.ReadFile")
		}

		_, err = f.WriteString(string(b))
		if err != nil {
			fmt.Printf("ioutil.WriteFile Fail: %s\n", err)
			return errors.Wrap(err, "ioutil.WriteFile")
		}
	}

	return nil
}
