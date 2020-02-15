package util

import (
	"io"
	"testing"
)

func TestURLHTMLToUTF8Encoding(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name        string
		args        args
		wantR       io.Reader
		wantName    string
		wantCertain bool
		wantErr     bool
	}{
		// TODO: Add test cases.
		{"GBK", args{URL: "https://www.wanbentxt.com/18868/"}, nil, "gbk", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotName, gotCertain, err := URLHTMLToUTF8Encoding(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLHTMLToUTF8Encoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("URLHTMLToUTF8Encoding() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotCertain != tt.wantCertain {
				t.Errorf("URLHTMLToUTF8Encoding() gotCertain = %v, want %v", gotCertain, tt.wantCertain)
			}
		})
	}
}

func TestFormatText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"Test", args{"      一行\n 　 兩行\n       三行\n\n"}, "  一行\n\n  兩行\n\n  三行"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatText(tt.args.text); got != tt.want {
				t.Errorf("FormatText() = %q, want %q", got, tt.want)
			}
		})
	}
}
