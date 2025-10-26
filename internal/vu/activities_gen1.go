package vu

import (
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalActivitiesGen1 parses Gen1 Activities data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
//
// Gen1 Activities structure (from Data Dictionary and Appendix 7, Section 2.2.6.3):
//
// ASN.1 Definition:
//
//	VuActivitiesFirstGen ::= SEQUENCE {
//	    timeReal                      TimeReal,                              -- 4 bytes
//	    odometerValueMidnight         OdometerShort,                         -- 3 bytes
//	    vuCardIWData                  VuCardIWDataFirstGen,                  -- 2 + (N * 129) bytes
//	    vuActivityDailyData           VuActivityDailyDataFirstGen,           -- 2 + (M * 2) bytes
//	    vuPlaceDailyWorkPeriodData    VuPlaceDailyWorkPeriodDataFirstGen,    -- 1 + (P * 28) bytes
//	    vuSpecificConditionData       VuSpecificConditionDataFirstGen,       -- 2 + (Q * 5) bytes
//	    signature                     SignatureFirstGen                      -- 128 bytes (RSA)
//	}
//
// Binary Layout:
// - TimeReal: 4 bytes (date of day downloaded)
// - OdometerValueMidnight: 3 bytes (OdometerShort)
// - VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
//   - Each VuCardIWRecordFirstGen: 129 bytes
//   - FullCardNumber: 18 bytes
//   - ManufacturerCode: 1 byte
//   - DownloadTime: 4 bytes
//   - ... (rest of record)
//
// - VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
//   - Each ActivityChangeInfo: 2 bytes
//
// - VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
//   - Each VuPlaceDailyWorkPeriodRecordFirstGen: 28 bytes
//   - FullCardNumber: 18 bytes
//   - PlaceRecord: 10 bytes
//
// - VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
//   - Each SpecificConditionRecord: 5 bytes
//   - TimeReal: 4 bytes
//   - SpecificConditionType: 1 byte
//
// - Signature: 128 bytes (RSA-1024)
//
// Note: This is a minimal implementation that validates the binary structure and stores raw_data.
// Full semantic parsing of all nested records is TODO.
func unmarshalActivitiesGen1(value []byte) (*vuv1.ActivitiesGen1, error) {
	// Split transfer value into data and signature
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	const signatureSize = 128
	if len(value) < signatureSize {
		return nil, fmt.Errorf("insufficient data for signature: need at least %d bytes, got %d", signatureSize, len(value))
	}

	dataSize := len(value) - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	activities := &vuv1.ActivitiesGen1{}
	activities.SetRawData(value) // Store complete transfer value for painting

	offset := 0
	var opts dd.UnmarshalOptions

	// TimeReal (4 bytes) - date of day downloaded
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for TimeReal")
	}
	timeReal, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("unmarshal TimeReal: %w", err)
	}
	activities.SetDateOfDay(timeReal)
	offset += 4

	// OdometerValueMidnight (3 bytes - OdometerShort)
	if offset+3 > len(data) {
		return nil, fmt.Errorf("insufficient data for OdometerValueMidnight")
	}
	odometer, err := opts.UnmarshalOdometer(data[offset : offset+3])
	if err != nil {
		return nil, fmt.Errorf("unmarshal OdometerValueMidnight: %w", err)
	}
	activities.SetOdometerMidnightKm(int32(odometer))
	offset += 3

	// VuCardIWData: 2 bytes (noOfIWRecords) + (noOfIWRecords * 129 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfIWRecords")
	}
	noOfIWRecords := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each CardIWRecord (129 bytes each for Gen1)
	cardIWRecords := make([]*ddv1.VuCardIWRecord, noOfIWRecords)
	for i := uint16(0); i < noOfIWRecords; i++ {
		const cardIWRecordSize = 129
		if offset+cardIWRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for CardIWRecord %d", i)
		}

		record, err := opts.UnmarshalVuCardIWRecord(data[offset : offset+cardIWRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal CardIWRecord %d: %w", i, err)
		}

		cardIWRecords[i] = record
		offset += cardIWRecordSize
	}
	activities.SetCardIwData(cardIWRecords)

	// VuActivityDailyData: 2 bytes (noOfActivityChanges) + (noOfActivityChanges * 2 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfActivityChanges")
	}
	noOfActivityChanges := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each ActivityChangeInfo (2 bytes each)
	activityChanges := make([]*ddv1.ActivityChangeInfo, noOfActivityChanges)
	for i := uint16(0); i < noOfActivityChanges; i++ {
		const activityChangeSize = 2
		if offset+activityChangeSize > len(data) {
			return nil, fmt.Errorf("insufficient data for ActivityChangeInfo %d", i)
		}

		activityChange, err := opts.UnmarshalActivityChangeInfo(data[offset : offset+activityChangeSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal activity change %d: %w", i, err)
		}
		activityChanges[i] = activityChange
		offset += activityChangeSize
	}
	activities.SetActivityChanges(activityChanges)

	// VuPlaceDailyWorkPeriodData: 1 byte (noOfPlaceRecords) + (noOfPlaceRecords * 28 bytes)
	// Note: Each record is 28 bytes (18 FullCardNumber + 10 PlaceRecord)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfPlaceRecords")
	}
	noOfPlaceRecords := data[offset]
	offset += 1

	// Parse each VuPlaceDailyWorkPeriodRecord using DD type (28 bytes each)
	// Extract only the PlaceRecord portion (VU records include FullCardNumber which we don't expose)
	placeRecords := make([]*ddv1.PlaceRecord, noOfPlaceRecords)
	for i := uint8(0); i < noOfPlaceRecords; i++ {
		const placeRecordSize = 28 // 18 bytes FullCardNumber + 10 bytes PlaceRecord
		if offset+placeRecordSize > len(data) {
			return nil, fmt.Errorf("insufficient data for PlaceRecord %d", i)
		}

		// Use DD type to parse the full VuPlaceDailyWorkPeriodRecord
		vuPlaceRecord, err := opts.UnmarshalVuPlaceDailyWorkPeriodRecord(data[offset : offset+placeRecordSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal VuPlaceDailyWorkPeriodRecord %d: %w", i, err)
		}

		// Extract just the PlaceRecord portion (DD type)
		placeRecords[i] = vuPlaceRecord.GetPlaceRecord()
		offset += placeRecordSize
	}
	activities.SetPlaces(placeRecords)

	// VuSpecificConditionData: 2 bytes (noOfSpecificConditionRecords) + (noOfSpecificConditionRecords * 5 bytes)
	if offset+2 > len(data) {
		return nil, fmt.Errorf("insufficient data for noOfSpecificConditionRecords")
	}
	noOfSpecificConditionRecords := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	// Parse each SpecificConditionRecord (5 bytes each)
	specificConditions := make([]*ddv1.SpecificConditionRecord, noOfSpecificConditionRecords)
	for i := uint16(0); i < noOfSpecificConditionRecords; i++ {
		const specificConditionSize = 5
		if offset+specificConditionSize > len(data) {
			return nil, fmt.Errorf("insufficient data for SpecificConditionRecord %d", i)
		}

		specificCondition, err := opts.UnmarshalSpecificConditionRecord(data[offset : offset+specificConditionSize])
		if err != nil {
			return nil, fmt.Errorf("unmarshal specific condition %d: %w", i, err)
		}
		specificConditions[i] = specificCondition
		offset += specificConditionSize
	}
	activities.SetSpecificConditions(specificConditions)

	// Store signature (extracted at the beginning)
	activities.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Activities Gen1 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
	}

	return activities, nil
}

// MarshalActivitiesGen1 marshals Gen1 Activities data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and has the correct length, it uses it as a canvas and paints semantic values over it.
// Otherwise, it creates a zero-filled canvas and encodes from semantic fields.
func (opts MarshalOptions) MarshalActivitiesGen1(activities *vuv1.ActivitiesGen1) ([]byte, error) {
	if activities == nil {
		return nil, fmt.Errorf("activities cannot be nil")
	}

	// Use raw_data if available (includes complete transfer with signature)
	// Full semantic marshalling requires implementing all record types
	raw := activities.GetRawData()
	if len(raw) > 0 {
		// raw_data contains complete transfer value (data + signature)
		return raw, nil
	}

	// TODO: Implement marshalling from semantic fields
	// This would require:
	// 1. Writing TimeReal
	// 2. Writing OdometerValueMidnight
	// 3. Writing VuCardIWData (count + records)
	// 4. Writing VuActivityDailyData (count + records)
	// 5. Writing VuPlaceDailyWorkPeriodData (count + records)
	// 6. Writing VuSpecificConditionData (count + records)
	// 7. Appending Signature
	return nil, fmt.Errorf("cannot marshal Activities Gen1 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeActivitiesGen1 anonymizes Gen1 Activities data.
// TODO: Implement full semantic anonymization (anonymize card numbers, timestamps, etc.).
func (opts AnonymizeOptions) anonymizeActivitiesGen1(activities *vuv1.ActivitiesGen1) *vuv1.ActivitiesGen1 {
	if activities == nil {
		return nil
	}
	result := proto.Clone(activities).(*vuv1.ActivitiesGen1)
	// Set signature to zero bytes (TV format: maintains structure)
	// Gen1 uses fixed 128-byte RSA-1024 signatures
	result.SetSignature(make([]byte, 128))

	// Note: We intentionally keep raw_data here because MarshalActivitiesGen1
	// currently requires raw_data (semantic marshalling not yet implemented).
	// Once semantic marshalling is implemented, we should clear raw_data and
	// implement full semantic anonymization of card_iw_data, activity_changes, etc.

	return result
}
