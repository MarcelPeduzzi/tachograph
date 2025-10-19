package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalCardLoadTypeEntryRecord parses a CardLoadTypeEntryRecord (5 bytes).
//
// The data type `CardLoadTypeEntryRecord` is specified in the Data Dictionary, Section 2.24b.
//
// ASN.1 Definition (Gen2 V2):
//
//	CardLoadTypeEntryRecord ::= SEQUENCE {
//	    timeStamp                       TimeReal,
//	    loadTypeEntered                 LoadType
//	}
//
// Binary Layout (fixed length, 5 bytes):
//   - Bytes 0-3: timeStamp (TimeReal)
//   - Byte 4: loadTypeEntered (LoadType)
func (opts UnmarshalOptions) UnmarshalCardLoadTypeEntryRecord(data []byte) (*ddv1.CardLoadTypeEntryRecord, error) {
	const (
		idxTimeStamp               = 0
		idxLoadTypeEntered         = 4
		lenCardLoadTypeEntryRecord = 5

		lenTimeReal = 4
		lenLoadType = 1
	)

	if len(data) != lenCardLoadTypeEntryRecord {
		return nil, fmt.Errorf("invalid data length for CardLoadTypeEntryRecord: got %d, want %d", len(data), lenCardLoadTypeEntryRecord)
	}

	record := &ddv1.CardLoadTypeEntryRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// timeStamp (4 bytes)
	timeStamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+lenTimeReal])
	if err != nil {
		return nil, fmt.Errorf("unmarshal time stamp: %w", err)
	}
	record.SetTimeStamp(timeStamp)

	// loadTypeEntered (1 byte)
	loadTypeEntered, err := UnmarshalEnum[ddv1.LoadType](data[idxLoadTypeEntered])
	if err != nil {
		return nil, fmt.Errorf("unmarshal load type entered: %w", err)
	}
	record.SetLoadTypeEntered(loadTypeEntered)

	return record, nil
}

// MarshalCardLoadTypeEntryRecord marshals a CardLoadTypeEntryRecord (5 bytes) to bytes.
func (opts MarshalOptions) MarshalCardLoadTypeEntryRecord(record *ddv1.CardLoadTypeEntryRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenCardLoadTypeEntryRecord = 5

	// Use raw data painting strategy if available
	var canvas [lenCardLoadTypeEntryRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenCardLoadTypeEntryRecord {
			return nil, fmt.Errorf("invalid raw_data length for CardLoadTypeEntryRecord: got %d, want %d", len(rawData), lenCardLoadTypeEntryRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// timeStamp (4 bytes)
	timeStampBytes, err := opts.MarshalTimeReal(record.GetTimeStamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal time stamp: %w", err)
	}
	copy(canvas[offset:offset+4], timeStampBytes)
	offset += 4

	// loadTypeEntered (1 byte)
	loadTypeEnteredByte, _ := MarshalEnum(record.GetLoadTypeEntered())
	canvas[offset] = loadTypeEnteredByte

	return canvas[:], nil
}
