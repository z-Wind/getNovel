package main

import (
	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/noveler"
)

// Noveler 抓取小說必需的 function
type Noveler interface {
	// 獲得所有章節的網址
	GetChapterURLs() ([]noveler.NovelChapter, error)
	// 獲得 章節的內容 下一頁的連結
	GetParseResult(req crawler.Request) (crawler.ParseResult, error)
	// 合併章節
	MergeContent(fileNames []string, fromPath, toPath string) error
}
