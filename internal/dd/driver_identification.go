package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalDriverIdentification parses driver identification data.
//
// The data type `DriverIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	driverIdentification SEQUENCE {
//	    driverIdentificationNumber IA5String(SIZE(14)),
//	    cardReplacementIndex CardReplacementIndex,
//	    cardRenewalIndex CardRenewalIndex
//	}
//
// Binary Layout (16 bytes):
//   - Driver Identification Number (14 bytes): IA5String
//   - Card Replacement Index (1 byte): IA5String
//   - Card Renewal Index (1 byte): IA5String
func (opts UnmarshalOptions) UnmarshalDriverIdentification(data []byte) (*ddv1.DriverIdentification, error) {
	const (
		lenDriverIdentification = 16
	)

	if len(data) != lenDriverIdentification {
		return nil, fmt.Errorf("invalid data length for DriverIdentification: got %d, want %d", len(data), lenDriverIdentification)
	}

	driverID := &ddv1.DriverIdentification{}

	// Parse driver identification number (14 bytes)
	identificationNumber, err := opts.UnmarshalIa5StringValue(data[0:14])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver identification number: %w", err)
	}
	driverID.SetDriverIdentificationNumber(identificationNumber)

	// Parse card replacement index (1 byte)
	replacementIndex, err := opts.UnmarshalIa5StringValue(data[14:15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card replacement index: %w", err)
	}
	driverID.SetCardReplacementIndex(replacementIndex)

	// Parse card renewal index (1 byte)
	renewalIndex, err := opts.UnmarshalIa5StringValue(data[15:16])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card renewal index: %w", err)
	}
	driverID.SetCardRenewalIndex(renewalIndex)

	return driverID, nil
}

// MarshalDriverIdentification marshals driver identification data to bytes.
//
// The data type `DriverIdentification` is specified in the Data Dictionary, Section 2.26.
//
// ASN.1 Definition:
//
//	driverIdentification SEQUENCE {
//	    driverIdentificationNumber IA5String(SIZE(14)),
//	    cardReplacementIndex CardReplacementIndex,
//	    cardRenewalIndex CardRenewalIndex
//	}
//
// Binary Layout (16 bytes):
//   - Driver Identification Number (14 bytes): IA5String
//   - Card Replacement Index (1 byte): IA5String
//   - Card Renewal Index (1 byte): IA5String
func (opts MarshalOptions) MarshalDriverIdentification(driverID *ddv1.DriverIdentification) ([]byte, error) {
	if driverID == nil {
		return nil, fmt.Errorf("driverID cannot be nil")
	}

	var dst []byte

	// Marshal driver identification number (14 bytes)
	idNumberBytes, err := opts.MarshalIa5StringValue(driverID.GetDriverIdentificationNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driver identification number: %w", err)
	}
	dst = append(dst, idNumberBytes...)

	// Marshal card replacement index (1 byte)
	replacementBytes, err := opts.MarshalIa5StringValue(driverID.GetCardReplacementIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card replacement index: %w", err)
	}
	dst = append(dst, replacementBytes...)

	// Marshal card renewal index (1 byte)
	renewalBytes, err := opts.MarshalIa5StringValue(driverID.GetCardRenewalIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card renewal index: %w", err)
	}
	dst = append(dst, renewalBytes...)

	return dst, nil
}

// AnonymizeDriverIdentification creates an anonymized copy of DriverIdentification,
// replacing the driver identification number with a safe, deterministic value while
// maintaining the correct format and length.
func AnonymizeDriverIdentification(driverID *ddv1.DriverIdentification) *ddv1.DriverIdentification {
	if driverID == nil {
		return nil
	}
	result := &ddv1.DriverIdentification{}
	// Anonymize driver identification number (IA5String, 14 bytes)
	testDriverID := &ddv1.Ia5StringValue{}
	testDriverID.SetValue("DRIVER00000001")
	testDriverID.SetLength(14)
	result.SetDriverIdentificationNumber(testDriverID)

	// Card replacement index (IA5String, 1 byte)
	replacementIndex := &ddv1.Ia5StringValue{}
	replacementIndex.SetValue("0")
	replacementIndex.SetLength(1)
	result.SetCardReplacementIndex(replacementIndex)

	// Card renewal index (IA5String, 1 byte)
	renewalIndex := &ddv1.Ia5StringValue{}
	renewalIndex.SetValue("0")
	renewalIndex.SetLength(1)
	result.SetCardRenewalIndex(renewalIndex)

	return result
}
