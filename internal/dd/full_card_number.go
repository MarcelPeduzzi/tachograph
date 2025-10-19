package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalFullCardNumber parses full card number data.
//
// The data type `FullCardNumber` is specified in the Data Dictionary, Section 2.73.
//
// ASN.1 Definition:
//
//	FullCardNumber ::= SEQUENCE {
//	    cardType EquipmentType,
//	    cardIssuingMemberState NationNumeric,
//	    cardNumber CardNumber
//	}
//
//	CardNumber ::= CHOICE {
//	    driverIdentification   SEQUENCE { ... },
//	    ownerIdentification    SEQUENCE { ... }
//	}
//
// Binary Layout (fixed length, 18 bytes):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (16 bytes): CardNumber CHOICE based on card type (padded to 16 bytes)
func (opts UnmarshalOptions) UnmarshalFullCardNumber(data []byte) (*ddv1.FullCardNumber, error) {
	const lenFullCardNumber = 18

	if len(data) != lenFullCardNumber {
		return nil, fmt.Errorf("invalid data length for FullCardNumber: got %d, want %d", len(data), lenFullCardNumber)
	}

	cardNumber := &ddv1.FullCardNumber{}

	// Parse card type (1 byte)
	cardType, err := UnmarshalEnum[ddv1.EquipmentType](data[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card type: %w", err)
	}
	cardNumber.SetCardType(cardType)

	// Parse issuing member state (1 byte)
	issuingState := data[1]
	cardNumber.SetCardIssuingMemberState(ddv1.NationNumeric(issuingState))

	// Parse card number based on card type (16 bytes, may have padding)
	cardNumberData := data[2:18]
	switch cardType {
	case ddv1.EquipmentType_DRIVER_CARD:
		// DriverIdentification is 14 bytes + 2 bytes padding
		driverID, err := opts.UnmarshalDriverIdentification(cardNumberData[:14])
		if err != nil {
			return nil, fmt.Errorf("failed to parse driver identification: %w", err)
		}
		cardNumber.SetDriverIdentification(driverID)
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		// OwnerIdentification is 16 bytes (no padding)
		ownerID, err := opts.UnmarshalOwnerIdentification(cardNumberData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse owner identification: %w", err)
		}
		cardNumber.SetOwnerIdentification(ownerID)
	default:
		return nil, fmt.Errorf("unsupported card type: %d", cardType)
	}

	return cardNumber, nil
}

// MarshalFullCardNumber marshals full card number data to bytes.
//
// The data type `FullCardNumber` is specified in the Data Dictionary, Section 2.73.
//
// ASN.1 Definition:
//
//	FullCardNumber ::= SEQUENCE {
//	    cardType EquipmentType,
//	    cardIssuingMemberState NationNumeric,
//	    cardNumber CardNumber
//	}
//
//	CardNumber ::= CHOICE {
//	    driverIdentification   SEQUENCE { ... },
//	    ownerIdentification    SEQUENCE { ... }
//	}
//
// Binary Layout (fixed length, 18 bytes):
//   - Card Type (1 byte): EquipmentType
//   - Issuing Member State (1 byte): NationNumeric
//   - Card Number (16 bytes): CardNumber CHOICE based on card type (padded to 16 bytes)
func (opts MarshalOptions) MarshalFullCardNumber(cardNumber *ddv1.FullCardNumber) ([]byte, error) {
	if cardNumber == nil {
		return nil, fmt.Errorf("cardNumber cannot be nil")
	}

	var dst []byte

	// Marshal card type (1 byte)
	cardTypeByte, err := MarshalEnum(cardNumber.GetCardType())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card type: %w", err)
	}
	dst = append(dst, cardTypeByte)

	// Marshal issuing member state (1 byte)
	dst = append(dst, byte(cardNumber.GetCardIssuingMemberState()))

	// Marshal card number based on card type (16 bytes with padding if needed)
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			// DriverIdentification is 14 bytes, pad to 16
			driverBytes, err := opts.MarshalDriverIdentification(driverID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal driver identification: %w", err)
			}
			dst = append(dst, driverBytes...)
			// Add 2 bytes padding
			dst = append(dst, 0x00, 0x00)
		} else {
			// Empty driver ID: 16 zero bytes
			dst = append(dst, make([]byte, 16)...)
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			// OwnerIdentification is 16 bytes (no padding needed)
			ownerBytes, err := opts.MarshalOwnerIdentification(ownerID)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal owner identification: %w", err)
			}
			dst = append(dst, ownerBytes...)
		} else {
			// Empty owner ID: 16 zero bytes
			dst = append(dst, make([]byte, 16)...)
		}
	default:
		// Unknown card type: 16 zero bytes
		dst = append(dst, make([]byte, 16)...)
	}

	return dst, nil
}

// MarshalFullCardNumberAsString marshals a FullCardNumber structure as a string representation.
// This is used for display purposes and has a maximum length constraint.
func (opts MarshalOptions) MarshalFullCardNumberAsString(cardNumber *ddv1.FullCardNumber, maxLen int) ([]byte, error) {
	if cardNumber == nil {
		return nil, fmt.Errorf("cardNumber cannot be nil")
	}

	// Handle the CardNumber CHOICE based on card type
	switch cardNumber.GetCardType() {
	case ddv1.EquipmentType_DRIVER_CARD:
		if driverID := cardNumber.GetDriverIdentification(); driverID != nil {
			return opts.MarshalIa5StringValue(driverID.GetDriverIdentificationNumber())
		}
	case ddv1.EquipmentType_WORKSHOP_CARD, ddv1.EquipmentType_COMPANY_CARD:
		if ownerID := cardNumber.GetOwnerIdentification(); ownerID != nil {
			return opts.MarshalIa5StringValue(ownerID.GetOwnerIdentification())
		}
	}

	return opts.MarshalStringValue(nil)
}
