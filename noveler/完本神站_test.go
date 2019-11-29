package noveler

import (
	"io"
	"net/http"
	"testing"

	"github.com/z-Wind/getNovel/util"
)

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
				t.Errorf("WanbentxtNoveler.getChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 || tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("WanbentxtNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestWanbentxtNoveler_GetText(t *testing.T) {
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

			got, err := tt.n.GetText(tt.args.html)
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
		fromPath string
		toPath   string
	}
	tests := []struct {
		name    string
		n       *WanbentxtNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", &WanbentxtNoveler{title: "瘟疫医生", author: "机器人瓦力", numPages: 10}, args{fromPath: "./temp", toPath: "./finish"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("WanbentxtNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
