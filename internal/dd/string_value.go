package dd

import (
	"fmt"
	"io"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalStringValue unmarshals a code-paged string value from binary data.
// The input should contain a code page byte followed by the encoded string data.
//
// The data type `StringValue` is specified in the Data Dictionary, Section 2.158.
//
// ASN.1 Definition:
//
//	StringValue ::= SEQUENCE {
//	    codePage    OCTET STRING (SIZE(1)),
//	    stringData  OCTET STRING (SIZE(0..255))
//	}
func (opts UnmarshalOptions) UnmarshalStringValue(input []byte) (*ddv1.StringValue, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("insufficient data for string value: %w", io.ErrUnexpectedEOF)
	}

	codePage := input[0]
	data := input[1:]

	var output ddv1.StringValue
	output.SetEncoding(getEncodingFromCodePage(codePage))
	// Store the entire input including the code page byte (aligns with raw data painting policy)
	if opts.PreserveRawData {
		output.SetRawData(input)
	}
	// Length field represents the string data length (not including code page)
	output.SetLength(int32(len(data)))

	// Decode the string based on the code page
	decoded, err := decodeWithCodePage(codePage, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode string with code page %d: %w", codePage, err)
	}
	output.SetValue(decoded)

	return &output, nil
}

// MarshalStringValue marshals a StringValue to bytes.
//
// This function handles code-paged string format defined as:
//
//	StringValue ::= SEQUENCE {
//	    codePage    OCTET STRING (SIZE(1)),
//	    stringData  OCTET STRING (SIZE(0..255))
//	}
//
// Binary Layout: codePage (1 byte) + stringData (variable or fixed-length bytes)
//
// If 'raw_data' is available, it is used directly (for round-trip fidelity).
// Otherwise, the 'value' string is encoded using the specified encoding and padded
// with spaces if a 'length' field is set (for fixed-length strings).
func (opts MarshalOptions) MarshalStringValue(sv *ddv1.StringValue) ([]byte, error) {
	// Handle nil
	if sv == nil {
		// Empty string value: code page 255 (EMPTY) + no data
		return []byte{0xFF}, nil
	}

	// Determine the expected total length (code page + string data)
	// Length field represents string data only (without code page)
	hasFixedLength := sv.HasLength()
	var totalLength int
	if hasFixedLength {
		// Total = 1 (code page) + length (string data)
		totalLength = 1 + int(sv.GetLength())
	}

	// Validate that raw_data has correct size if present
	if sv.HasRawData() {
		rawData := sv.GetRawData()
		if len(rawData) < 1 {
			return nil, fmt.Errorf("raw_data must include at least the code page byte")
		}
		if hasFixedLength && len(rawData) != totalLength {
			return nil, fmt.Errorf("raw_data length (%d) does not match expected total length (%d = 1 + %d)", len(rawData), totalLength, sv.GetLength())
		}
	}

	// Determine the code page byte from the encoding field
	codePage := getCodePageFromEncoding(sv.GetEncoding())

	// Use raw data painting approach if raw_data is available
	if raw := sv.GetRawData(); len(raw) > 0 {
		// Allocate canvas from raw_data
		canvas := make([]byte, len(raw))
		copy(canvas, raw)

		// Paint only the code page byte at offset 0 (from semantic encoding field)
		// The string data at offset 1+ is already correct in raw_data
		canvas[0] = codePage

		// Note: We do NOT re-encode from the value field because:
		// 1. The value field is UTF-8 (for display), while raw_data is in the original encoding
		// 2. Re-encoding from UTF-8 can produce different byte lengths (e.g., ISO-8859-2 → UTF-8 → ISO-8859-2)
		// 3. The raw_data preserves the exact original bytes, which is what we want for round-trip fidelity

		return canvas, nil
	}

	// Fallback: encode from semantic fields (no raw_data available)
	value := sv.GetValue()
	encoded, err := encodeWithCodePage(codePage, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode string with code page %d: %w", codePage, err)
	}

	// Create result buffer
	var result []byte
	// Write code page byte
	result = append(result, codePage)

	// Handle fixed-length strings (pad with spaces)
	if hasFixedLength {
		dataLength := int(sv.GetLength())
		if len(encoded) > dataLength {
			return nil, fmt.Errorf("encoded string length (%d) exceeds allowed length (%d)", len(encoded), dataLength)
		}
		for i := 0; i < len(encoded); i++ {
			result = append(result, encoded[i])
		}
		for i := len(encoded); i < dataLength; i++ {
			result = append(result, ' ')
		}
		return result, nil
	}

	// Handle variable-length strings (no padding)
	return append(result, encoded...), nil
}
