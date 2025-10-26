package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalEventsAndFaultsGen2V1 parses Gen2 V1 Events and Faults data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen2 V1 Events and Faults structure uses RecordArray format.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalEventsAndFaultsGen2V1(value []byte) (*vuv1.EventsAndFaultsGen2V1, error) {
	// Split transfer value into data and signature
	// Gen2 uses variable-length ECDSA signatures stored as SignatureRecordArray
	// We use the sizeOf function to determine where to split
	totalSize, signatureSize, err := sizeOfEventsAndFaultsGen2V1(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	eventsAndFaults := &vuv1.EventsAndFaultsGen2V1{}
	eventsAndFaults.SetRawData(value) // Store complete transfer value for painting

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

	// Skip all record arrays
	// VuFaultRecordArray
	if err := skipRecordArray("VuFault"); err != nil {
		return nil, err
	}
	// VuEventRecordArray
	if err := skipRecordArray("VuEvent"); err != nil {
		return nil, err
	}
	// VuOverSpeedingControlRecordArray
	if err := skipRecordArray("VuOverSpeedingControl"); err != nil {
		return nil, err
	}
	// VuTimeAdjustmentRecordArray
	if err := skipRecordArray("VuTimeAdjustment"); err != nil {
		return nil, err
	}

	// Store signature (extracted at the beginning)
	eventsAndFaults.SetSignature(signature)

	if offset != len(data) {
		return nil, fmt.Errorf("Events and Faults Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return eventsAndFaults, nil
}

// MarshalEventsAndFaultsGen2V1 marshals Gen2 V1 Events and Faults data using raw data painting.
func (opts MarshalOptions) MarshalEventsAndFaultsGen2V1(eventsAndFaults *vuv1.EventsAndFaultsGen2V1) ([]byte, error) {
	if eventsAndFaults == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	raw := eventsAndFaults.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Events and Faults Gen2 V1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeEventsAndFaultsGen2V1 anonymizes Gen2 V1 Events and Faults data.
// TODO: Implement full semantic anonymization (anonymize event/fault records, timestamps, etc.).
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen2V1(ef *vuv1.EventsAndFaultsGen2V1) *vuv1.EventsAndFaultsGen2V1 {
	if ef == nil {
		return nil
	}
	result := proto.Clone(ef).(*vuv1.EventsAndFaultsGen2V1)
	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})

	// Note: We intentionally keep raw_data here because MarshalEventsAndFaultsGen2V1
	// currently requires raw_data (semantic marshalling not yet implemented).

	return result
}
