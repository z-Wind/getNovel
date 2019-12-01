package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/noveler"
)

func Test_record(t *testing.T) {
	tmpPath := "temp"
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}
	url := "http://123"
	order := "00001"
	req := crawler.Request{
		Item: noveler.NovelChapter{Order: order, URL: url},
	}

	r := NewRecord()
	if got := r.checkExistOrAdd(req); got {
		t.Errorf("r.checkExist() = %v, want %v", got, false)
	}
	r.done(req.Item.(noveler.NovelChapter))
	if got := r.checkExistOrAdd(req); !got {
		t.Errorf("r.checkExist() = %v, want %v", got, true)
	}
	if err := r.saveExist(tmpPath); err != nil {
		t.Errorf("r.saveExist() error = %v, wantErr %v", err, nil)
	}

	r = NewRecord()
	if _, err := r.loadExist(tmpPath); err != nil {
		t.Errorf("r.saveExist() error = %v, wantErr %v", err, nil)
	}
	req2 := crawler.Request{
		Item: noveler.NovelChapter{Order: order, URL: url},
	}
	fmt.Printf("%+v\n", r.taskDone)
	if got := r.checkExistOrAdd(req2); !got {
		t.Errorf("r.checkExist() = %v, want %v", got, true)
	}

	// 移除暫存檔
	if err := os.RemoveAll(tmpPath); err != nil {
		t.Error(err)
	}

}
