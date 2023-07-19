package noveler

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

// UUkanshuNoveler UU看書網的 Noveler
type UUkanshuNoveler struct {
	URL    string
	title  string
	author string
}

// NewUUkanshuNoveler 建立 UUkanshuNoveler
func NewUUkanshuNoveler(url string) *UUkanshuNoveler {
	var noveler UUkanshuNoveler
	noveler.URL = url

	return &noveler
}

// GetInfo 獲得小說基本資料
func (n *UUkanshuNoveler) GetInfo() error {
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	n.title = dom.Find("dd.jieshao_content > h1 > a").Text()
	n.title = strings.ReplaceAll(n.title, "最新章节", "")
	n.title = strings.Trim(n.title, " ")
	n.author = dom.Find("dd.jieshao_content > h2 > a").Text()

	return nil
}

// GetChapterURLs 獲得所有章節的網址
func (n *UUkanshuNoveler) GetChapterURLs() ([]NovelChapter, error) {
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
	dom.Find("ul#chapterList a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u.Path = href
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%s-%010d", n.title, i), URL: u.String()})
			fmt.Printf("NovelPage %010d: %s\n", i, u.String())
		}
	})

	// reverse order
	for i, j := 0, len(chapters)-1; i < j; i, j = i+1, j-1 {
		chapters[i].Order, chapters[j].Order = chapters[j].Order, chapters[i].Order
	}

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *UUkanshuNoveler) GetParseResult(req concurrencyengine.Request) (concurrencyengine.ParseResult, error) {
	return getParseResult(n, req)
}

// GetName 回傳目前抓取的小說名字
func (n *UUkanshuNoveler) GetName() string {
	novelName := fmt.Sprintf("%s-作者：%s", n.title, n.author)
	return novelName
}

// getNextPage 獲得下一頁的連結
func (n *UUkanshuNoveler) getNextPage(html io.Reader, req concurrencyengine.Request) ([]concurrencyengine.Request, error) {
	requests := []concurrencyengine.Request{}

	return requests, nil
}

// getText 獲得章節的內容
func (n *UUkanshuNoveler) getText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("h1#timu").Text()
	chapterTitle = strings.TrimSpace(chapterTitle)

	dom.Find("div.ad_content").Remove()
	text, err := dom.Find("div#contentbox").Html()
	if err != nil {
		return "", errors.Wrap(err, `dom.Find("div#contentbox").Html`)
	}
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = regexp.MustCompile(`[wｗ]{3}[．\.][ｕu][ｕu][ｋk][ａa][ｎn][ｓs][ｈh][ｕu][．\.][ｃc][ｏo][ｍm]`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`[ｕu][ｕu]看书\w+?\n`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`[ｕu][ｕu]看书`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`<.*>`).ReplaceAllString(text, "")
	text = util.FormatText(text)

	return util.MergeTitle(text, chapterTitle), nil
}

// MergeContent 合併章節
func (n *UUkanshuNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := n.GetName() + ".txt"
	return mergeContent(novelName, fileNames, fromPath, toPath)
}
