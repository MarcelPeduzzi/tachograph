package hexdump

import (
	"bytes"
	"encoding/hex"
)

// Marshal converts binary data into hexdump format matching `hexdump -C`.
// The output format is:
//
//	00000000  48 65 6c 6c 6f 20 57 6f  72 6c 64 21              |Hello World!|
//	0000000c  01 02 03                                          |...|
//
// Each line contains 16 bytes of data with:
//   - 8-digit hex offset (lowercase, zero-padded)
//   - Two spaces separator
//   - Hex bytes (lowercase, space-separated, double space after byte 8)
//   - ASCII representation (printable chars or '.' for non-printable)
func Marshal(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	const bytesPerLine = 16

	for offset := 0; offset < len(data); offset += bytesPerLine {
		// Write offset
		buf.WriteString(hex.EncodeToString([]byte{
			byte(offset >> 24),
			byte(offset >> 16),
			byte(offset >> 8),
			byte(offset),
		}))
		buf.WriteString("  ")

		// Get the chunk for this line (up to 16 bytes)
		end := offset + bytesPerLine
		if end > len(data) {
			end = len(data)
		}
		chunk := data[offset:end]

		// Write hex bytes
		for i, b := range chunk {
			if i == 8 {
				buf.WriteString(" ") // Extra space after 8th byte
			}
			buf.WriteString(hex.EncodeToString([]byte{b}))
			buf.WriteString(" ")
		}

		// Pad hex section to fixed width (50 chars for full line, matching hexdump -C)
		// Each byte takes 3 chars (2 hex + 1 space), plus 1 extra space at position 8
		// Full line: 16*3 + 1 = 49 chars, plus 1 space before | = 50 chars total
		hexChars := len(chunk)*3 + 1 // +1 for extra space at byte 8
		if len(chunk) <= 8 {
			hexChars = len(chunk) * 3
		}
		padding := 50 - hexChars // 50 = 16*3 + 1 (extra space) + 1 (space before |)
		for i := 0; i < padding; i++ {
			buf.WriteByte(' ')
		}

		// Write ASCII column
		buf.WriteByte('|')
		for _, b := range chunk {
			if b >= 0x20 && b <= 0x7e {
				buf.WriteByte(b)
			} else {
				buf.WriteByte('.')
			}
		}
		buf.WriteString("|\n")
	}

	return buf.Bytes(), nil
}
