package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/noveler"
)

type record struct {
	taskDone map[noveler.NovelChapter]bool
	lock     sync.Mutex
}

// newRecord 建立 record
func newRecord() *record {
	var r record

	r.taskDone = make(map[noveler.NovelChapter]bool)

	return &r
}

// loadExist 讀取記錄資料
func (r *record) loadExist(dirPath string) ([]noveler.NovelChapter, error) {
	filePath := path.Join(dirPath, "record.dat")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	err = json.Unmarshal(b, &r.taskDone)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal(")
	}

	fmt.Printf("load from %s\n", filePath)

	var chapters []noveler.NovelChapter
	for k := range r.taskDone {
		chapters = append(chapters, k)
	}

	return chapters, nil
}

// checkExistOrAdd 確認連結是否存在，不存在就加入
func (r *record) checkExistOrAdd(req interface{}) bool {
	key := req.(crawler.Request).Item.(noveler.NovelChapter)
	_, ok := r.taskDone[key]

	if !ok {
		r.lock.Lock()
		r.taskDone[key] = false
		r.lock.Unlock()
	}

	return ok
}

// checkDone 確認連結是否完成
func (r *record) checkDone(req interface{}) bool {
	order := req.(crawler.Request).Item.(noveler.NovelChapter).Order
	key := req.(crawler.Request).Item.(noveler.NovelChapter)

	fmt.Printf("NovelPage %s: %s Done\n", order, key)

	return r.taskDone[key]
}

// done 任務已完成
func (r *record) done(chapter noveler.NovelChapter) {
	key := chapter

	r.lock.Lock()
	r.taskDone[key] = true
	r.lock.Unlock()
}

// saveExist 儲存記錄資料
func (r *record) saveExist(dirPath string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	jsonString, err := json.Marshal(&r.taskDone)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	filePath := path.Join(dirPath, "record.dat")
	err = ioutil.WriteFile(filePath, jsonString, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "ioutil.WriteFile")
	}
	fmt.Printf("save to %s\n", filePath)

	return nil
}
