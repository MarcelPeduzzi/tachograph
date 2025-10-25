package hexdump

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty",
			data: []byte{},
		},
		{
			name: "single byte",
			data: []byte{0x42},
		},
		{
			name: "hello world",
			data: []byte("Hello World!"),
		},
		{
			name: "all printable ASCII",
			data: []byte("The quick brown fox jumps over the lazy dog."),
		},
		{
			name: "binary data",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		},
		{
			name: "exactly 16 bytes",
			data: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f},
		},
		{
			name: "17 bytes",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		},
		{
			name: "32 bytes (two full lines)",
			data: []byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
		},
		{
			name: "mixed printable and non-printable",
			data: []byte{'A', 0x00, 'B', 0xff, 'C', 0x0a, 'D'},
		},
		{
			name: "all non-printable",
			data: []byte{0x00, 0x01, 0x02, 0x7f, 0x80, 0xff},
		},
		{
			name: "larger data (100 bytes)",
			data: func() []byte {
				data := make([]byte, 100)
				for i := range data {
					data[i] = byte(i)
				}
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to hexdump
			marshaled, err := Marshal(tt.data)
			if err != nil {
				t.Fatalf("Marshal() unexpected error: %v", err)
			}
			// Unmarshal back to binary
			unmarshaled, err := Unmarshal(marshaled)
			if err != nil {
				t.Fatalf("Unmarshal() unexpected error: %v", err)
			}
			// Compare original with roundtrip result
			if diff := cmp.Diff(tt.data, unmarshaled); diff != "" {
				t.Errorf("Roundtrip mismatch (-want +got):\n%s", diff)
				t.Logf("Marshaled output:\n%s", string(marshaled))
			}
		})
	}
}
