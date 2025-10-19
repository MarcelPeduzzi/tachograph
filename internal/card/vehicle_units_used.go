package card

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalVehicleUnitsUsed unmarshals vehicle units used data from a card EF.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
//
//	CardVehicleUnitsUsed ::= SEQUENCE {
//	    vehicleUnitPointerNewestRecord     INTEGER(0..NoOfCardVehicleUnitRecords-1),
//	    cardVehicleUnitRecords             SET SIZE(NoOfCardVehicleUnitRecords) OF CardVehicleUnitRecord
//	}
//
//	CardVehicleUnitRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    manufacturerCode                   ManufacturerCode,     -- 1 byte
//	    deviceID                           OCTET STRING(SIZE(1)), -- 1 byte
//	    vuSoftwareVersion                  VuSoftwareVersion     -- 4 bytes
//	}
//
// Binary structure:
//   - 2 bytes: vehicleUnitPointerNewestRecord
//   - N * 10 bytes: fixed-size array of CardVehicleUnitRecord (N determined by data length)
//
// Typical sizes:
//   - Driver cards: 2002 bytes (200 records)
//   - Control cards: 82 bytes (8 records)
func (opts UnmarshalOptions) unmarshalVehicleUnitsUsed(data []byte) (*cardv1.VehicleUnitsUsed, error) {
	const (
		idxNewestRecordPointer         = 0
		lenNewestRecordPointer         = 2
		lenCardVehicleUnitRecord       = 10
		lenCardVehicleUnitsUsedMinimum = lenNewestRecordPointer
	)

	if len(data) < lenCardVehicleUnitsUsedMinimum {
		return nil, fmt.Errorf("invalid data length for CardVehicleUnitsUsed: got %d, want at least %d", len(data), lenCardVehicleUnitsUsedMinimum)
	}

	// Validate that the records section is a multiple of record size
	recordsDataLen := len(data) - lenNewestRecordPointer
	if recordsDataLen%lenCardVehicleUnitRecord != 0 {
		return nil, fmt.Errorf("invalid records data length for CardVehicleUnitsUsed: got %d bytes, not a multiple of %d", recordsDataLen, lenCardVehicleUnitRecord)
	}

	var target cardv1.VehicleUnitsUsed

	// Parse newest record pointer
	newestRecordPointer := binary.BigEndian.Uint16(data[idxNewestRecordPointer:])
	target.SetVehicleUnitPointerNewestRecord(int32(newestRecordPointer))

	// Parse records using bufio.Scanner pattern
	records, err := opts.unmarshalCardVehicleUnitRecords(data[lenNewestRecordPointer:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle unit records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// splitCardVehicleUnitRecord is a bufio.SplitFunc for parsing CardVehicleUnitRecord entries.
func splitCardVehicleUnitRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const lenCardVehicleUnitRecord = 10

	if len(data) < lenCardVehicleUnitRecord {
		if atEOF {
			return 0, nil, nil
		}
		return 0, nil, nil
	}

	return lenCardVehicleUnitRecord, data[:lenCardVehicleUnitRecord], nil
}

// unmarshalCardVehicleUnitRecords parses the fixed-size array of vehicle unit records.
func (opts UnmarshalOptions) unmarshalCardVehicleUnitRecords(data []byte) ([]*cardv1.VehicleUnitsUsed_Record, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(splitCardVehicleUnitRecord)

	var records []*cardv1.VehicleUnitsUsed_Record
	for scanner.Scan() {
		record, err := opts.unmarshalCardVehicleUnitRecord(scanner.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal vehicle unit record: %w", err)
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return records, nil
}

// unmarshalCardVehicleUnitRecord unmarshals a single CardVehicleUnitRecord.
//
// Binary structure (10 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 1 byte: ManufacturerCode
//   - 1 byte: deviceID
//   - 4 bytes: VuSoftwareVersion (IA5String)
func (opts UnmarshalOptions) unmarshalCardVehicleUnitRecord(data []byte) (*cardv1.VehicleUnitsUsed_Record, error) {
	const (
		idxTimeStamp         = 0
		idxManufacturerCode  = 4
		idxDeviceID          = 5
		idxVuSoftwareVersion = 6
		lenRecord            = 10
	)

	if len(data) != lenRecord {
		return nil, fmt.Errorf("invalid data length for CardVehicleUnitRecord: got %d, want %d", len(data), lenRecord)
	}

	var record cardv1.VehicleUnitsUsed_Record

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimeStamp : idxTimeStamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse manufacturer code (1 byte)
	record.SetManufacturerCode(int32(data[idxManufacturerCode]))

	// Parse device ID (1 byte)
	record.SetDeviceId(data[idxDeviceID : idxDeviceID+1])

	// Parse VU software version (4 bytes, IA5String)
	record.SetVuSoftwareVersion(data[idxVuSoftwareVersion : idxVuSoftwareVersion+4])

	return &record, nil
}

// MarshalCardVehicleUnitsUsed marshals vehicle units used data.
//
// The data type `CardVehicleUnitsUsed` is specified in the Data Dictionary, Section 2.40.
//
// ASN.1 Definition:
//
//	CardVehicleUnitsUsed ::= SEQUENCE {
//	    vehicleUnitPointerNewestRecord     INTEGER(0..NoOfCardVehicleUnitRecords-1),
//	    cardVehicleUnitRecords             SET SIZE(NoOfCardVehicleUnitRecords) OF CardVehicleUnitRecord
//	}
//
//	CardVehicleUnitRecord ::= SEQUENCE {
//	    timeStamp                          TimeReal,             -- 4 bytes
//	    manufacturerCode                   ManufacturerCode,     -- 1 byte
//	    deviceID                           OCTET STRING(SIZE(1)), -- 1 byte
//	    vuSoftwareVersion                  VuSoftwareVersion     -- 4 bytes
//	}
//
// Binary structure:
//   - 2 bytes: vehicleUnitPointerNewestRecord
//   - N * 10 bytes: fixed-size array of CardVehicleUnitRecord
//
// The number of records (N) is determined from the original data, not explicitly stored.
func (opts MarshalOptions) MarshalCardVehicleUnitsUsed(vehicleUnits *cardv1.VehicleUnitsUsed) ([]byte, error) {
	if vehicleUnits == nil {
		return nil, nil
	}

	var dst []byte

	// Append newest record pointer (2 bytes)
	newestRecordPointer := vehicleUnits.GetVehicleUnitPointerNewestRecord()
	dst = binary.BigEndian.AppendUint16(dst, uint16(newestRecordPointer))

	// Append all vehicle unit records
	// The binary format is a fixed-size array, so we write exactly what we have
	records := vehicleUnits.GetRecords()
	for _, record := range records {
		recordBytes, err := opts.MarshalCardVehicleUnitRecord(record)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vehicle unit record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}

// MarshalCardVehicleUnitRecord marshals a single vehicle unit record.
//
// Binary structure (10 bytes):
//   - 4 bytes: TimeReal timestamp
//   - 1 byte: ManufacturerCode
//   - 1 byte: deviceID
//   - 4 bytes: VuSoftwareVersion (IA5String)
func (opts MarshalOptions) MarshalCardVehicleUnitRecord(record *cardv1.VehicleUnitsUsed_Record) ([]byte, error) {
	if record == nil {
		return nil, nil
	}

	var dst []byte

	// Append timestamp (TimeReal - 4 bytes)

	timestampBytes, err := opts.MarshalTimeReal(record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal timestamp: %w", err)
	}
	dst = append(dst, timestampBytes...)

	// Append manufacturer code (1 byte)
	manufacturerCode := record.GetManufacturerCode()
	if manufacturerCode < 0 || manufacturerCode > 255 {
		return nil, fmt.Errorf("invalid manufacturer code: %d", manufacturerCode)
	}
	dst = append(dst, byte(manufacturerCode))

	// Append device ID (1 byte)
	deviceID := record.GetDeviceId()
	if len(deviceID) > 1 {
		return nil, fmt.Errorf("device ID too long: %d bytes", len(deviceID))
	}
	if len(deviceID) == 1 {
		dst = append(dst, deviceID[0])
	} else {
		dst = append(dst, 0x00)
	}

	// Append VU software version (4 bytes)
	vuSoftwareVersion := record.GetVuSoftwareVersion()
	if len(vuSoftwareVersion) > 4 {
		return nil, fmt.Errorf("VU software version too long: %d bytes", len(vuSoftwareVersion))
	}
	if len(vuSoftwareVersion) == 4 {
		dst = append(dst, vuSoftwareVersion...)
	} else {
		// Pad with zeros if shorter than 4 bytes
		padded := make([]byte, 4)
		copy(padded, vuSoftwareVersion)
		dst = append(dst, padded...)
	}

	return dst, nil
}

// anonymizeVehicleUnitsUsed creates an anonymized copy of VehicleUnitsUsed,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeVehicleUnitsUsed(vu *cardv1.VehicleUnitsUsed) *cardv1.VehicleUnitsUsed {
	if vu == nil {
		return nil
	}

	result := &cardv1.VehicleUnitsUsed{}

	// Preserve the pointer to newest record
	result.SetVehicleUnitPointerNewestRecord(vu.GetVehicleUnitPointerNewestRecord())

	// Anonymize each record
	originalRecords := vu.GetRecords()
	anonymizedRecords := make([]*cardv1.VehicleUnitsUsed_Record, len(originalRecords))
	for i, record := range originalRecords {
		anonymizedRecords[i] = opts.anonymizeVehicleUnitRecord(record, i)
	}
	result.SetRecords(anonymizedRecords)

	// Preserve signature if present
	if vu.HasSignature() {
		result.SetSignature(vu.GetSignature())
	}

	return result
}

// anonymizeVehicleUnitRecord anonymizes a single vehicle unit record.
// Uses index to create sequential timestamps.
func (opts AnonymizeOptions) anonymizeVehicleUnitRecord(record *cardv1.VehicleUnitsUsed_Record, index int) *cardv1.VehicleUnitsUsed_Record {
	if record == nil {
		// Return a zero-filled record for nil entries
		zeroRecord := &cardv1.VehicleUnitsUsed_Record{}
		zeroRecord.SetManufacturerCode(0)
		zeroRecord.SetDeviceId([]byte{0})
		zeroRecord.SetVuSoftwareVersion([]byte{0, 0, 0, 0})
		return zeroRecord
	}

	result := &cardv1.VehicleUnitsUsed_Record{}

	// Replace timestamp with sequential test timestamps
	// Base: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	// Increment by 1 hour per record
	baseEpoch := int64(1577836800)
	ts := record.GetTimestamp()
	if ts != nil && ts.Seconds != 0 {
		// Non-zero timestamp - replace with sequential test timestamp
		result.SetTimestamp(&timestamppb.Timestamp{
			Seconds: baseEpoch + int64(index)*3600,
			Nanos:   0,
		})
	}
	// else: Zero timestamp - leave unset (nil)

	// Replace manufacturer code with test value 0x40 for non-zero values
	if record.GetManufacturerCode() != 0 {
		result.SetManufacturerCode(0x40)
	} else {
		result.SetManufacturerCode(0)
	}

	// Replace device ID with test value 0x00
	if len(record.GetDeviceId()) > 0 {
		result.SetDeviceId([]byte{0x00})
	}

	// Replace VU software version with "0000" for non-zero values
	swVersion := record.GetVuSoftwareVersion()
	if len(swVersion) > 0 && !bytes.Equal(swVersion, []byte{0, 0, 0, 0}) {
		result.SetVuSoftwareVersion([]byte("0000"))
	} else {
		result.SetVuSoftwareVersion([]byte{0, 0, 0, 0})
	}

	return result
}
