package dd

// readInt24 reads a 3-byte big-endian signed integer.
//
// This function handles the conversion from 3-byte binary data to a signed 32-bit
// integer, with proper sign extension. The tachograph protocol uses 24-bit signed
// integers for various fields including geographic coordinates.
//
// Parameters:
//   - data: A 3-byte slice containing the big-endian 24-bit integer
//
// Returns:
//   - A signed 32-bit integer with the value sign-extended from 24 bits
//
// The sign extension works by checking bit 23 (the most significant bit of the
// 24-bit value). If it's set, the upper 8 bits are filled with 1s to maintain
// the correct negative value representation in 32-bit two's complement.
func readInt24(data []byte) int32 {
	// Read as unsigned 24-bit value
	val := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	// Sign extend from 24 bits to 32 bits
	// If bit 23 is set (negative number), extend with 1s
	if val&0x800000 != 0 {
		val |= 0xFF000000
	}
	return int32(val)
}

// marshalInt24 converts a signed 32-bit integer to 3-byte big-endian bytes.
//
// This function handles the conversion from a signed 32-bit integer to 3-byte
// binary data in big-endian format. Only the lower 24 bits of the input value
// are written, with higher bits being truncated.
//
// Parameters:
//   - val: A signed 32-bit integer to convert
//
// Returns:
//   - A 3-byte slice containing the big-endian representation of the lower 24 bits
//
// Note: Values outside the 24-bit signed range (-8,388,608 to 8,388,607) will
// be truncated to fit within 24 bits. This is the expected behavior for the
// tachograph protocol which uses 24-bit fields.
func marshalInt24(val int32) []byte {
	// Write the lower 24 bits in big-endian order
	return []byte{byte(val >> 16), byte(val >> 8), byte(val)}
}
