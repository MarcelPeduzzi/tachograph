package card

import (
	"encoding/binary"
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalCardControlActivityData unmarshals control activity data from a card EF.
//
// The data type `CardControlActivityDataRecord` is specified in the Data Dictionary, Section 2.15.
//
// ASN.1 Definition:
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
func (opts UnmarshalOptions) unmarshalControlActivityData(data []byte) (*cardv1.ControlActivityData, error) {
	const (
		lenCardControlActivityDataRecord = 46 // CardControlActivityDataRecord total size
	)

	if len(data) < lenCardControlActivityDataRecord {
		return nil, fmt.Errorf("insufficient data for control activity data")
	}
	var target cardv1.ControlActivityData
	controlTime := binary.BigEndian.Uint32(data[1:5])
	if controlTime == 0 {
		target.SetValid(false)
		target.SetRawData(data)
		return &target, nil
	}
	target.SetValid(true)

	offset := 0

	// Read control type (1 byte)
	if offset+1 > len(data) {
		return nil, fmt.Errorf("insufficient data for control type")
	}
	controlType, err := opts.UnmarshalControlType(data[offset : offset+1])
	if err != nil {
		return nil, fmt.Errorf("failed to read control type: %w", err)
	}
	target.SetControlType(controlType)
	offset++

	// Read control time (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control time")
	}
	controlTimestamp, err2 := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err2 != nil {
		return nil, fmt.Errorf("failed to parse control time: %w", err2)
	}
	target.SetControlTime(controlTimestamp)
	offset += 4

	// Read control card number (18 bytes) - this should be parsed as a proper FullCardNumberAndGeneration
	// For now, create a basic structure - this needs proper protocol parsing
	fullCardNumberAndGeneration := &ddv1.FullCardNumberAndGeneration{}

	fullCardNumber := &ddv1.FullCardNumber{}
	fullCardNumber.SetCardType(ddv1.EquipmentType_DRIVER_CARD)

	// Read the card number as IA5 string
	if offset+18 > len(data) {
		return nil, fmt.Errorf("insufficient data for control card number")
	}
	cardNumberStr, err := opts.UnmarshalIa5StringValue(data[offset : offset+18])
	if err != nil {
		return nil, fmt.Errorf("failed to read control card number: %w", err)
	}
	offset += 18

	// Create driver identification with the card number
	driverID := &ddv1.DriverIdentification{}
	driverID.SetDriverIdentificationNumber(cardNumberStr)
	fullCardNumber.SetDriverIdentification(driverID)

	// Set the full card number in the generation wrapper
	fullCardNumberAndGeneration.SetFullCardNumber(fullCardNumber)
	// Default to Generation 1 for now - this should be determined from context
	fullCardNumberAndGeneration.SetGeneration(ddv1.Generation_GENERATION_1)

	target.SetControlCardNumber(fullCardNumberAndGeneration)

	// Read vehicle registration (15 bytes: 1 byte nation + 14 bytes number)
	if offset+15 > len(data) {
		return nil, fmt.Errorf("insufficient data for vehicle registration")
	}
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[offset : offset+15])
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicle registration: %w", err)
	}
	offset += 15
	target.SetControlVehicleRegistration(vehicleReg)

	// Read control download period begin (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period begin")
	}
	controlDownloadPeriodBegin, err3 := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err3 != nil {
		return nil, fmt.Errorf("failed to parse control download period begin: %w", err3)
	}
	target.SetControlDownloadPeriodBegin(controlDownloadPeriodBegin)
	offset += 4

	// Read control download period end (4 bytes)
	if offset+4 > len(data) {
		return nil, fmt.Errorf("insufficient data for control download period end")
	}
	controlDownloadPeriodEnd, err4 := opts.UnmarshalTimeReal(data[offset : offset+4])
	if err4 != nil {
		return nil, fmt.Errorf("failed to parse control download period end: %w", err4)
	}
	target.SetControlDownloadPeriodEnd(controlDownloadPeriodEnd)
	// offset += 4 // Not needed as this is the last field

	return &target, nil
}

// MarshalCardControlActivityData marshals control activity data.
//
// The data type `CardControlActivityDataRecord` is specified in the Data Dictionary, Section 2.15.
//
// ASN.1 Definition:
//
//	CardControlActivityDataRecord ::= SEQUENCE {
//	    controlType                        ControlType,
//	    controlTime                        TimeReal,
//	    controlCardNumber                  FullCardNumber,
//	    controlVehicleRegistration         VehicleRegistrationIdentification,
//	    controlDownloadPeriodBegin         TimeReal,
//	    controlDownloadPeriodEnd           TimeReal
//	}
func (opts MarshalOptions) MarshalCardControlActivityData(controlData *cardv1.ControlActivityData) ([]byte, error) {
	if controlData == nil {
		return nil, nil
	}

	if !controlData.GetValid() {
		// Non-valid record: use preserved raw data
		rawData := controlData.GetRawData()
		if len(rawData) != 46 {
			// Fallback to zeros if raw data is invalid
			return make([]byte, 46), nil
		}
		return rawData, nil
	}

	var data []byte

	// Valid record: serialize semantic data
	// Control type (1 byte)
	controlType := controlData.GetControlType()
	var controlTypeByte byte
	if controlType != nil {
		// Build bitmask from boolean fields
		// Structure: 'cvpdexxx'B
		// - 'c': card downloading
		// - 'v': VU downloading
		// - 'p': printing
		// - 'd': display
		// - 'e': calibration checking
		if controlType.GetCardDownloading() {
			controlTypeByte |= 0x80 // bit 7
		}
		if controlType.GetVuDownloading() {
			controlTypeByte |= 0x40 // bit 6
		}
		if controlType.GetPrinting() {
			controlTypeByte |= 0x20 // bit 5
		}
		if controlType.GetDisplay() {
			controlTypeByte |= 0x10 // bit 4
		}
		if controlType.GetCalibrationChecking() {
			controlTypeByte |= 0x08 // bit 3
		}
	}
	data = append(data, controlTypeByte)

	// Control time (4 bytes)

	timeBytes, err := opts.MarshalTimeReal(controlData.GetControlTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal control time: %w", err)
	}
	data = append(data, timeBytes...)

	// Control card number (18 bytes)
	cardNumberBytes, err := opts.MarshalFullCardNumberAsString(controlData.GetControlCardNumber().GetFullCardNumber(), 18)
	if err != nil {
		return nil, err
	}
	data = append(data, cardNumberBytes...)

	// Vehicle registration (15 bytes total: 1 byte nation + 14 bytes number)
	vehicleRegBytes, err := opts.MarshalVehicleRegistration(controlData.GetControlVehicleRegistration())
	if err != nil {
		return nil, err
	}
	data = append(data, vehicleRegBytes...)

	// Control download period begin (4 bytes)
	beginBytes, err := opts.MarshalTimeReal(controlData.GetControlDownloadPeriodBegin())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal download period begin: %w", err)
	}
	data = append(data, beginBytes...)

	// Control download period end (4 bytes)
	endBytes, err := opts.MarshalTimeReal(controlData.GetControlDownloadPeriodEnd())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal download period end: %w", err)
	}
	data = append(data, endBytes...)

	return data, nil
}

// anonymizeControlActivityData creates an anonymized copy of ControlActivityData,
// replacing sensitive information with static, deterministic test values.
func (opts AnonymizeOptions) anonymizeControlActivityData(ca *cardv1.ControlActivityData) *cardv1.ControlActivityData {
	if ca == nil {
		return nil
	}

	anonymized := &cardv1.ControlActivityData{}
	anonymized.SetValid(ca.GetValid())

	if !ca.GetValid() {
		anonymized.SetRawData(ca.GetRawData())
		return anonymized
	}

	// Preserve control type (categorical)
	anonymized.SetControlType(ca.GetControlType())

	// Static test timestamp: 2020-01-01 00:00:00 UTC
	anonymized.SetControlTime(&timestamppb.Timestamp{Seconds: 1577836800})

	// Anonymize control card number
	if cardNum := ca.GetControlCardNumber(); cardNum != nil {
		anonymizedCardNum := &ddv1.FullCardNumberAndGeneration{}
		if fcn := cardNum.GetFullCardNumber(); fcn != nil {
			anonymizedFCN := &ddv1.FullCardNumber{}
			anonymizedFCN.SetCardType(fcn.GetCardType())
			anonymizedFCN.SetCardIssuingMemberState(ddv1.NationNumeric_FINLAND)

			// Anonymize driver or owner identification
			if driverID := fcn.GetDriverIdentification(); driverID != nil {
				anonymizedDriverID := &ddv1.DriverIdentification{}
				if driverID.GetDriverIdentificationNumber() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("CTRL-DRV-001")

					sv.SetLength(14)
					anonymizedDriverID.SetDriverIdentificationNumber(sv)
				}
				if driverID.GetCardReplacementIndex() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("0")

					sv.SetLength(1)
					anonymizedDriverID.SetCardReplacementIndex(sv)
				}
				if driverID.GetCardRenewalIndex() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("0")

					sv.SetLength(1)
					anonymizedDriverID.SetCardRenewalIndex(sv)
				}
				anonymizedFCN.SetDriverIdentification(anonymizedDriverID)
			} else if ownerID := fcn.GetOwnerIdentification(); ownerID != nil {
				anonymizedOwnerID := &ddv1.OwnerIdentification{}
				if ownerID.GetOwnerIdentification() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("CTRL-OWN-001")

					sv.SetLength(13)
					anonymizedOwnerID.SetOwnerIdentification(sv)
				}
				if ownerID.GetConsecutiveIndex() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("0")

					sv.SetLength(1)
					anonymizedOwnerID.SetConsecutiveIndex(sv)
				}
				if ownerID.GetReplacementIndex() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("0")

					sv.SetLength(1)
					anonymizedOwnerID.SetReplacementIndex(sv)
				}
				if ownerID.GetRenewalIndex() != nil {
					sv := &ddv1.Ia5StringValue{}
					sv.SetValue("0")

					sv.SetLength(1)
					anonymizedOwnerID.SetRenewalIndex(sv)
				}
				anonymizedFCN.SetOwnerIdentification(anonymizedOwnerID)
			}

			anonymizedCardNum.SetFullCardNumber(anonymizedFCN)
		}
		anonymized.SetControlCardNumber(anonymizedCardNum)
	}

	// Anonymize vehicle registration
	if vehicleReg := ca.GetControlVehicleRegistration(); vehicleReg != nil {
		anonymizedReg := &ddv1.VehicleRegistrationIdentification{}
		anonymizedReg.SetNation(ddv1.NationNumeric_FINLAND)

		// VehicleRegistrationNumber is: 1 byte code page + 13 bytes data
		testRegNum := &ddv1.StringValue{}
		testRegNum.SetValue("TEST-VRN")
		testRegNum.SetEncoding(ddv1.Encoding_ISO_8859_1) // Code page 1 (Latin-1)
		testRegNum.SetLength(13)                         // Length of data bytes (not including code page)
		anonymizedReg.SetNumber(testRegNum)

		anonymized.SetControlVehicleRegistration(anonymizedReg)
	}

	// Static download period
	anonymized.SetControlDownloadPeriodBegin(&timestamppb.Timestamp{Seconds: 1577836800})
	anonymized.SetControlDownloadPeriodEnd(&timestamppb.Timestamp{Seconds: 1577836800 + 86400}) // 1 day later

	// Signature field left unset (nil) - TLV marshaller will omit the signature block

	return anonymized
}
