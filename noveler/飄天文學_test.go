package noveler

import (
	"io"
	"testing"

	"github.com/z-Wind/concurrencyengine"
	"github.com/z-Wind/getNovel/util"
)

func TestPtwxzNoveler_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		n       *PtwxzNoveler
		want    PtwxzNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9795/index.html"}, PtwxzNoveler{title: "原来我是妖二代", author: "卖报小郎君"}, false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9640/"}, PtwxzNoveler{title: "全球高武", author: "老鹰吃小鸡"}, false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.piaotia.com/bookinfo/9/9640.html"}, PtwxzNoveler{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.n.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("PtwxzNoveler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("PtwxzNoveler = %#v, want %#v", tt.n, tt.want)
			}
		})
	}
}
func TestPtwxzNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *PtwxzNoveler
		want    PtwxzNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9795/index.html"}, PtwxzNoveler{title: "原來我是妖二代", author: "賣報小郎君"}, false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9640/"}, PtwxzNoveler{title: "全球高武", author: "老鹰吃小鸡"}, false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.piaotia.com/bookinfo/9/9640.html"}, PtwxzNoveler{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("PtwxzNoveler.GetChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) == 0 {
				t.Errorf("PtwxzNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestPtwxzNoveler_getText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *PtwxzNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"原來我是妖二代", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9795/6694578.html"}, args{}, "", false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9640/6523755.html"}, args{}, "", false},
		{"全球高武", &PtwxzNoveler{URL: "https://www.ptwxz.com/html/9/9640/6770680.html"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _, _, err := util.URLHTMLToUTF8Encoding(tt.n.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("util.URLHTMLToUTF8Encoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.args.html = r

			got, err := tt.n.getText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("PtwxzNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("PtwxzNoveler.getText() = %v, want %v", got, tt.want)
			} else {
				t.Logf("PtwxzNoveler.getText() = %v", got)
			}
		})
	}
}

func TestPtwxzNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		n       *PtwxzNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test",
			&PtwxzNoveler{title: "原來我是妖二代", author: "賣報小郎君"},
			args{fileNames: []string{"1.txt", "2.txt"}, fromPath: "./temp", toPath: "./finish"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("PtwxzNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPtwxzNoveler_GetParseResult(t *testing.T) {
	type args struct {
		req concurrencyengine.Request
	}
	tests := []struct {
		name    string
		n       *PtwxzNoveler
		args    args
		want    concurrencyengine.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test",
			&PtwxzNoveler{},
			args{
				req: concurrencyengine.Request{
					Item: NovelChapter{
						URL:   "https://www.ptwxz.com/html/9/9795/6694578.html",
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
				t.Errorf("PtwxzNoveler.GetParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Done || len(got.ExtraRequests) != 0 || len(got.RedoRequests) != 0 {
				t.Errorf("PtwxzNoveler.GetParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
