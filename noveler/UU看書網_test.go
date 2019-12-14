package noveler

import (
	"io"
	"net/http"
	"testing"

	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/util"
)

func TestUUkanshuNoveler_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		n       *UUkanshuNoveler
		want    UUkanshuNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &UUkanshuNoveler{URL: "https://www.uukanshu.com/b/81074/"}, UUkanshuNoveler{title: "原来我是妖二代", author: "卖报小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUkanshuNoveler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("UUkanshuNoveler = %v, want %v", tt.n, tt.want)
			}
		})
	}
}
func TestUUkanshuNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *UUkanshuNoveler
		want    UUkanshuNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &UUkanshuNoveler{URL: "https://www.uukanshu.com/b/81074/"}, UUkanshuNoveler{title: "原來我是妖二代", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUkanshuNoveler.GetChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("UUkanshuNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestUUkanshuNoveler_getText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *UUkanshuNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &UUkanshuNoveler{URL: "https://www.uukanshu.com/b/81074/42310.html"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.n.URL)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.getText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("UUkanshuNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("UUkanshuNoveler.getText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUkanshuNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		n       *UUkanshuNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test",
			&UUkanshuNoveler{title: "原來我是妖二代", author: "賣報小郎君"},
			args{fileNames: []string{"1.txt", "2.txt"}, fromPath: "./temp", toPath: "./finish"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("UUkanshuNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUUkanshuNoveler_GetParseResult(t *testing.T) {
	type args struct {
		req crawler.Request
	}
	tests := []struct {
		name    string
		n       *UUkanshuNoveler
		args    args
		want    crawler.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&UUkanshuNoveler{},
			args{
				req: crawler.Request{
					Item: NovelChapter{
						URL:   "https://www.uukanshu.com/b/81074/42310.html",
						Order: "0001",
					}}},
			crawler.ParseResult{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetParseResult(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UUkanshuNoveler.GetParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.DoneN != -len(got.Requests)+1 {
				t.Errorf("UUkanshuNoveler.GetParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
