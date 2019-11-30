package noveler

import (
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
