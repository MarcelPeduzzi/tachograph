package card

import (
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDrivingLicenceInfo parses the binary data for an EF_Driving_Licence_Info record.
//
// The data type `CardDrivingLicenceInformation` is specified in the Data Dictionary, Section 2.18.
//
// ASN.1 Definition:
//
//	CardDrivingLicenceInformation ::= SEQUENCE {
//	    drivingLicenceIssuingAuthority     Name,
//	    drivingLicenceIssuingNation        NationNumeric,
//	    drivingLicenceNumber               Name
//	}
func (opts UnmarshalOptions) unmarshalDrivingLicenceInfo(data []byte) (*cardv1.DrivingLicenceInfo, error) {
	const (
		lenCardDrivingLicenceInformation = 53 // CardDrivingLicenceInformation total size
	)

	if len(data) < lenCardDrivingLicenceInformation {
		return nil, errors.New("not enough data for DrivingLicenceInfo")
	}
	var dli cardv1.DrivingLicenceInfo
	offset := 0

	// Read driving licence issuing authority (36 bytes)
	if offset+36 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing authority")
	}
	authority, err := opts.UnmarshalStringValue(data[offset : offset+36])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence issuing authority: %w", err)
	}
	dli.SetDrivingLicenceIssuingAuthority(authority)
	offset += 36

	// Read driving licence issuing nation (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence issuing nation")
	}
	if nation, err := dd.UnmarshalEnum[ddv1.NationNumeric](data[offset]); err == nil {
		dli.SetDrivingLicenceIssuingNation(nation)
	} else {
		// Value not recognized - set UNRECOGNIZED (no unrecognized field for this type)
		dli.SetDrivingLicenceIssuingNation(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}
	offset++

	// Read driving licence number (16 bytes)
	if offset+16 > len(data) {
		return nil, fmt.Errorf("insufficient data for driving licence number")
	}
	licenceNumber, err := opts.UnmarshalIa5StringValue(data[offset : offset+16])
	if err != nil {
		return nil, fmt.Errorf("failed to read driving licence number: %w", err)
	}
	dli.SetDrivingLicenceNumber(licenceNumber)
	// offset += 16 // Not needed as this is the last field

	return &dli, nil
}

// MarshalDrivingLicenceInfo marshals the binary representation of DrivingLicenceInfo to bytes.
//
// The data type `CardDrivingLicenceInformation` is specified in the Data Dictionary, Section 2.18.
//
// ASN.1 Definition:
//
//	CardDrivingLicenceInformation ::= SEQUENCE {
//	    drivingLicenceIssuingAuthority     Name,
//	    drivingLicenceIssuingNation        NationNumeric,
//	    drivingLicenceNumber               Name
//	}
func (opts MarshalOptions) MarshalDrivingLicenceInfo(dli *cardv1.DrivingLicenceInfo) ([]byte, error) {
	if dli == nil {
		return nil, nil
	}

	var dst []byte

	authorityBytes, err := opts.MarshalStringValue(dli.GetDrivingLicenceIssuingAuthority())
	if err != nil {
		return nil, err
	}
	dst = append(dst, authorityBytes...)

	// Marshal nation enum to protocol value
	nationByte, err := dd.MarshalEnum(dli.GetDrivingLicenceIssuingNation())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nation: %w", err)
	}
	dst = append(dst, nationByte)

	licenceNumberBytes, err := opts.MarshalIa5StringValue(dli.GetDrivingLicenceNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driving licence number: %w", err)
	}
	dst = append(dst, licenceNumberBytes...)

	return dst, nil
}

// anonymizeDrivingLicenceInfo creates an anonymized copy of DrivingLicenceInfo,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeDrivingLicenceInfo(dli *cardv1.DrivingLicenceInfo) *cardv1.DrivingLicenceInfo {
	if dli == nil {
		return nil
	}

	anonymized := &cardv1.DrivingLicenceInfo{}

	// Create DD anonymize options
	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	// Anonymize issuing authority
	if dli.GetDrivingLicenceIssuingAuthority() != nil {
		anonymized.SetDrivingLicenceIssuingAuthority(ddOpts.AnonymizeStringValue(dli.GetDrivingLicenceIssuingAuthority()))
	}

	// Preserve country (structural info)
	anonymized.SetDrivingLicenceIssuingNation(dli.GetDrivingLicenceIssuingNation())

	// Anonymize licence number
	if dli.GetDrivingLicenceNumber() != nil {
		anonymized.SetDrivingLicenceNumber(ddOpts.AnonymizeIa5StringValue(dli.GetDrivingLicenceNumber()))
	}

	// Signature field left unset (nil) - TLV marshaller will omit the signature block

	return anonymized
}
