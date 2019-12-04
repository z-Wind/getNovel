package noveler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/util"
)

// WanbentxtNoveler 完本神站的 Noveler
type WanbentxtNoveler struct {
	URL      string
	title    string
	author   string
	numPages int
}

// NewWanbentxtNoveler 建立 WanbentxtNoveler
func NewWanbentxtNoveler(url string) *WanbentxtNoveler {
	var noveler WanbentxtNoveler
	noveler.URL = url

	return &noveler
}

// GetChapterURLs 獲得所有章節的網址
func (n *WanbentxtNoveler) GetChapterURLs() ([]NovelChapter, error) {
	// Create a new context with a deadline
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	resp, err := util.HTTPGetwithContext(ctx, n.URL)
	if err != nil {
		fmt.Printf("GetChapterURLs: HTTPGetwithContext(%s): %s\n", n.URL, err)
		return nil, errors.Wrap(err, "HTTPGetwithContext")
	}
	defer resp.Body.Close()

	// 編碼成 UTF8，goquery 指定編碼
	r, name, certain, err := util.ToUTF8Encoding(resp.Body)
	if err != nil {
		fmt.Printf("ToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return nil, errors.Wrap(err, "DetermineEncodingFromReader")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	u, err := url.Parse(n.URL)
	if err != nil {
		return nil, errors.Wrap(err, "url.Parse")
	}

	var chapters []NovelChapter
	dom.Find("div.chapter > ul > li > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u.Path = href
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%010d", i), URL: u.String()})
			fmt.Printf("NovelPage %010d: %s\n", i, u.String())
		}
	})

	n.title = dom.Find("div.detailTitle > h1").Text()
	n.title = strings.Trim(n.title, " ")
	n.author = dom.Find("div.writer").Text()
	n.numPages = len(chapters)

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *WanbentxtNoveler) GetParseResult(req crawler.Request) (crawler.ParseResult, error) {
	// Request the HTML page
	// Create a new context with a deadline
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	resp, err := util.HTTPGetwithContext(ctx, req.Item.(NovelChapter).URL)
	if err != nil {
		// fmt.Printf("ParseResult: HTTPGetwithContext(%s): %s\n", req.URL, err)
		return crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{req},
			DoneN:    0,
		}, errors.Wrap(err, "util.HTTPGetwithContext")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// fmt.Printf("ParseResult: HTTPGetwithContext(%s): status code error: %d %s\n", req.URL, resp.StatusCode, resp.Status)
		return crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{req},
			DoneN:    0,
		}, fmt.Errorf("util.HTTPGetwithContext(%s): status code error: %d %s", req.Item.(NovelChapter).URL, resp.StatusCode, resp.Status)
	}

	r, name, certain, err := util.ToUTF8Encoding(resp.Body)
	if err != nil {
		fmt.Printf("GetParseResult: util.ToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{req},
			DoneN:    0,
		}, errors.Wrap(err, "util.ToUTF8Encoding")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		// fmt.Printf("ParseResult: ioutil.ReadAll: err:%s\n", err)
		return crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{req},
			DoneN:    0,
		}, errors.Wrap(err, "ioutil.ReadAll")
	}

	requests, err := n.GetNextPage(bytes.NewReader(b), req)
	if err != nil {
		return crawler.ParseResult{
			Item:     nil,
			Requests: []crawler.Request{req},
			DoneN:    0,
		}, errors.Wrap(err, "GetNextPage")
	}

	text, err := n.GetText(bytes.NewReader(b))
	if err != nil {
		return crawler.ParseResult{
			Item:     nil,
			Requests: append(requests, req),
			DoneN:    -len(requests),
		}, errors.Wrap(err, "GetText")
	}

	return crawler.ParseResult{
		Item: NovelChapterHTML{
			NovelChapter: NovelChapter{
				Order: req.Item.(NovelChapter).Order,
				URL:   req.Item.(NovelChapter).URL,
			},
			Text: text},
		Requests: requests,
		DoneN:    -len(requests) + 1,
	}, nil
}

// GetNextPage 獲得下一頁的連結
func (n *WanbentxtNoveler) GetNextPage(html io.Reader, req crawler.Request) ([]crawler.Request, error) {
	requests := []crawler.Request{}

	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return requests, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	if s := dom.Find("span.next"); s.Text() == "下一页" {
		href, ok := s.Parent().Attr("href")
		if !ok {
			log.Fatal(goquery.OuterHtml(s.Parent()))
		}
		order := req.Item.(NovelChapter).Order + "-1"

		requests = append(requests, crawler.Request{
			Item: NovelChapter{
				Order: order,
				URL:   href,
			},
			ParseFunc: n.GetParseResult,
		})

		fmt.Printf("NovelPage %s: %s\n", order, href)
	}

	return requests, nil
}

// GetText 獲得章節的內容
func (n *WanbentxtNoveler) GetText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("div.readerTitle").Text()
	text := dom.Find("div.readerCon").Text()

	return fmt.Sprintf("%s\n\n%s\n\n\n\n\n", chapterTitle, text), nil
}

// MergeContent 合併章節
func (n *WanbentxtNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := fmt.Sprintf("%s-作者：%s.txt", n.title, n.author)
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
