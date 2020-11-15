package noveler

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

// WanbentxtNoveler 完本神站的 Noveler
type WanbentxtNoveler struct {
	URL    string
	title  string
	author string
}

// NewWanbentxtNoveler 建立 WanbentxtNoveler
func NewWanbentxtNoveler(url string) *WanbentxtNoveler {
	var noveler WanbentxtNoveler
	noveler.URL = url

	return &noveler
}

// GetInfo 獲得小說基本資料
func (n *WanbentxtNoveler) GetInfo() error {
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	n.title = dom.Find("div.detailTitle > h1").Text()
	n.title = strings.Trim(n.title, " ")
	n.author = dom.Find("div.writer").Text()

	return nil
}

// GetChapterURLs 獲得所有章節的網址
func (n *WanbentxtNoveler) GetChapterURLs() ([]NovelChapter, error) {
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return nil, errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
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
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%s-%010d", n.title, i), URL: u.String()})
			fmt.Printf("NovelPage %010d: %s\n", i, u.String())
		}
	})

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *WanbentxtNoveler) GetParseResult(req concurrencyengine.Request) (concurrencyengine.ParseResult, error) {
	return getParseResult(n, req)
}

// GetName 回傳目前抓取的小說名字
func (n *WanbentxtNoveler) GetName() string {
	novelName := fmt.Sprintf("%s-作者：%s", n.title, n.author)
	return novelName
}

// getNextPage 獲得下一頁的連結
func (n *WanbentxtNoveler) getNextPage(html io.Reader, req concurrencyengine.Request) ([]concurrencyengine.Request, error) {
	requests := []concurrencyengine.Request{}

	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return requests, errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	if s := dom.Find("span.next"); s.Text() == "下一页" {
		href, ok := s.Parent().Attr("href")
		if !ok {
			outer, err := goquery.OuterHtml(s.Parent())
			concurrencyengine.ELog.Fatalf("%v %v\n", outer, err)
		}
		order := req.Item.(NovelChapter).Order + "-1"

		requests = append(requests, concurrencyengine.Request{
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

// getText 獲得章節的內容
func (n *WanbentxtNoveler) getText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("div.readerTitle").Text()
	chapterTitle = strings.TrimSpace(chapterTitle)

	text := dom.Find("div.readerCon").Text()
	text = strings.ReplaceAll(text, "一秒记住【完本神站】手机用户输入地址：m.wanbentxt.com", "")
	text = strings.ReplaceAll(text, "-->>本章未完，点击下一页继续阅读", "")
	text = util.FormatText(text)

	return fmt.Sprintf("%s\n\n%s\n\n\n\n\n", chapterTitle, text), nil
}

// MergeContent 合併章節
func (n *WanbentxtNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := n.GetName() + ".txt"
	return mergeContent(novelName, fileNames, fromPath, toPath)
}
