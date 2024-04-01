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

// PtwxzNoveler 飄天文學的 Noveler
type PtwxzNoveler struct {
	URL    string
	title  string
	author string
}

// NewPtwxzNoveler 建立 PtwxzNoveler
func NewPtwxzNoveler(url string) *PtwxzNoveler {
	var noveler PtwxzNoveler
	noveler.URL = url

	return &noveler
}

// GetInfo 獲得小說基本資料
func (n *PtwxzNoveler) GetInfo() error {
	if strings.Contains(n.URL, "bookinfo") {
		return errors.New("You are using the wrong URL. Please provide the chapter list, which should look like /html/9/9795/index.html")
	}
	r, name, certain, err := util.URLHTMLToUTF8Encoding(n.URL)
	if err != nil {
		fmt.Printf("URLHTMLToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return errors.Wrap(err, "util.URLHTMLToUTF8Encoding")
	}

	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	n.title = dom.Find("div.title h1").Text()
	n.title = strings.ReplaceAll(n.title, "最新章节", "")
	n.title = strings.Trim(n.title, " ")

	dom.Find("div.mainbody > div.list a").Remove()
	n.author = dom.Find("div.mainbody > div.list").Text()
	n.author = strings.ReplaceAll(n.author, "作者：", "")
	n.author = strings.ReplaceAll(n.author, "\u00a0", "")
	n.author = strings.ReplaceAll(n.author, "\n", "")
	n.author = strings.Trim(n.author, " ")

	return nil
}

// GetChapterURLs 獲得所有章節的網址
func (n *PtwxzNoveler) GetChapterURLs() ([]NovelChapter, error) {
	if strings.Contains(n.URL, "bookinfo") {
		return nil, errors.New("You are using the wrong URL. Please provide the chapter list, which should look like /html/9/9795/index.html")
	}

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
	sel := dom.Find("div.centent li a")
	for i := range sel.Nodes {
		s := sel.Eq(i)
		if href, ok := s.Attr("href"); ok {
			link, err := util.ToAbsoluteURL(u.String(), href)
			if err != nil {
				return chapters, errors.Wrap(err, "util.ToAbsoluteURL")
			}
			chapters = append(chapters, NovelChapter{Order: fmt.Sprintf("%s-%010d", n.title, i), URL: link})
			fmt.Printf("NovelPage %010d: %s\n", i, link)
		}
	}

	return chapters, nil
}

// GetParseResult 獲得 章節的內容 & 下一頁的連結
func (n *PtwxzNoveler) GetParseResult(req concurrencyengine.Request) (concurrencyengine.ParseResult, error) {
	return getParseResult(n, req)
}

// GetName 回傳目前抓取的小說名字
func (n *PtwxzNoveler) GetName() string {
	novelName := fmt.Sprintf("%s-作者：%s", n.title, n.author)
	return novelName
}

// getNextPage 獲得下一頁的連結
func (n *PtwxzNoveler) getNextPage(html io.Reader, req concurrencyengine.Request) ([]concurrencyengine.Request, error) {
	requests := []concurrencyengine.Request{}

	return requests, nil
}

// getText 獲得章節的內容
func (n *PtwxzNoveler) getText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	dom.Find("body h1 a").Remove()
	chapterTitle := dom.Find("body h1").Text()
	chapterTitle = strings.TrimSpace(chapterTitle)

	dom.Find("body div").Remove()
	dom.Find("body script").Remove()
	dom.Find("body > table").Remove()
	dom.Find("body > h1").Remove()
	dom.Find("body center").Remove()

	text, err := dom.Find("body").Html()
	if err != nil {
		return "", errors.Wrap(err, "dom.Find")
	}
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = regexp.MustCompile(`.*<.*>.*`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`.*小说网友请提示:.*`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`&gt;`).ReplaceAllString(text, "")
	text = util.FormatText(text)

	return util.MergeTitle(text, chapterTitle), nil
}

// MergeContent 合併章節
func (n *PtwxzNoveler) MergeContent(fileNames []string, fromPath, toPath string) error {
	novelName := n.GetName() + ".txt"
	return mergeContent(novelName, fileNames, fromPath, toPath)
}
