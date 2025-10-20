package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalOverviewGen2V2 parses Gen2 V2 Overview data from the complete transfer value.
//
// Gen2 V2 Overview structure is identical to Gen2 V1 with one addition:
// VehicleRegistrationNumberRecordArray is inserted between VehicleIdentificationNumberRecordArray
// and CurrentDateTimeRecordArray.
//
// ASN.1 Definition:
//
//	VuOverviewSecondGenV2 ::= SEQUENCE {
//	    memberStateCertificateRecordArray                MemberStateCertificateRecordArray,
//	    vuCertificateRecordArray                         VuCertificateRecordArray,
//	    vehicleIdentificationNumberRecordArray           VehicleIdentificationNumberRecordArray,
//	    vehicleRegistrationNumberRecordArray             VehicleRegistrationNumberRecordArray,   -- NEW in V2
//	    currentDateTimeRecordArray                       CurrentDateTimeRecordArray,
//	    vuDownloadablePeriodRecordArray                  VuDownloadablePeriodRecordArray,
//	    cardSlotsStatusRecordArray                       CardSlotsStatusRecordArray,
//	    vuDownloadActivityDataRecordArray                VuDownloadActivityDataRecordArray,
//	    vuCompanyLocksRecordArray                        VuCompanyLocksRecordArray,
//	    vuControlActivityRecordArray                     VuControlActivityRecordArray,
//	    signatureRecordArray                             SignatureRecordArray
//	}
//
// Each RecordArray has a 5-byte header:
//
//	recordType (1 byte) + recordSize (2 bytes, big-endian) + noOfRecords (2 bytes, big-endian)
//
// Note: This is a minimal implementation that stores raw_data for round-trip fidelity.
// Full semantic parsing of all RecordArrays is TODO.
func unmarshalOverviewGen2V2(value []byte) (*vuv1.OverviewGen2V2, error) {
	overview := &vuv1.OverviewGen2V2{}
	overview.SetRawData(value)

	// For now, store the raw data and validate structure by skipping through all record arrays
	offset := 0

	// Helper to skip a RecordArray
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(value, offset)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		offset += size
		return nil
	}

	// MemberStateCertificateRecordArray
	if err := skipRecordArray("MemberStateCertificate"); err != nil {
		return nil, err
	}

	// VUCertificateRecordArray
	if err := skipRecordArray("VUCertificate"); err != nil {
		return nil, err
	}

	// VehicleIdentificationNumberRecordArray
	if err := skipRecordArray("VehicleIdentificationNumber"); err != nil {
		return nil, err
	}

	// VehicleRegistrationNumberRecordArray (Gen2 V2 addition)
	if err := skipRecordArray("VehicleRegistrationNumber"); err != nil {
		return nil, err
	}

	// CurrentDateTimeRecordArray
	if err := skipRecordArray("CurrentDateTime"); err != nil {
		return nil, err
	}

	// VuDownloadablePeriodRecordArray
	if err := skipRecordArray("VuDownloadablePeriod"); err != nil {
		return nil, err
	}

	// CardSlotsStatusRecordArray
	if err := skipRecordArray("CardSlotsStatus"); err != nil {
		return nil, err
	}

	// VuDownloadActivityDataRecordArray
	if err := skipRecordArray("VuDownloadActivityData"); err != nil {
		return nil, err
	}

	// VuCompanyLocksRecordArray
	if err := skipRecordArray("VuCompanyLocks"); err != nil {
		return nil, err
	}

	// VuControlActivityRecordArray
	if err := skipRecordArray("VuControlActivity"); err != nil {
		return nil, err
	}

	// SignatureRecordArray is now handled separately in raw parsing, not part of value

	// Verify we consumed exactly the right amount of data
	if offset != len(value) {
		return nil, fmt.Errorf("Overview Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(value))
	}

	// TODO: Implement full semantic parsing of all record arrays
	// For now, raw_data contains all the information needed for round-trip testing

	return overview, nil
}

// MarshalOverviewGen2V2 marshals Gen2 V2 Overview data using raw data painting.
//
// This function implements the raw data painting pattern: if raw_data is available
// and matches the structure, it uses it as the output. Otherwise, it would need to
// construct from semantic fields (currently not implemented).
func (opts MarshalOptions) MarshalOverviewGen2V2(overview *vuv1.OverviewGen2V2) ([]byte, error) {
	if overview == nil {
		return nil, fmt.Errorf("overview cannot be nil")
	}

	// For Gen2 structures with RecordArrays, raw data painting is straightforward:
	// we use the raw_data if available
	raw := overview.GetRawData()
	if len(raw) > 0 {
		return raw, nil
	}

	// TODO: Implement marshalling from semantic fields
	// This would require constructing all RecordArrays from semantic data
	return nil, fmt.Errorf("cannot marshal Overview Gen2 V2 without raw_data (semantic marshalling not yet implemented)")
}

// anonymizeOverviewGen2V2 anonymizes Gen2 V2 Overview data.
func (opts AnonymizeOptions) anonymizeOverviewGen2V2(overview *vuv1.OverviewGen2V2) *vuv1.OverviewGen2V2 {
	if overview == nil {
		return nil
	}

	result := proto.Clone(overview).(*vuv1.OverviewGen2V2)

	// Anonymize VIN
	if vin := result.GetVehicleIdentificationNumber(); vin != nil {
		vin.SetValue("TESTVIN1234567890")
	}

	// Anonymize VRN
	if vrn := result.GetVehicleRegistrationNumber(); vrn != nil {
		vrn.SetValue("TEST123")
	}

	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})
	result.SetRawData(nil)

	return result
}
