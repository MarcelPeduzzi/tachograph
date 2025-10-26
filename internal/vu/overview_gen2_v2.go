package vu

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// unmarshalOverviewGen2V2 parses Gen2 V2 Overview data from the complete transfer value.
//
// This function accepts the complete transfer value including the signature appended
// at the end, as specified in Appendix 7, Section 2.2.6.
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
	// Split transfer value into data and signature
	// Gen2 uses variable-length ECDSA signatures stored as SignatureRecordArray
	// We use the sizeOf function to determine where to split
	totalSize, signatureSize, err := sizeOfOverviewGen2V2(value)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize != len(value) {
		return nil, fmt.Errorf("size mismatch: calculated %d, got %d", totalSize, len(value))
	}

	dataSize := totalSize - signatureSize
	data := value[:dataSize]
	signature := value[dataSize:]

	overview := &vuv1.OverviewGen2V2{}
	overview.SetRawData(value) // Store complete transfer value for painting

	// For now, store the raw data and validate structure by skipping through all record arrays
	offset := 0

	// Helper to skip a RecordArray
	skipRecordArray := func(name string) error {
		size, err := sizeOfRecordArray(data, offset)
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

	// Store signature (extracted at the beginning)
	overview.SetSignature(signature)

	// Verify we consumed exactly the right amount of data
	if offset != len(data) {
		return nil, fmt.Errorf("Overview Gen2 V2 parsing mismatch: parsed %d bytes, expected %d", offset, len(data))
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
		// raw_data contains complete transfer value (data + signature)
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

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize VIN
	if vin := result.GetVehicleIdentificationNumber(); vin != nil {
		result.SetVehicleIdentificationNumber(ddOpts.AnonymizeIa5StringValue(vin))
	}

	// Anonymize VRN (just an IA5String in Gen2 V2)
	if vrn := result.GetVehicleRegistrationNumber(); vrn != nil {
		result.SetVehicleRegistrationNumber(ddOpts.AnonymizeIa5StringValue(vrn))
	}

	// Set signature to empty bytes (TV format: maintains structure)
	// Gen2 uses variable-length ECDSA signatures
	result.SetSignature([]byte{})

	// Anonymize download activities
	var anonymizedDownloadActivities []*vuv1.OverviewGen2V2_DownloadActivity
	for _, activity := range result.GetDownloadActivities() {
		anonActivity := proto.Clone(activity).(*vuv1.OverviewGen2V2_DownloadActivity)
		// Anonymize card number and generation
		if anonActivity.GetFullCardNumberAndGeneration() != nil {
			anonActivity.SetFullCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(anonActivity.GetFullCardNumberAndGeneration()))
		}
		// Anonymize company/workshop name
		if anonActivity.GetCompanyOrWorkshopName() != nil {
			anonActivity.SetCompanyOrWorkshopName(ddOpts.AnonymizeStringValue(anonActivity.GetCompanyOrWorkshopName()))
		}
		anonymizedDownloadActivities = append(anonymizedDownloadActivities, anonActivity)
	}
	result.SetDownloadActivities(anonymizedDownloadActivities)

	// Anonymize company locks
	var anonymizedCompanyLocks []*vuv1.OverviewGen2V2_CompanyLock
	for _, lock := range result.GetCompanyLocks() {
		anonLock := proto.Clone(lock).(*vuv1.OverviewGen2V2_CompanyLock)
		// Anonymize company name
		if anonLock.GetCompanyName() != nil {
			anonLock.SetCompanyName(ddOpts.AnonymizeStringValue(anonLock.GetCompanyName()))
		}
		// Anonymize company address
		if anonLock.GetCompanyAddress() != nil {
			anonLock.SetCompanyAddress(ddOpts.AnonymizeStringValue(anonLock.GetCompanyAddress()))
		}
		// Anonymize company card number and generation
		if anonLock.GetCompanyCardNumberAndGeneration() != nil {
			anonLock.SetCompanyCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(anonLock.GetCompanyCardNumberAndGeneration()))
		}
		anonymizedCompanyLocks = append(anonymizedCompanyLocks, anonLock)
	}
	result.SetCompanyLocks(anonymizedCompanyLocks)

	// Anonymize control activities
	var anonymizedControlActivities []*vuv1.OverviewGen2V2_ControlActivity
	for _, activity := range result.GetControlActivities() {
		anonActivity := proto.Clone(activity).(*vuv1.OverviewGen2V2_ControlActivity)
		// Anonymize control card number and generation
		if anonActivity.GetControlCardNumberAndGeneration() != nil {
			anonActivity.SetControlCardNumberAndGeneration(ddOpts.AnonymizeFullCardNumberAndGeneration(anonActivity.GetControlCardNumberAndGeneration()))
		}
		anonymizedControlActivities = append(anonymizedControlActivities, anonActivity)
	}
	result.SetControlActivities(anonymizedControlActivities)

	// Clear raw_data
	result.ClearRawData()

	return result
}
