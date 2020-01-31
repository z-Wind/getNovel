package noveler

import (
	"io"
	"net/http"
	"testing"

	"github.com/z-Wind/getNovel/crawler"
	"github.com/z-Wind/getNovel/util"
)

func TestWanbentxtNoveler_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		want    WanbentxtNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"瘟疫医生", &WanbentxtNoveler{URL: "https://www.wanbentxt.com/18868/"}, WanbentxtNoveler{title: "瘟疫医生", author: "机器人瓦力"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("WanbentxtNoveler = %v, want %v", tt.n, tt.want)
			}
		})
	}
}
func TestWanbentxtNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		want    WanbentxtNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"瘟疫医生", &WanbentxtNoveler{URL: "https://www.wanbentxt.com/18868/"}, WanbentxtNoveler{title: "瘟疫医生", author: "机器人瓦力"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.GetChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("WanbentxtNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestWanbentxtNoveler_getText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &WanbentxtNoveler{URL: "https://www.wanbentxt.com/18868/12250657.html"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.n.URL)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.getText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("WanbentxtNoveler.getText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWanbentxtNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test",
			&WanbentxtNoveler{title: "瘟疫医生", author: "机器人瓦力"},
			args{fileNames: []string{"1.txt", "2.txt"}, fromPath: "./temp", toPath: "./finish"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWanbentxtNoveler_getNextPage(t *testing.T) {
	type args struct {
		html io.Reader
		req  crawler.Request
	}
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		url     string
		args    args
		want    []crawler.Request
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&WanbentxtNoveler{},
			"https://www.wanbentxt.com/8895/5687694.html",
			args{req: crawler.Request{Item: NovelChapter{Order: "0001"}}},
			nil,
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
			if got[0].Item.(NovelChapter).Order != tt.args.req.Item.(NovelChapter).Order+"-1" {
				t.Errorf("WanbentxtNoveler.getNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWanbentxtNoveler_GetParseResult(t *testing.T) {
	type args struct {
		req crawler.Request
	}
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		args    args
		want    crawler.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&WanbentxtNoveler{},
			args{
				req: crawler.Request{
					Item: NovelChapter{
						URL:   "https://www.wanbentxt.com/8895/5687694.html",
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
				t.Errorf("WanbentxtNoveler.GetParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.DoneN != 1 || len(got.Requests) != 1 {
				t.Errorf("WanbentxtNoveler.GetParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
