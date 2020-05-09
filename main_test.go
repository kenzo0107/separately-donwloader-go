package main

import (
	"testing"
)

func TestDownload(t *testing.T) {
	type args struct {
		targetURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully download",
			args: args{
				targetURL: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "failed to download",
			args: args{
				targetURL: "http://%$#'a.b.c.d.e",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetURL = tt.args.targetURL
			if err := download(); (err != nil) != tt.wantErr {
				t.Errorf("download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesize(t *testing.T) {
	type args struct {
		targetURL string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "access to http://example.com",
			args: args{
				targetURL: "http://example.com",
			},
			want:    648,
			wantErr: false,
		},
		{
			name: "cannot access to not exist http://%$#'a.b.c.d.e",
			args: args{
				targetURL: "http://%$#'a.b.c.d.e",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "cannot access to not exist http://a.b.c.d.e",
			args: args{
				targetURL: "http://a.b.c.d.e",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filesize(tt.args.targetURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("filesize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filesize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeAccessRequest(t *testing.T) {
	dst = "test"

	type args struct {
		min int
		max int
		i   int
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "access to http://example.com",
			args: args{
				min: 0,
				max: 5,
				i:   0,
				url: "http://example.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := rangeAccessRequest(tt.args.min, tt.args.max, tt.args.i, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("rangeAccessRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBindwithFiles(t *testing.T) {
	type args struct {
		filename string
		filesize int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "access to http://example.com",
			args: args{
				filename: "0",
				filesize: 6,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bindwithFiles(tt.args.filename, tt.args.filesize); (err != nil) != tt.wantErr {
				t.Errorf("bindwithFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
