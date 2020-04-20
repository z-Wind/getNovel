package main

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/z-Wind/getNovel/noveler"
)

func jsonUnmarshal(b []byte) (map[interface{}]bool, error) {
	result := make(map[interface{}]bool)

	m := make(map[noveler.NovelChapter]bool)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return result, errors.Wrap(err, "json.Unmarshal(")
	}

	for key, item := range m {
		result[key] = item

	}
	return result, nil
}

func jsonMarshal(m map[interface{}]bool) ([]byte, error) {
	result := make(map[noveler.NovelChapter]bool)

	for key, item := range m {
		result[key.(noveler.NovelChapter)] = item
	}
	jsonBytes, err := json.Marshal(&result)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	return jsonBytes, nil
}
