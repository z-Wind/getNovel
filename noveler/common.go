package noveler

import (
	"bytes"
	"fmt"
	
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

// getParseResult 獲得 章節的內容 & 下一頁的連結
func getParseResult(novel Noveler, req concurrencyengine.Request) (concurrencyengine.ParseResult, error) {
	parseResult := concurrencyengine.ParseResult{
		Item:          nil,
		ExtraRequests: []concurrencyengine.Request{},
		RedoRequests:  []concurrencyengine.Request{},
		Done:          false,
	}

	url := req.Item.(NovelChapter).URL
	r, name, certain, err := util.URLHTMLToUTF8Encoding(url)
	if err != nil {
		fmt.Printf("GetParseResult: util.URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		parseResult.RedoRequests = append(parseResult.RedoRequests, req)
		parseResult.Done = false
		return parseResult, errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	b, err := io.ReadAll(r)
	if err != nil {
		parseResult.RedoRequests = append(parseResult.RedoRequests, req)
		parseResult.Done = false
		return parseResult, errors.Wrap(err, "io.ReadAll")
	}

	requests, err := novel.getNextPage(bytes.NewReader(b), req)
	if err != nil {
		parseResult.RedoRequests = append(parseResult.RedoRequests, req)
		parseResult.Done = false
		return parseResult, errors.Wrap(err, "GetNextPage")
	}
	parseResult.ExtraRequests = append(parseResult.ExtraRequests, requests...)

	text, err := novel.getText(bytes.NewReader(b))
	if err != nil {
		parseResult.RedoRequests = append(parseResult.RedoRequests, req)
		parseResult.Done = false
		return parseResult, errors.Wrap(err, "GetText")
	}

	parseResult.Item = NovelChapterHTML{
		NovelChapter: NovelChapter{
			Order: req.Item.(NovelChapter).Order,
			URL:   req.Item.(NovelChapter).URL,
		},
		Text: text}
	parseResult.Done = true
	return parseResult, nil
}

// mergeContent 合併章節
func mergeContent(novelName string, fileNames []string, fromPath, toPath string) error {
	savePath := path.Join(toPath, novelName)

	f, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("os.ReadFile Fail: %s\n", err)
		return errors.Wrap(err, "os.ReadFile")
	}

	defer f.Close()

	for _, fName := range fileNames {
		fPath := path.Join(fromPath, fName)

		b, err := os.ReadFile(fPath)
		if err != nil {
			fmt.Printf("os.ReadFile Fail: %s\n", err)
			return errors.Wrap(err, "os.ReadFile")
		}

		_, err = f.WriteString(string(b))
		if err != nil {
			fmt.Printf("os.WriteFile Fail: %s\n", err)
			return errors.Wrap(err, "os.WriteFile")
		}
	}

	return nil
}
