package hexdump

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "empty",
			input: "",
			want:  []byte{},
		},
		{
			name:  "standard format with ASCII",
			input: "00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21              |Hello World!|\n",
			want:  []byte("Hello World!"),
		},
		{
			name:  "without ASCII column",
			input: "00000000  48 65 6c 6c 6f\n",
			want:  []byte("Hello"),
		},
		{
			name: "multiple lines",
			input: strings.Join(
				[]string{
					"00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|",
					"00000010  10                                                |.|",
				}, "\n") + "\n",
			want: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		},
		{
			name:  "different spacing (no double space at byte 8)",
			input: "00000000  48 65 6c 6c 6f 20 57 6f 72 6c 64 21\n",
			want:  []byte("Hello World!"),
		},
		{
			name:  "uppercase hex",
			input: "00000000  48 65 6C 6C 6F\n",
			want:  []byte("Hello"),
		},
		{
			name:  "with empty lines",
			input: "00000000  48 65 6c\n\n00000003  6c 6f\n",
			want:  []byte("Hello"),
		},
		{
			name:  "compact format (no spaces between bytes)",
			input: "00000000  48656c6c6f\n",
			want:  []byte("Hello"),
		},
		{
			name:  "wrong offsets (should be ignored)",
			input: "99999999  48 65 6c 6c 6f\n",
			want:  []byte("Hello"),
		},
		{
			name: "32-byte input across two lines",
			input: "00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|\n" +
				"00000010  10 11 12 13 14 15 16 17  18 19 1a 1b 1c 1d 1e 1f  |................|\n",
			want: []byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("Unmarshal() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unmarshal() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
