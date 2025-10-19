package card

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalCardCurrentUsage unmarshals current usage data from a card EF.
//
// The data type `CardCurrentUse` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
func (opts UnmarshalOptions) unmarshalCurrentUsage(data []byte) (*cardv1.CurrentUsage, error) {
	const (
		lenCardCurrentUse = 19 // 4 bytes time + 15 bytes vehicle registration
	)

	if len(data) < lenCardCurrentUse {
		return nil, fmt.Errorf("insufficient data for current usage")
	}
	var target cardv1.CurrentUsage
	offset := 0

	// Read session open time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for session open time")
	}
	sessionOpenTime, err := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse session open time: %w", err)
	}
	target.SetSessionOpenTime(sessionOpenTime)
	offset += 4

	// Read session open vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+15 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration")
	}
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	// offset += 15 // Not needed as this is the last field
	target.SetSessionOpenVehicle(vehicleReg)
	return &target, nil
}

// MarshalCurrentUsage marshals current usage data to bytes.
//
// The data type `CardCurrentUse` is specified in the Data Dictionary, Section 2.16.
//
// ASN.1 Definition:
//
//	CardCurrentUse ::= SEQUENCE {
//	    sessionOpenTime                   TimeReal,
//	    sessionOpenVehicle                VehicleRegistrationIdentification
//	}
func (opts MarshalOptions) MarshalCurrentUsage(currentUsage *cardv1.CurrentUsage) ([]byte, error) {
	if currentUsage == nil {
		return nil, nil
	}

	var data []byte

	// Session open time (4 bytes)
	timeBytes, err := opts.MarshalTimeReal(currentUsage.GetSessionOpenTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session open time: %w", err)
	}
	data = append(data, timeBytes...)

	// Session open vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	vehicleReg := currentUsage.GetSessionOpenVehicle()
	if vehicleReg != nil {
		regBytes, err := opts.MarshalVehicleRegistration(vehicleReg)
		if err != nil {
			return nil, err
		}
		data = append(data, regBytes...)
	} else {
		// No vehicle registration - pad with zeros
		data = append(data, make([]byte, 15)...)
	}

	return data, nil
}

// anonymizeCurrentUsage creates an anonymized copy of CurrentUsage,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeCurrentUsage(cu *cardv1.CurrentUsage) *cardv1.CurrentUsage {
	if cu == nil {
		return nil
	}

	anonymized := &cardv1.CurrentUsage{}

	// Use static test timestamp: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	anonymized.SetSessionOpenTime(&timestamppb.Timestamp{Seconds: 1577836800})

	// Anonymize vehicle registration
	if vehicleReg := cu.GetSessionOpenVehicle(); vehicleReg != nil {
		anonymizedReg := &ddv1.VehicleRegistrationIdentification{}

		// Country → FINLAND (always)
		anonymizedReg.SetNation(ddv1.NationNumeric_FINLAND)

		// Registration number → static test value
		// VehicleRegistrationNumber is: 1 byte code page + 13 bytes data
		testRegNum := &ddv1.StringValue{}
		testRegNum.SetValue("TEST-123")
		testRegNum.SetEncoding(ddv1.Encoding_ISO_8859_1) // Code page 1 (Latin-1)
		testRegNum.SetLength(13)                         // Length of data bytes (not including code page)
		anonymizedReg.SetNumber(testRegNum)

		anonymized.SetSessionOpenVehicle(anonymizedReg)
	}

	return anonymized
}
