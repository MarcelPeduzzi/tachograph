package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalTechnicalDataGen1 parses Gen1 Technical Data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Technical Data structure (from Data Dictionary and Appendix 7, Section 2.2.6.7):
//
// ASN.1 Definition:
//
//	VuTechnicalDataFirstGen ::= SEQUENCE {
//	    vuApprovalNumber                VuApprovalNumber,
//	    vuSoftwareIdentification        VuSoftwareIdentification,
//	    vuManufacturerName              VuManufacturerName,
//	    vuManufacturerAddress           VuManufacturerAddress,
//	    vuPartNumber                    VuPartNumber,
//	    vuSerialNumber                  ExtendedSerialNumber,
//	    sensorPaired                    SensorPaired,
//	    signature                       SignatureFirstGen
//	}
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalTechnicalDataGen1(value []byte) (*vuv1.TechnicalDataGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	technicalData := &vuv1.TechnicalDataGen1{}
	technicalData.SetRawData(value) // Store complete transfer value for painting

	// TODO: Implement full semantic parsing of data portion
	// For now, we just validate we have some data
	_ = data // Will be used when semantic parsing is implemented

	// Store signature (extracted at the beginning)
	technicalData.SetSignature(signature)

	return technicalData, nil
}

// MarshalTechnicalDataGen1 marshals Gen1 Technical Data using raw data painting.
func (opts MarshalOptions) MarshalTechnicalDataGen1(technicalData *vuv1.TechnicalDataGen1) ([]byte, error) {
	if technicalData == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	raw := technicalData.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Technical Data Gen1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeTechnicalDataGen1 anonymizes Gen1 Technical Data.
// TODO: Implement full semantic anonymization (anonymize VIN, VRN, sensor IDs, etc.).
func (opts AnonymizeOptions) anonymizeTechnicalDataGen1(td *vuv1.TechnicalDataGen1) *vuv1.TechnicalDataGen1 {
	if td == nil {
		return nil
	}
	result := proto.Clone(td).(*vuv1.TechnicalDataGen1)
	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Note: We intentionally keep raw_data here because MarshalTechnicalDataGen1
	// currently requires raw_data (semantic marshalling not yet implemented).
	// Once semantic marshalling is implemented, we should clear raw_data and
	// implement full semantic anonymization of VIN, VRN, sensor data, etc.

	return result
}
