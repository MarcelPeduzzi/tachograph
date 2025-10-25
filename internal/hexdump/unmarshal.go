package hexdump

import (
	"bufio"
	"bytes"
	"encoding/hex"
)

// Unmarshal converts hexdump format back to binary data.
// It accepts any hexdump format with offsets and hex bytes, ignoring:
//   - Offset values (not validated)
//   - ASCII columns (anything after the hex data)
//   - Empty lines
//   - Arbitrary spacing between hex bytes
//
// This makes it forgiving and able to parse dumps from various sources.
func Unmarshal(data []byte) ([]byte, error) {
	result := []byte{}
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// Buffer for collecting hex characters (reused across lines)
	hexBuf := make([]byte, 0, 32)

	for scanner.Scan() {
		line := scanner.Bytes()

		// Trim leading and trailing whitespace
		line = bytes.TrimSpace(line)

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Find the separator between offset and hex data (two spaces)
		sepIdx := bytes.Index(line, []byte("  "))
		if sepIdx == -1 {
			// No separator found, skip this line
			continue
		}

		// Extract everything after the offset
		hexPart := line[sepIdx+2:]

		// Remove ASCII column if present (everything from '|' onwards)
		if pipeIdx := bytes.IndexByte(hexPart, '|'); pipeIdx != -1 {
			hexPart = hexPart[:pipeIdx]
		}

		// Extract hex bytes by filtering out spaces
		hexBuf = hexBuf[:0] // Reset buffer, keep capacity
		for _, b := range hexPart {
			// Skip whitespace, keep only hex characters
			if b != ' ' && b != '\t' {
				hexBuf = append(hexBuf, b)
			}
		}

		if len(hexBuf) == 0 {
			continue
		}

		// Decode hex bytes
		decoded := make([]byte, hex.DecodedLen(len(hexBuf)))
		n, err := hex.Decode(decoded, hexBuf)
		if err != nil {
			// Skip lines that can't be decoded as hex
			continue
		}

		result = append(result, decoded[:n]...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
