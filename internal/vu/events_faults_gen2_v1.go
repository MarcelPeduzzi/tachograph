package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalEventsAndFaultsGen2V1 parses Gen2 V1 Events and Faults data from the complete transfer value.
//
// Gen2 V1 Events and Faults structure uses RecordArray format.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalEventsAndFaultsGen2V1(value []byte) (*vuv1.EventsAndFaultsGen2V1, error) {
	eventsAndFaults := &vuv1.EventsAndFaultsGen2V1{}
	eventsAndFaults.SetRawData(value)

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

	// SignatureRecordArray is now handled separately in raw parsing, not part of value

	if offset != len(value) {
		return nil, fmt.Errorf("Events and Faults Gen2 V1 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
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
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Events and Faults Gen2 V1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeEventsAndFaultsGen2V1 anonymizes Gen2 V1 Events and Faults data.
// TODO: Implement full anonymization logic for Gen2 V1 events/faults.
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen2V1(ef *vuv1.EventsAndFaultsGen2V1) *vuv1.EventsAndFaultsGen2V1 {
	if ef == nil {
		return nil
	}
	result := proto.Clone(ef).(*vuv1.EventsAndFaultsGen2V1)
	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})
	result.SetRawData(nil)
	return result
}
