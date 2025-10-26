package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalEventsAndFaultsGen1 parses Gen1 Events and Faults data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Events and Faults structure (from Data Dictionary and Appendix 7, Section 2.2.6.4 and 2.2.6.5):
//
// ASN.1 Definition:
//
//	VuEventsAndFaultsFirstGen ::= SEQUENCE {
//	    vuFaultData          VuFaultDataFirstGen,
//	    vuEventData          VuEventDataFirstGen,
//	    vuOverSpeedingControlData    VuOverSpeedingControlData,
//	    vuTimeAdjustmentData VuTimeAdjustmentDataFirstGen,
//	    signature            SignatureFirstGen
//	}
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Full semantic parsing is TODO.
func unmarshalEventsAndFaultsGen1(value []byte) (*vuv1.EventsAndFaultsGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	eventsAndFaults := &vuv1.EventsAndFaultsGen1{}
	eventsAndFaults.SetRawData(value) // Store complete transfer value for painting

	// TODO: Implement full semantic parsing of data portion
	// For now, we just validate we have some data
	_ = data // Will be used when semantic parsing is implemented

	// Store signature (extracted at the beginning)
	eventsAndFaults.SetSignature(signature)

	return eventsAndFaults, nil
}

// MarshalEventsAndFaultsGen1 marshals Gen1 Events and Faults data using raw data painting.
func (opts MarshalOptions) MarshalEventsAndFaultsGen1(eventsAndFaults *vuv1.EventsAndFaultsGen1) ([]byte, error) {
	if eventsAndFaults == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	raw := eventsAndFaults.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	// TODO: Implement marshalling from semantic fields
	return nil, fmt.Errorf("cannot marshal Events and Faults Gen1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeEventsAndFaultsGen1 anonymizes Gen1 Events and Faults data.
// TODO: Implement full semantic anonymization (anonymize event/fault records, timestamps, etc.).
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen1(ef *vuv1.EventsAndFaultsGen1) *vuv1.EventsAndFaultsGen1 {
	if ef == nil {
		return nil
	}
	result := proto.Clone(ef).(*vuv1.EventsAndFaultsGen1)
	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Note: We intentionally keep raw_data here because MarshalEventsAndFaultsGen1
	// currently requires raw_data (semantic marshalling not yet implemented).
	// Once semantic marshalling is implemented, we should clear raw_data and
	// implement full semantic anonymization of event/fault records.

	return result
}
