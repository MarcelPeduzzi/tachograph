package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalDetailedSpeedGen1 parses Gen1 Detailed Speed data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
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
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	detailedSpeed := &vuv1.DetailedSpeedGen1{}
	detailedSpeed.SetRawData(value) // Store complete transfer value for painting

	// TODO: Implement full semantic parsing of data portion
	// For now, we just validate we have some data
	_ = data // Will be used when semantic parsing is implemented

	// Store signature (extracted at the beginning)
	detailedSpeed.SetSignature(signature)

	return detailedSpeed, nil
}

// MarshalDetailedSpeedGen1 marshals Gen1 Detailed Speed data using raw data painting.
func (opts MarshalOptions) MarshalDetailedSpeedGen1(detailedSpeed *vuv1.DetailedSpeedGen1) ([]byte, error) {
	if detailedSpeed == nil {
		return nil, fmt.Errorf("detailedSpeed cannot be nil")
	}

	raw := detailedSpeed.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Detailed Speed Gen1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeDetailedSpeedGen1 anonymizes Gen1 Detailed Speed data.
// TODO: Implement full semantic anonymization (anonymize speed records if needed).
func (opts AnonymizeOptions) anonymizeDetailedSpeedGen1(ds *vuv1.DetailedSpeedGen1) *vuv1.DetailedSpeedGen1 {
	if ds == nil {
		return nil
	}
	result := proto.Clone(ds).(*vuv1.DetailedSpeedGen1)
	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Note: We intentionally keep raw_data here because MarshalDetailedSpeedGen1
	// currently requires raw_data (semantic marshalling not yet implemented).
	// Once semantic marshalling is implemented, we should clear raw_data.

	return result
}
