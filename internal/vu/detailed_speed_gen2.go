package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// unmarshalDetailedSpeedGen2 parses Gen2 Detailed Speed data from the complete transfer value.
//
// Gen2 Detailed Speed structure uses RecordArray format.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Gen2 has no V2 variant - both V1 and V2 use the same structure.
func unmarshalDetailedSpeedGen2(value []byte) (*vuv1.DetailedSpeedGen2, error) {
	detailedSpeed := &vuv1.DetailedSpeedGen2{}
	detailedSpeed.SetRawData(value)

	// Validate structure by skipping through all record arrays
	offset := 0
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
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

	// SignatureRecordArray is now handled separately in raw parsing, not part of value

	if offset != len(value) {
		return nil, fmt.Errorf("Detailed Speed Gen2 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
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
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Detailed Speed Gen2 without raw_data (semantic marshalling not yet implemented)")
}
