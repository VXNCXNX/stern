package stern

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestConsumeFileTail(t *testing.T) {
	logLines := `line 1
line 2
line 3
line 4`
	tmpl := template.Must(template.New("").Parse(`{{printf "%s\n" .Message}}`))

	tests := []struct {
		name      string
		resumeReq *ResumeRequest
		expected  []byte
	}{
		{
			name: "normal",
			expected: []byte(`line 1
line 2
line 3
line 4
`),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			tail := NewFileTail(tmpl, nil, out, io.Discard, &TailOptions{})
			if err := tail.ConsumeReader(bufio.NewReader(strings.NewReader(logLines))); err != nil {
				t.Fatalf("%d: unexpected err %v", i, err)
			}

			if !bytes.Equal(tt.expected, out.Bytes()) {
				t.Errorf("%d: expected %s, but actual %s", i, tt.expected, out)
			}
		})
	}
}

func TestConsumeFileTailTimestamps(t *testing.T) {
	logLine := "test line\n"
	tmpl := template.Must(template.New("").Parse(`{{if .Timestamp}}{{.Timestamp}} {{end}}{{printf "%s\n" .Message}}`))

	tests := []struct {
		name       string
		timestamps bool
		expected   string
	}{
		{
			name:       "with timestamps",
			timestamps: true,
			// TimestampFormatDefault always contains a "T" separator, e.g.
			// "2006-01-02T15:04:05.000000000Z07:00".
			expected: "T",
		},
		{
			name:       "without timestamps",
			timestamps: false,
			expected:   "test line\n",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			tail := NewFileTail(tmpl, nil, out, io.Discard, &TailOptions{Timestamps: tt.timestamps, Location: time.UTC})
			if err := tail.ConsumeReader(bufio.NewReader(strings.NewReader(logLine))); err != nil {
				t.Fatalf("%d: unexpected err %v", i, err)
			}

			if tt.timestamps {
				if !strings.Contains(out.String(), tt.expected) {
					t.Errorf("%d: expected output to contain %q, but got %q", i, tt.expected, out.String())
				}
			} else if out.String() != tt.expected {
				t.Errorf("%d: expected %q, but actual %q", i, tt.expected, out.String())
			}
		})
	}
}
