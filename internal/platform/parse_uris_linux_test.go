//go:build linux

package platform

import (
	"reflect"
	"testing"
)

func TestParseFileURIs(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "single file uri",
			in:   "file:///tmp/report.pdf\n",
			want: []string{"/tmp/report.pdf"},
		},
		{
			name: "multiple uris",
			in:   "file:///tmp/a.txt\nfile:///tmp/b.txt\n",
			want: []string{"/tmp/a.txt", "/tmp/b.txt"},
		},
		{
			name: "crlf uri list",
			in:   "file:///tmp/a.txt\r\nfile:///tmp/b.txt\r\n",
			want: []string{"/tmp/a.txt", "/tmp/b.txt"},
		},
		{
			name: "non file scheme ignored",
			in:   "https://example.com/file.txt\nfile:///tmp/c.txt\n",
			want: []string{"/tmp/c.txt"},
		},
		{
			name: "malformed and empty lines",
			in:   "\nnot a uri\nfile:///tmp/d.txt\n",
			want: []string{"/tmp/d.txt"},
		},
		{
			name: "percent encoded path",
			in:   "file:///tmp/My%20File.txt\n",
			want: []string{"/tmp/My File.txt"},
		},
		{
			name: "empty input",
			in:   "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFileURIs(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseFileURIs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
