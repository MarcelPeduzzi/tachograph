package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalEventsAndFaultsGen1 parses Gen1 Events and Faults data from the complete transfer value.
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
	eventsAndFaults := &vuv1.EventsAndFaultsGen1{}
	eventsAndFaults.SetRawData(value)

	// TODO: Implement full semantic parsing
	// For now, validate that we have enough data for the structure
	if len(value) < 128 { // At minimum, signature is 128 bytes
		return nil, fmt.Errorf("insufficient data for Events and Faults Gen1")
	}

	// Store the signature (last 128 bytes)
	signatureStart := len(value) - 128
	eventsAndFaults.SetSignature(value[signatureStart:])

	return eventsAndFaults, nil
}

// MarshalEventsAndFaultsGen1 marshals Gen1 Events and Faults data using raw data painting.
func (opts MarshalOptions) MarshalEventsAndFaultsGen1(eventsAndFaults *vuv1.EventsAndFaultsGen1) ([]byte, error) {
	if eventsAndFaults == nil {
		return nil, fmt.Errorf("eventsAndFaults cannot be nil")
	}

	raw := eventsAndFaults.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	// TODO: Implement marshalling from semantic fields
	return nil, fmt.Errorf("cannot marshal Events and Faults Gen1 without raw_data (semantic marshalling not yet implemented)")
}



// anonymizeEventsAndFaultsGen1 anonymizes Gen1 Events and Faults data.
// TODO: Implement full anonymization logic for Gen1 events/faults.
func (opts AnonymizeOptions) anonymizeEventsAndFaultsGen1(ef *vuv1.EventsAndFaultsGen1) *vuv1.EventsAndFaultsGen1 {
	if ef == nil {
		return nil
	}
	result := proto.Clone(ef).(*vuv1.EventsAndFaultsGen1)
	result.SetSignature(nil)
	result.SetRawData(nil)
	return result
}
