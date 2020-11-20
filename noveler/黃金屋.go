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

// HjwzwNoveler 黃金屋的 Noveler
type HjwzwNoveler struct {
	URL    string
	title  string
	author string
}

// NewHjwzwNoveler 建立 HjwzwNoveler
func NewHjwzwNoveler(url string) *HjwzwNoveler {
	var noveler HjwzwNoveler
	noveler.URL = url

	return &noveler
}

// GetInfo 獲得小說基本資料
func (n *HjwzwNoveler) GetInfo() error {
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	n.title = dom.Find("td > h1").Text()
	n.title = strings.Trim(n.title, " ")
	n.author = dom.Find("body > div:first-child > table:nth-of-type(7) tr:nth-child(2) a:first-child").Text()

	return nil
}

// GetChapterURLs 獲得所有章節的網址
func (n *HjwzwNoveler) GetChapterURLs() ([]NovelChapter, error) {
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
	dom.Find("div#tbchapterlist a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u.Path = href
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%s-%010d", n.title, i), URL: u.String()})
			fmt.Printf("NovelPage %010d: %s\n", i, u.String())
		}
	})

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *HjwzwNoveler) GetParseResult(req concurrencyengine.Request) (concurrencyengine.ParseResult, error) {
	return getParseResult(n, req)
}

// GetName 回傳目前抓取的小說名字
func (n *HjwzwNoveler) GetName() string {
	novelName := fmt.Sprintf("%s-作者：%s", n.title, n.author)
	return novelName
}

// getNextPage 獲得下一頁的連結
func (n *HjwzwNoveler) getNextPage(html io.Reader, req concurrencyengine.Request) ([]concurrencyengine.Request, error) {
	requests := []concurrencyengine.Request{}

	return requests, nil
}

// getText 獲得章節的內容
func (n *HjwzwNoveler) getText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("table:nth-of-type(7) h1").Text()
	chapterTitle = strings.TrimSpace(chapterTitle)

	// text := dom.Find("table:nth-of-type(7) div:nth-of-type(5)").Text()
	dom.Find("div#Pan_Ad1").Remove()
	text := dom.Find("table:nth-of-type(7) div:nth-of-type(4)").Text()
	text = strings.ReplaceAll(text, "請記住本站域名: 黃金屋", "")
	text = strings.ReplaceAll(text, "，歡迎訪問大家讀書院", "")
	text = util.FormatText(text)

	return util.MergeTitle(text,chapterTitle), nil
}

// MergeContent 合併章節
func (n *HjwzwNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := n.GetName() + ".txt"
	return mergeContent(novelName, fileNames, fromPath, toPath)
}
