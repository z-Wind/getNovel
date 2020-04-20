package noveler

import (
	"io"
	"net/http"
	"testing"

	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

func TestHjwzwNoveler_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		n       *HjwzwNoveler
		want    HjwzwNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &HjwzwNoveler{URL: "https://tw.hjwzw.com/Book/Chapter/37176"}, HjwzwNoveler{title: "原來我是妖二代", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("HjwzwNoveler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("HjwzwNoveler = %v, want %v", tt.n, tt.want)
			}
		})
	}
}
func TestHjwzwNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *HjwzwNoveler
		want    HjwzwNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &HjwzwNoveler{URL: "https://tw.hjwzw.com/Book/Chapter/37176"}, HjwzwNoveler{title: "原來我是妖二代", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("HjwzwNoveler.GetChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("HjwzwNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestHjwzwNoveler_getText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *HjwzwNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &HjwzwNoveler{URL: "https://tw.hjwzw.com/Book/Read/37176,15476202"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.n.URL)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.getText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("HjwzwNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("HjwzwNoveler.getText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHjwzwNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		n       *HjwzwNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test",
			&HjwzwNoveler{title: "原來我是妖二代", author: "賣報小郎君"},
			args{fileNames: []string{"1.txt", "2.txt"}, fromPath: "./temp", toPath: "./finish"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("HjwzwNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHjwzwNoveler_GetParseResult(t *testing.T) {
	type args struct {
		req concurrencyengine.Request
	}
	tests := []struct {
		name    string
		n       *HjwzwNoveler
		args    args
		want    concurrencyengine.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&HjwzwNoveler{},
			args{
				req: concurrencyengine.Request{
					Item: NovelChapter{
						URL:   "https://tw.hjwzw.com/Book/Chapter/37176",
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
				t.Errorf("HjwzwNoveler.GetParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Done || len(got.ExtraRequests) != 0 || len(got.RedoRequests) != 0 {
				t.Errorf("HjwzwNoveler.GetParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
