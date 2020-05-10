package main

import "testing"

func TestMain(t *testing.T) {
	type args struct {
		targetURL string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "access to http://example.com",
			args: args{
				targetURL: "http://example.com/",
			},
		},
	}
	for _, tt := range tests {
		targetURL = tt.args.targetURL
		t.Run(tt.name, func(t *testing.T) {
			_main()
		})
	}
}

func BenchmarkMain(b *testing.B) {
	b.ResetTimer()
	_main()
}
