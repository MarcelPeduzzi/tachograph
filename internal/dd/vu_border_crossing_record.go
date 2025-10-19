package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuBorderCrossingRecord parses a VuBorderCrossingRecord (57 bytes).
//
// The data type `VuBorderCrossingRecord` is specified in the Data Dictionary, Section 2.203a.
//
// ASN.1 Definition (Gen2 V2):
//
//	VuBorderCrossingRecord ::= SEQUENCE {
//	    cardNumberAndGenDriverSlot      FullCardNumberAndGeneration,
//	    cardNumberAndGenCodriverSlot    FullCardNumberAndGeneration,
//	    countryLeft                     NationNumeric,
//	    countryEntered                  NationNumeric,
//	    gnssPlaceAuthRecord             GNSSPlaceAuthRecord,
//	    vehicleOdometerValue            OdometerShort
//	}
//
// Binary Layout (fixed length, 57 bytes):
//   - Bytes 0-19: cardNumberAndGenDriverSlot (FullCardNumberAndGeneration)
//   - Bytes 20-39: cardNumberAndGenCodriverSlot (FullCardNumberAndGeneration)
//   - Byte 40: countryLeft (NationNumeric)
//   - Byte 41: countryEntered (NationNumeric)
//   - Bytes 42-53: gnssPlaceAuthRecord (GNSSPlaceAuthRecord)
//   - Bytes 54-56: vehicleOdometerValue (OdometerShort)
func (opts UnmarshalOptions) UnmarshalVuBorderCrossingRecord(data []byte) (*ddv1.VuBorderCrossingRecord, error) {
	const (
		idxCardNumberDriverSlot   = 0
		idxCardNumberCodriverSlot = 20
		idxCountryLeft            = 40
		idxCountryEntered         = 41
		idxGnssPlaceAuthRecord    = 42
		idxVehicleOdometerValue   = 54
		lenVuBorderCrossingRecord = 57

		lenFullCardNumberAndGeneration = 20
		lenNationNumeric               = 1
		lenGNSSPlaceAuthRecord         = 12
		lenOdometerShort               = 3
	)

	if len(data) != lenVuBorderCrossingRecord {
		return nil, fmt.Errorf("invalid data length for VuBorderCrossingRecord: got %d, want %d", len(data), lenVuBorderCrossingRecord)
	}

	record := &ddv1.VuBorderCrossingRecord{}
	if opts.PreserveRawData {
		record.SetRawData(data)
	}

	// cardNumberAndGenDriverSlot (20 bytes)
	cardNumberDriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberDriverSlot : idxCardNumberDriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number driver slot: %w", err)
	}
	record.SetCardNumberDriverSlot(cardNumberDriverSlot)

	// cardNumberAndGenCodriverSlot (20 bytes)
	cardNumberCodriverSlot, err := opts.UnmarshalFullCardNumberAndGeneration(data[idxCardNumberCodriverSlot : idxCardNumberCodriverSlot+lenFullCardNumberAndGeneration])
	if err != nil {
		return nil, fmt.Errorf("unmarshal card number codriver slot: %w", err)
	}
	record.SetCardNumberCodriverSlot(cardNumberCodriverSlot)

	// countryLeft (1 byte)
	countryLeft, err := UnmarshalEnum[ddv1.NationNumeric](data[idxCountryLeft])
	if err != nil {
		return nil, fmt.Errorf("unmarshal country left: %w", err)
	}
	record.SetCountryLeft(countryLeft)

	// countryEntered (1 byte)
	countryEntered, err := UnmarshalEnum[ddv1.NationNumeric](data[idxCountryEntered])
	if err != nil {
		return nil, fmt.Errorf("unmarshal country entered: %w", err)
	}
	record.SetCountryEntered(countryEntered)

	// gnssPlaceAuthRecord (12 bytes)
	gnssPlaceAuthRecord, err := opts.UnmarshalGNSSPlaceAuthRecord(data[idxGnssPlaceAuthRecord : idxGnssPlaceAuthRecord+lenGNSSPlaceAuthRecord])
	if err != nil {
		return nil, fmt.Errorf("unmarshal GNSS place auth record: %w", err)
	}
	record.SetGnssPlaceAuthRecord(gnssPlaceAuthRecord)

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerValue, err := opts.UnmarshalOdometer(data[idxVehicleOdometerValue : idxVehicleOdometerValue+lenOdometerShort])
	if err != nil {
		return nil, fmt.Errorf("unmarshal vehicle odometer value: %w", err)
	}
	record.SetVehicleOdometerKm(int32(vehicleOdometerValue))

	return record, nil
}

// MarshalVuBorderCrossingRecord marshals a VuBorderCrossingRecord (57 bytes) to bytes.
func (opts MarshalOptions) MarshalVuBorderCrossingRecord(record *ddv1.VuBorderCrossingRecord) ([]byte, error) {
	if record == nil {
		return nil, fmt.Errorf("record cannot be nil")
	}

	const lenVuBorderCrossingRecord = 57

	// Use raw data painting strategy if available
	var canvas [lenVuBorderCrossingRecord]byte
	if rawData := record.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenVuBorderCrossingRecord {
			return nil, fmt.Errorf("invalid raw_data length for VuBorderCrossingRecord: got %d, want %d", len(rawData), lenVuBorderCrossingRecord)
		}
		copy(canvas[:], rawData)
	}

	offset := 0

	// cardNumberAndGenDriverSlot (20 bytes)
	cardNumberDriverSlotBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetCardNumberDriverSlot())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card number driver slot: %w", err)
	}
	copy(canvas[offset:offset+20], cardNumberDriverSlotBytes)
	offset += 20

	// cardNumberAndGenCodriverSlot (20 bytes)
	cardNumberCodriverSlotBytes, err := opts.MarshalFullCardNumberAndGeneration(record.GetCardNumberCodriverSlot())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card number codriver slot: %w", err)
	}
	copy(canvas[offset:offset+20], cardNumberCodriverSlotBytes)
	offset += 20

	// countryLeft (1 byte)
	countryLeftByte, _ := MarshalEnum(record.GetCountryLeft())
	canvas[offset] = countryLeftByte
	offset += 1

	// countryEntered (1 byte)
	countryEnteredByte, _ := MarshalEnum(record.GetCountryEntered())
	canvas[offset] = countryEnteredByte
	offset += 1

	// gnssPlaceAuthRecord (12 bytes)
	gnssPlaceAuthRecordBytes, err := opts.MarshalGNSSPlaceAuthRecord(record.GetGnssPlaceAuthRecord())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GNSS place auth record: %w", err)
	}
	copy(canvas[offset:offset+12], gnssPlaceAuthRecordBytes)
	offset += 12

	// vehicleOdometerValue (3 bytes)
	vehicleOdometerBytes, err := opts.MarshalOdometer(record.GetVehicleOdometerKm())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle odometer value: %w", err)
	}
	copy(canvas[offset:offset+3], vehicleOdometerBytes)

	return canvas[:], nil
}
