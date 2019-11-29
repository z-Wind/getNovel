package noveler

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/util"
)

// CzbooksNoveler 小說狂人的 Noveler
type CzbooksNoveler struct {
	URL      string
	title    string
	author   string
	numPages int
}

// GetChapterURLs 獲得所有章節的網址
func (n *CzbooksNoveler) GetChapterURLs() ([]NovelChapter, error) {
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
	dom.Find("ul.nav.chapter-list > li > a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u.Opaque = href
			chapters = append(chapters, NovelChapter{Order: i, URL: u.String()})
			fmt.Printf("NovelPage %d: %s\n", i, u.String())
		}
	})

	n.title = dom.Find("span.title").Text()
	n.title = strings.Trim(n.title, " ")
	n.author = strings.Replace(dom.Find("span.author").Text(), "作者: ", "", 1)
	n.numPages = len(chapters)

	return chapters, nil
}

// GetText 獲得章節的內容
func (n *CzbooksNoveler) GetText(html io.Reader) (string, error) {
	dom, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocumentFromReader")
	}

	chapterTitle := dom.Find("div.name").Text()
	text := dom.Find("div.content").Text()

	return fmt.Sprintf("%s\n\n%s", chapterTitle, text), nil
}

// MergeContent 合併章節
func (n *CzbooksNoveler) MergeContent(fromPath, toPath string) error {
	novelName := fmt.Sprintf("%s-作者：%s.txt", n.title, n.author)
	savePath := path.Join(toPath, novelName)

	f, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("ioutil.ReadFile Fail: %s\n", err)
		return errors.Wrap(err, "ioutil.ReadFile")
	}

	defer f.Close()

	for i := 0; i < n.numPages; i++ {
		fName := fmt.Sprintf("%d.txt", i)
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
