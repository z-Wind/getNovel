package noveler

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/z-Wind/getNovel/crawler"
)

// Timeout Context 的 timeout
const Timeout = 60 * time.Second

// Noveler 抓取小說必需的 function
type Noveler interface {
	// 得到小說名字
	GetName() string
	// 獲得小說基本資料
	GetInfo() error
	// 獲得和所有章節的網址
	GetChapterURLs() ([]NovelChapter, error)
	// 獲得 章節的內容 下一頁的連結
	GetParseResult(req crawler.Request) (crawler.ParseResult, error)
	// 合併章節
	MergeContent(fileNames []string, fromPath, toPath string) error

	// getNextPage 獲得下一頁的連結
	getNextPage(html io.Reader, req crawler.Request) ([]crawler.Request, error)
	// getText 獲得章節的內容
	getText(html io.Reader) (string, error)
}

// NovelChapter 小說章節網址
type NovelChapter struct {
	Order string
	URL   string
}

// NovelChapterHTML 小說章節的 HTML
type NovelChapterHTML struct {
	NovelChapter
	Text string
}

// MarshalText 為了 map 的 key，所以 NovelChapter 不是 pointer
func (n NovelChapter) MarshalText() ([]byte, error) {
	key := fmt.Sprintf("%s,%s", n.Order, n.URL)

	return []byte(key), nil
}

// UnmarshalText 為了 unmarshal 後能填入值，所以是 *NovelChapter
func (n *NovelChapter) UnmarshalText(text []byte) error {
	str := strings.Split(string(text), ",")

	n.Order = str[0]
	n.URL = str[1]

	return nil
}
