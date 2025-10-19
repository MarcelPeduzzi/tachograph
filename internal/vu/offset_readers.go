package vu

import (
	"encoding/binary"
	"io"
)

// VU-specific offset-based binary parsing functions for reading structured data from byte slices.
//
// These functions provide a low-level, offset-based approach to parsing binary data.
// They are used extensively throughout the VU package for parsing complex structures
// where the bufio.Scanner pattern is not suitable (e.g., when parsing non-contiguous
// fields or when the structure layout is not uniform).
//
// NOTE: The preferred pattern for contiguous binary data parsing is bufio.Scanner
// with custom SplitFunc (see AGENTS.md for details). These offset-based functions
// should only be used when the scanner pattern is not applicable.
//
// Usage pattern:
//   offset := 0
//   value, offset, err := readUint8FromBytes(data, offset)
//   if err != nil { return err }
//   // Continue with next field...

// readUint8FromBytes reads a single byte from a byte slice at the given offset
func readUint8FromBytes(data []byte, offset int) (uint8, int, error) {
	if offset >= len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	return data[offset], offset + 1, nil
}

// readBytesFromBytes reads n bytes from a byte slice at the given offset
func readBytesFromBytes(data []byte, offset int, n int) ([]byte, int, error) {
	if offset+n > len(data) {
		return nil, offset, io.ErrUnexpectedEOF
	}
	result := make([]byte, n)
	copy(result, data[offset:offset+n])
	return result, offset + n, nil
}

// readVuTimeRealFromBytes reads a TimeReal value (4 bytes) and converts to Unix timestamp
func readVuTimeRealFromBytes(data []byte, offset int) (int64, int, error) {
	if offset+4 > len(data) {
		return 0, offset, io.ErrUnexpectedEOF
	}
	value := binary.BigEndian.Uint32(data[offset:])
	// TimeReal is seconds since 00:00:00 UTC, 1 January 1970
	return int64(value), offset + 4, nil
}
