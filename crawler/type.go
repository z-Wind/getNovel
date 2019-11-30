package main

import (
	"github.com/z-Wind/getNovel/noveler"
	"io"
)

// Request 需執行的任務
type Request noveler.NovelChapter

// ParseResult worker 回傳的執行結果
type ParseResult struct {
	Item     *noveler.NovelChapterHTML
	Requests []Request
	// 已執行完的任務數，用來扣除用
	doneN int
}

// Scheduler 調配工作
type Scheduler interface {
	Submit(Request)
	WorkerReady(chan Request)
	Run()
}

// Noveler 抓取小說必需的 function
type Noveler interface {
	// 獲得所有章節的網址
	GetChapterURLs() ([]noveler.NovelChapter, error)
	// 獲得章節的內容
	GetText(html io.Reader) (string, error)
	// 合併章節
	MergeContent(fromPath, toPath string) error
}
