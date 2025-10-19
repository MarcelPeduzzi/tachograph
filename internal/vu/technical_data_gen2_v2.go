package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalTechnicalDataGen2V2 parses Gen2 V2 Technical Data from the complete transfer value.
//
// Gen2 V2 Technical Data structure is identical to Gen2 V1.
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
func unmarshalTechnicalDataGen2V2(value []byte) (*vuv1.TechnicalDataGen2V2, error) {
	technicalData := &vuv1.TechnicalDataGen2V2{}
	technicalData.SetRawData(value)

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
	if err := skipRecordArray("VuApprovalNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuSoftwareIdentification"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuManufacturerName"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuManufacturerAddress"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuPartNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("VuSerialNumber"); err != nil {
		return nil, err
	}
	if err := skipRecordArray("SensorPaired"); err != nil {
		return nil, err
	}

	// SignatureRecordArray is now handled separately in raw parsing, not part of value

	if offset != len(value) {
		return nil, fmt.Errorf("Technical Data Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	return technicalData, nil
}

// MarshalTechnicalDataGen2V2 marshals Gen2 V2 Technical Data using raw data painting.
func (opts MarshalOptions) MarshalTechnicalDataGen2V2(technicalData *vuv1.TechnicalDataGen2V2) ([]byte, error) {
	if technicalData == nil {
		return nil, fmt.Errorf("technicalData cannot be nil")
	}

	raw := technicalData.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	return nil, fmt.Errorf("cannot marshal Technical Data Gen2 V2 without raw_data (semantic marshalling not yet implemented)")
}


// anonymizeTechnicalDataGen2V2 anonymizes Gen2 V2 Technical Data.
// TODO: Implement full anonymization logic for Gen2 V2 technical data.
func (opts AnonymizeOptions) anonymizeTechnicalDataGen2V2(td *vuv1.TechnicalDataGen2V2) *vuv1.TechnicalDataGen2V2 {
	if td == nil {
		return nil
	}
	result := proto.Clone(td).(*vuv1.TechnicalDataGen2V2)
	result.SetSignature(nil)
	result.SetRawData(nil)
	return result
}
