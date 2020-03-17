package noveler

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/util"
)

// CzbooksNoveler 小說狂人的 Noveler
type CzbooksNoveler struct {
	URL    string
	title  string
	author string
}

// NewCzbooksNoveler 建立 CzbooksNoveler
func NewCzbooksNoveler(url string) *CzbooksNoveler {
	var noveler CzbooksNoveler
	noveler.URL = url

	return &noveler
}

// GetInfo 獲得小說基本資料
func (n *CzbooksNoveler) GetInfo() error {
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	n.title = dom.Find("span.title").Text()
	n.title = strings.TrimSpace(n.title)
	n.author = strings.Replace(dom.Find("span.author").Text(), "作者: ", "", 1)
	n.author = strings.TrimSpace(n.author)

	return nil
}

// GetChapterURLs 獲得所有章節的網址
func (n *CzbooksNoveler) GetChapterURLs() ([]NovelChapter, error) {
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
	dom.Find("ul.nav.chapter-list > li > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u.Opaque = href
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%s-%010d", n.title, i), URL: u.String()})
			fmt.Printf("NovelPage %010d: %s\n", i, u.String())
		}
	})

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *CzbooksNoveler) GetParseResult(req crawler.Request) (crawler.ParseResult, error) {
	return getParseResult(n, req)
}

// GetName 回傳目前抓取的小說名字
func (n *CzbooksNoveler) GetName() string {
	novelName := fmt.Sprintf("%s-作者：%s", n.title, n.author)
	return novelName
}

// getNextPage 獲得下一頁的連結
func (n *CzbooksNoveler) getNextPage(html io.Reader, req crawler.Request) ([]crawler.Request, error) {
	return []crawler.Request{}, nil
}

// getText 獲得章節的內容
func (n *CzbooksNoveler) getText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("div.name").Text()

	text := dom.Find("div.content").Text()
	text = util.FormatText(text)

	return fmt.Sprintf("%s\n\n%s\n\n\n\n\n", chapterTitle, text), nil
}

// MergeContent 合併章節
func (n *CzbooksNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := n.GetName() + ".txt"
	return mergeContent(novelName, fileNames, fromPath, toPath)
}
