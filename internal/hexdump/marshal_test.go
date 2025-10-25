package hexdump

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty",
			input: []byte{},
			want:  "",
		},
		{
			name:  "single byte",
			input: []byte{0x48},
			want:  "00000000  48                                                |H|\n",
		},
		{
			name:  "hello world",
			input: []byte("Hello World!"),
			want:  "00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21              |Hello World!|\n",
		},
		{
			name:  "three bytes",
			input: []byte{0x01, 0x02, 0x03},
			want:  "00000000  01 02 03                                          |...|\n",
		},
		{
			name:  "exactly 16 bytes",
			input: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
			want:  "00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|\n",
		},
		{
			name:  "17 bytes (two lines)",
			input: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			want: "00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|\n" +
				"00000010  10                                                |.|\n",
		},
		{
			name:  "non-printable characters",
			input: []byte{0x00, 0x01, 0x1f, 0x20, 0x7e, 0x7f, 0xff},
			want:  "00000000  00 01 1f 20 7e 7f ff                              |... ~..|\n",
		},
		{
			name:  "mixed printable and non-printable",
			input: []byte{'A', 0x00, 'B', 0xff, 'C'},
			want:  "00000000  41 00 42 ff 43                                    |A.B.C|\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() unexpected error: %v", err)
			}
			gotStr := string(got)
			if diff := cmp.Diff(tt.want, gotStr); diff != "" {
				t.Errorf("Marshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
