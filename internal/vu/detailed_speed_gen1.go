package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalDetailedSpeedGen1 parses Gen1 Detailed Speed data from the complete transfer value.
//
// Gen1 Detailed Speed structure (from Data Dictionary and Appendix 7, Section 2.2.6.6):
//
// ASN.1 Definition:
//
//	VuDetailedSpeedFirstGen ::= SEQUENCE {
//	    vuDetailedSpeedBlocks      VuDetailedSpeedBlocksFirstGen,
//	    signature                  SignatureFirstGen
//	}
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalDetailedSpeedGen1(value []byte) (*vuv1.DetailedSpeedGen1, error) {
	detailedSpeed := &vuv1.DetailedSpeedGen1{}
	detailedSpeed.SetRawData(value)

	// TODO: Implement full semantic parsing
	// For now, validate that we have enough data for the structure
	if len(value) < 128 { // At minimum, signature is 128 bytes
		return nil, fmt.Errorf("insufficient data for Detailed Speed Gen1")
	}

	// Store the signature (last 128 bytes)
	signatureStart := len(value) - 128
	detailedSpeed.SetSignature(value[signatureStart:])

	return detailedSpeed, nil
}

// MarshalDetailedSpeedGen1 marshals Gen1 Detailed Speed data using raw data painting.
func (opts MarshalOptions) MarshalDetailedSpeedGen1(detailedSpeed *vuv1.DetailedSpeedGen1) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	raw := detailedSpeed.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Detailed Speed Gen1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeDetailedSpeedGen1 anonymizes Gen1 Detailed Speed data.
// TODO: Implement full anonymization logic for Gen1 detailed speed.
func (opts AnonymizeOptions) anonymizeDetailedSpeedGen1(ds *vuv1.DetailedSpeedGen1) *vuv1.DetailedSpeedGen1 {
	if ds == nil {
		return nil
	}
	result := proto.Clone(ds).(*vuv1.DetailedSpeedGen1)
	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))
	result.SetRawData(nil)
	return result
}
