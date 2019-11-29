package util

import (
	"io"
	"testing"
)

func TestURLToUTF8Encoding(t *testing.T) {
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
			_, gotName, gotCertain, err := URLToUTF8Encoding(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlToUtf8Encoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("UrlToUtf8Encoding() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotCertain != tt.wantCertain {
				t.Errorf("UrlToUtf8Encoding() gotCertain = %v, want %v", gotCertain, tt.wantCertain)
			}
		})
	}
}
