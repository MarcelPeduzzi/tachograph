package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuPlaceDailyWorkPeriodRecord parses a Generation 1 VuPlaceDailyWorkPeriodRecord (28 bytes).
//
// The data type `VuPlaceDailyWorkPeriodRecord` is specified in the Data Dictionary, Section 2.219.
//
// ASN.1 Definition (Gen1):
//
//	VuPlaceDailyWorkPeriodRecord ::= SEQUENCE {
//	    fullCardNumber              FullCardNumber,
//	    placeRecord                 PlaceRecord
//	}
//
// Binary Layout (fixed length, 28 bytes):
//   - Bytes 0-17: fullCardNumber (FullCardNumber)
//   - Bytes 18-27: placeRecord (PlaceRecord)
func (opts UnmarshalOptions) UnmarshalVuPlaceDailyWorkPeriodRecord(data []byte) (*ddv1.VuPlaceDailyWorkPeriodRecord, error) {
	const (
		idxFullCardNumber               = 0
		idxPlaceRecord                  = 18
		lenVuPlaceDailyWorkPeriodRecord = 28

		lenFullCardNumber = 18
		lenPlaceRecord    = 10
	)

	if len(data) != lenVuPlaceDailyWorkPeriodRecord {
		return nil, fmt.Errorf("invalid data length for VuPlaceDailyWorkPeriodRecord: got %d, want %d", len(data), lenVuPlaceDailyWorkPeriodRecord)
	}

	record := &ddv1.VuPlaceDailyWorkPeriodRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// fullCardNumber (18 bytes)
	fullCardNumber, err := opts.UnmarshalFullCardNumber(data[idxFullCardNumber : idxFullCardNumber+lenFullCardNumber])
	if err != nil {
		return nil, fmt.Errorf("unmarshal full card number: %w", err)
	}
	record.SetFullCardNumber(fullCardNumber)

	// placeRecord (10 bytes)
	placeRecord, err := opts.UnmarshalPlaceRecord(data[idxPlaceRecord : idxPlaceRecord+lenPlaceRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal place record: %w", err)
	}
	record.SetPlaceRecord(placeRecord)

	return record, nil
}

// MarshalVuPlaceDailyWorkPeriodRecord marshals a VuPlaceDailyWorkPeriodRecord (28 bytes) to bytes.
func (opts MarshalOptions) MarshalVuPlaceDailyWorkPeriodRecord(record *ddv1.VuPlaceDailyWorkPeriodRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuPlaceDailyWorkPeriodRecord = 28

	// Use raw data painting strategy if available
	var canvas [lenVuPlaceDailyWorkPeriodRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenVuPlaceDailyWorkPeriodRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuPlaceDailyWorkPeriodRecord: got %d, want %d", len(rawData), lenVuPlaceDailyWorkPeriodRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// fullCardNumber (18 bytes)
	fullCardNumberBytes, err := opts.MarshalFullCardNumber(record.GetFullCardNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal full card number: %w", err)
	}
	copy(canvas[offset:offset+18], fullCardNumberBytes)
	offset += 18

	// placeRecord (10 bytes)
	placeRecordBytes, err := opts.MarshalPlaceRecord(record.GetPlaceRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal place record: %w", err)
	}
	copy(canvas[offset:offset+10], placeRecordBytes)

	return canvas[:], nil
}
