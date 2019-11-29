package noveler

import (
	"io"
	"time"
)

// Timeout Context 的 timeout
const Timeout = 20 * time.Second

// NovelChapter 小說章節網址
type NovelChapter struct {
	Order int
	URL   string
}

// NovelChapterHTML 小說章節的 HTML
type NovelChapterHTML struct {
	*NovelChapter
	HTML io.Reader
}
