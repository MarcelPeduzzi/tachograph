package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalDetailedSpeedGen2 parses Gen2 Detailed Speed data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen2 Detailed Speed structure uses RecordArray format.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Gen2 has no V2 variant - both V1 and V2 use the same structure.
func unmarshalDetailedSpeedGen2(value []byte) (*vuv1.DetailedSpeedGen2, error) {
	// Split transfer value into data and signature
	// Gen2 uses variable-length ECDSA signatures stored as SignatureRecordArray
	// We use the sizeOf function to determine where to split
	totalSize, signatureSize, err := sizeOfDetailedSpeedGen2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	detailedSpeed := &vuv1.DetailedSpeedGen2{}
	detailedSpeed.SetRawData(value) // Store complete transfer value for painting

	// Validate structure by skipping through all record arrays
	offset := 0
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(data, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// VuDetailedSpeedRecordArray
	if err := skipRecordArray("VuDetailedSpeed"); err != nil {
		return nil, err
	}

	// Store signature (extracted at the beginning)
	detailedSpeed.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Detailed Speed Gen2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return detailedSpeed, nil
}

// MarshalDetailedSpeedGen2 marshals Gen2 Detailed Speed data using raw data painting.
func (opts MarshalOptions) MarshalDetailedSpeedGen2(detailedSpeed *vuv1.DetailedSpeedGen2) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	raw := detailedSpeed.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Detailed Speed Gen2 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeDetailedSpeedGen2 anonymizes Gen2 Detailed Speed data.
// TODO: Implement full semantic anonymization (anonymize speed records if needed).
func (opts AnonymizeOptions) anonymizeDetailedSpeedGen2(ds *vuv1.DetailedSpeedGen2) *vuv1.DetailedSpeedGen2 {
	if ds == nil {
		return nil
	}
	result := proto.Clone(ds).(*vuv1.DetailedSpeedGen2)
	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})

	// Note: We intentionally keep raw_data here because MarshalDetailedSpeedGen2
	// currently requires raw_data (semantic marshalling not yet implemented).

	return result
}
