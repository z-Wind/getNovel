package noveler

import (
	"fmt"
	"strings"
	"time"
)

// Timeout Context 的 timeout
const Timeout = 60 * time.Second

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
