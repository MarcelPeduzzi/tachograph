package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuPlaceDailyWorkPeriodRecordG2 parses a Generation 2 version 1 VuPlaceDailyWorkPeriodRecord (41 bytes).
//
// The data type `VuPlaceDailyWorkPeriodRecord` is specified in the Data Dictionary, Section 2.219.
//
// ASN.1 Definition (Gen2 V1):
//
//	VuPlaceDailyWorkPeriodRecord ::= SEQUENCE {
//	    fullCardNumberAndGeneration FullCardNumberAndGeneration,
//	    placeRecord                 PlaceRecord
//	}
//
// Binary Layout (fixed length, 41 bytes):
//   - Bytes 0-19: fullCardNumberAndGeneration (FullCardNumberAndGeneration)
//   - Bytes 20-40: placeRecord (PlaceRecordG2)
func (opts UnmarshalOptions) UnmarshalVuPlaceDailyWorkPeriodRecordG2(data []byte) (*ddv1.VuPlaceDailyWorkPeriodRecordG2, error) {
	const (
		idxFullCardNumber                 = 0
		idxPlaceRecord                    = 20
		lenVuPlaceDailyWorkPeriodRecordG2 = 41

		lenFullCardNumberAndGeneration = 20
		lenPlaceRecordG2               = 21
	)

	if len(data) != lenVuPlaceDailyWorkPeriodRecordG2 {
		return nil, fmt.Errorf("invalid data length for VuPlaceDailyWorkPeriodRecordG2: got %d, want %d", len(data), lenVuPlaceDailyWorkPeriodRecordG2)
	}

	record := &ddv1.VuPlaceDailyWorkPeriodRecordG2{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// fullCardNumberAndGeneration (20 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxFullCardNumber : idxFullCardNumber+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number and generation: %w", err)
	}
	record.SetFullCardNumber(fullCardNumber)

	// placeRecord (21 bytes)
	placeRecord, err := opts.UnmarshalPlaceRecordG2(data[idxPlaceRecord : idxPlaceRecord+lenPlaceRecordG2])
	if err != nil {
		return nil, fmt.Errorf("unmarshal place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	return record, nil
}

// MarshalVuPlaceDailyWorkPeriodRecordG2 marshals a VuPlaceDailyWorkPeriodRecordG2 (41 bytes) to bytes.
func (opts MarshalOptions) MarshalVuPlaceDailyWorkPeriodRecordG2(record *ddv1.VuPlaceDailyWorkPeriodRecordG2) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuPlaceDailyWorkPeriodRecordG2 = 41

	// Use raw data painting strategy if available
	var canvas [lenVuPlaceDailyWorkPeriodRecordG2]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenVuPlaceDailyWorkPeriodRecordG2 {
			return nil, fmt.Errorf("invalid raw_data length for VuPlaceDailyWorkPeriodRecordG2: got %d, want %d", len(rawData), lenVuPlaceDailyWorkPeriodRecordG2)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// fullCardNumberAndGeneration (20 bytes)
	fullCardNumberBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetFullCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal full card number and generation: %w", err)
	}
	copy(canvas[offset:offset+20], fullCardNumberBytes)
	offset += 20

	// placeRecord (21 bytes)
	placeRecordBytes, err := opts.MarshalPlaceRecordG2(record.GetPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal place record: %w", err)
	}
	copy(canvas[offset:offset+21], placeRecordBytes)

	return canvas[:], nil
}
