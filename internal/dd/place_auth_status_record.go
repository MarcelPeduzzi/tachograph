package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPlaceAuthStatusRecord parses a PlaceAuthStatusRecord (5 bytes).
//
// The data type `PlaceAuthStatusRecord` is specified in the Data Dictionary, Section 2.116b.
//
// ASN.1 Definition (Gen2 V2):
//
//	PlaceAuthStatusRecord ::= SEQUENCE {
//	    entryTime                       TimeReal,
//	    authenticationStatus            PositionAuthenticationStatus
//	}
//
// Binary Layout (fixed length, 5 bytes):
//   - Bytes 0-3: entryTime (TimeReal)
//   - Byte 4: authenticationStatus (PositionAuthenticationStatus)
func (opts UnmarshalOptions) UnmarshalPlaceAuthStatusRecord(data []byte) (*ddv1.PlaceAuthStatusRecord, error) {
	const (
		idxEntryTime             = 0
		idxAuthenticationStatus  = 4
		lenPlaceAuthStatusRecord = 5

		lenTimeReal                     = 4
		lenPositionAuthenticationStatus = 1
	)

	if len(data) != lenPlaceAuthStatusRecord {
		return nil, fmt.Errorf("invalid data length for PlaceAuthStatusRecord: got %d, want %d", len(data), lenPlaceAuthStatusRecord)
	}

	record := &ddv1.PlaceAuthStatusRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// entryTime (4 bytes)
	entryTime, err := opts.UnmarshalTimeReal(data[idxEntryTime : idxEntryTime+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal entry time: %w", err)
	}
	record.SetEntryTime(entryTime)

	// authenticationStatus (1 byte)
	authenticationStatus, err := UnmarshalEnum[ddv1.PositionAuthenticationStatus](data[idxAuthenticationStatus])
	if err != nil {
		return nil, fmt.Errorf("unmarshal authentication status: %w", err)
	}
	record.SetAuthenticationStatus(authenticationStatus)

	return record, nil
}

// MarshalPlaceAuthStatusRecord marshals a PlaceAuthStatusRecord (5 bytes) to bytes.
func (opts MarshalOptions) MarshalPlaceAuthStatusRecord(record *ddv1.PlaceAuthStatusRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenPlaceAuthStatusRecord = 5

	// Use raw data painting strategy if available
	var canvas [lenPlaceAuthStatusRecord]byte
	if record.HasRawData() {
		rawData := record.GetRawData()
		if len(rawData) != lenPlaceAuthStatusRecord {
			return nil, fmt.Errorf("invalid raw_data length for PlaceAuthStatusRecord: got %d, want %d", len(rawData), lenPlaceAuthStatusRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// entryTime (4 bytes)
	entryTimeBytes, err := opts.MarshalTimeReal(record.GetEntryTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry time: %w", err)
	}
	copy(canvas[offset:offset+4], entryTimeBytes)
	offset += 4

	// authenticationStatus (1 byte)
	authenticationStatusByte, _ := MarshalEnum(record.GetAuthenticationStatus())
	canvas[offset] = authenticationStatusByte

	return canvas[:], nil
}
