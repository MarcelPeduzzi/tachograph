package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalVehiclesUsedG2 unmarshals vehicles used data from a Gen2 card EF.
func (opts UnmarshalOptions) unmarshalVehiclesUsedG2(data []byte) (*cardv1.VehiclesUsedG2, error) {
	const (
		lenMinEfVehiclesUsed = 2 // Minimum EF_Vehicles_Used record size
	)

	if len(data) < lenMinEfVehiclesUsed {
		return nil, fmt.Errorf("insufficient data for vehicles used: got %d bytes, need at least %d", len(data), lenMinEfVehiclesUsed)
	}

	var target cardv1.VehiclesUsedG2

	// Save complete raw data for painting
	target.SetRawData(data)

	r := bytes.NewReader(data)

	// Read newest record pointer (2 bytes)
	var newestRecordIndex uint16
	if err := binary.Read(r, binary.BigEndian, &newestRecordIndex); err != nil {
		return nil, fmt.Errorf("failed to read newest record index: %w", err)
	}

	target.SetNewestRecordIndex(int32(newestRecordIndex))

	// Parse Gen2 vehicle records (48 bytes each)
	records, err := opts.unmarshalVehicleRecordsGen2(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gen2 vehicle records: %w", err)
	}
	target.SetRecords(records)

	return &target, nil
}

// unmarshalVehicleRecordsGen2 parses Gen2 vehicle records (48 bytes each).
func (opts UnmarshalOptions) unmarshalVehicleRecordsGen2(r *bytes.Reader) ([]*ddv1.CardVehicleRecordG2, error) {
	const lenCardVehicleRecord = 48

	var records []*ddv1.CardVehicleRecordG2
	for r.Len() >= lenCardVehicleRecord {
		recordBytes := make([]byte, lenCardVehicleRecord)
		if _, err := r.Read(recordBytes); err != nil {
			break // Stop parsing on error, but return what we have
		}

		record, err := opts.UnmarshalOptions.UnmarshalCardVehicleRecordG2(recordBytes)
		if err != nil {
			return records, fmt.Errorf("failed to parse Gen2 vehicle record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// MarshalVehiclesUsedG2 marshals Gen2 vehicles used data.
func (opts MarshalOptions) MarshalVehiclesUsedG2(vehiclesUsed *cardv1.VehiclesUsedG2) ([]byte, error) {
	if vehiclesUsed == nil {
		return nil, nil
	}

	// Calculate expected size: 2 bytes (pointer) + N records × 48 bytes
	const recordSize = 48
	numRecords := len(vehiclesUsed.GetRecords())
	expectedSize := 2 + (numRecords * recordSize)

	// Use raw_data as canvas if available and correct size
	if rawData := vehiclesUsed.GetRawData(); len(rawData) == expectedSize {
		// Make a copy to use as canvas
		canvas := make([]byte, expectedSize)
		copy(canvas, rawData)

		// Paint newest record index over canvas
		binary.BigEndian.PutUint16(canvas[0:2], uint16(vehiclesUsed.GetNewestRecordIndex()))

		// Paint each record over canvas

		offset := 2
		for _, record := range vehiclesUsed.GetRecords() {
			recordBytes, err := opts.MarshalCardVehicleRecordG2(record)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Gen2 vehicle record: %w", err)
			}
			if len(recordBytes) != recordSize {
				return nil, fmt.Errorf("invalid Gen2 vehicle record size: got %d, want %d", len(recordBytes), recordSize)
			}
			copy(canvas[offset:offset+recordSize], recordBytes)
			offset += recordSize
		}

		return canvas, nil
	}

	// Fall back to building from scratch
	var dst []byte
	newestRecordIndex := uint16(vehiclesUsed.GetNewestRecordIndex())
	dst = binary.BigEndian.AppendUint16(dst, newestRecordIndex)

	for _, record := range vehiclesUsed.GetRecords() {
		recordBytes, err := opts.MarshalCardVehicleRecordG2(record)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Gen2 vehicle record: %w", err)
		}
		dst = append(dst, recordBytes...)
	}

	return dst, nil
}
