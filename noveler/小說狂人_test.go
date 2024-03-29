package noveler

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

func TestCzbooksNoveler_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		want    CzbooksNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &CzbooksNoveler{URL: "https://czbooks.net/n/u5a6m"}, CzbooksNoveler{title: "《原來我是妖二代》", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("CzbooksNoveler = %v, want %v", tt.n, tt.want)
			}
		})
	}
}

func TestCzbooksNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		want    CzbooksNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"瘟疫医生", &CzbooksNoveler{URL: "https://czbooks.net/n/u5a6m"}, CzbooksNoveler{title: "《原來我是妖二代》", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.GetChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("CzbooksNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestCzbooksNoveler_getText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &CzbooksNoveler{URL: "https://czbooks.net/n/u5a6m/uj6h"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.n.URL)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.getText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("CzbooksNoveler.getText() = %v, want %v", got, tt.want)
			}else {
				t.Logf("CzbooksNoveler.getText() = %v", got)
			}
		})
	}
}

func TestCzbooksNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test",
			&CzbooksNoveler{title: "《原來我是妖二代》", author: "賣報小郎君"},
			args{fileNames: []string{"1.txt", "2.txt"}, fromPath: "./temp", toPath: "./finish"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCzbooksNoveler_getNextPage(t *testing.T) {
	type args struct {
		html io.Reader
		req  concurrencyengine.Request
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		url     string
		args    args
		want    []concurrencyengine.Request
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&CzbooksNoveler{},
			"https://czbooks.net/n/u5a6m/uj6h",
			args{req: concurrencyengine.Request{Item: NovelChapter{Order: "0001"}}},
			[]concurrencyengine.Request{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.url)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.getNextPage(tt.args.html, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.getNextPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WanbentxtNoveler.getNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCzbooksNoveler_GetParseResult(t *testing.T) {
	type args struct {
		req concurrencyengine.Request
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		args    args
		want    concurrencyengine.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&CzbooksNoveler{},
			args{
				req: concurrencyengine.Request{
					Item: NovelChapter{
						URL:   "https://czbooks.net/n/u5a6m/uj6h",
						Order: "0001",
					}}},
			concurrencyengine.ParseResult{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetParseResult(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.GetParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Done || len(got.ExtraRequests) != 0 || len(got.RedoRequests) != 0 {
				t.Errorf("WanbentxtNoveler.GetParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
