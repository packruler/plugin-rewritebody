package compressutil_test

import (
	"bytes"
	"testing"

	"github.com/packruler/rewrite-body/compressutil"
)

type TestStruct struct {
	desc        string
	input       []byte
	expected    []byte
	encoding    string
	shouldMatch bool
}

var (
	deflatedBytes = []byte{74, 203, 207, 87, 200, 44, 86, 40, 201, 72, 85, 200, 75, 45, 87, 72, 74, 44, 2, 4, 0, 0, 255, 255}
	gzippedBytes  = []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 74, 203, 207, 87, 200, 44, 86, 40, 201, 72, 85, 200, 75, 45, 87, 72, 74, 44, 2, 4, 0, 0, 255, 255, 251, 28, 166, 187, 18, 0, 0, 0}
	normalBytes   = []byte("foo is the new bar")
)

func TestEncode(t *testing.T) {
	tests := []TestStruct{
		{
			desc:        "should support identity",
			input:       normalBytes,
			expected:    normalBytes,
			encoding:    "identity",
			shouldMatch: true,
		},
		{
			desc:        "should support gzip",
			input:       normalBytes,
			expected:    gzippedBytes,
			encoding:    "gzip",
			shouldMatch: false,
		},
		{
			desc:        "should support deflate",
			input:       normalBytes,
			expected:    deflatedBytes,
			encoding:    "deflate",
			shouldMatch: false,
		},
		{
			desc:        "should NOT support brotli",
			input:       normalBytes,
			expected:    normalBytes,
			encoding:    "br",
			shouldMatch: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			output, err := compressutil.Encode([]byte(test.input), test.encoding)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			isBad := !bytes.Equal([]byte(test.expected), output)

			if isBad {
				t.Errorf("expected error got body: %v\n wanted: %v", output, []byte(test.expected))
			}

			if test.shouldMatch {
				isBad = !bytes.Equal([]byte(test.input), output)
			} else {
				isBad = bytes.Equal([]byte(test.input), output)
			}
			if isBad {
				t.Errorf("match error got body: %v\n wanted: %v", output, []byte(test.input))
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []TestStruct{
		{
			desc:        "should support identity",
			input:       normalBytes,
			expected:    normalBytes,
			encoding:    "identity",
			shouldMatch: true,
		},
		{
			desc:        "should support gzip",
			input:       gzippedBytes,
			expected:    normalBytes,
			encoding:    "gzip",
			shouldMatch: false,
		},
		{
			desc:        "should support deflate",
			input:       deflatedBytes,
			expected:    normalBytes,
			encoding:    "deflate",
			shouldMatch: false,
		},
		{
			desc:        "should NOT support brotli",
			input:       normalBytes,
			expected:    normalBytes,
			encoding:    "br",
			shouldMatch: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			output, err := compressutil.Decode(bytes.NewBuffer([]byte(test.input)), test.encoding)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			isBad := !bytes.Equal([]byte(test.expected), output)

			if isBad {
				t.Errorf("expected error got body: %v\n wanted: %v", output, []byte(test.expected))
			}

			if test.shouldMatch {
				isBad = !bytes.Equal([]byte(test.input), output)
			} else {
				isBad = bytes.Equal([]byte(test.input), output)
			}
			if isBad {
				t.Errorf("match error got body: %s\n wanted: %s", output, []byte(test.input))
			}
		})
	}
}

func compressString(value string, encoding string) string {
	compressed, _ := compressutil.Encode([]byte(value), encoding)

	return string(compressed)
}
