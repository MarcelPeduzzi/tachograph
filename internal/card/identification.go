package card

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalDriverCardIdentification parses the binary data for an EF_Identification record
// from a driver card.
//
// The data type `DriverCardIdentification` combines two structures:
//   - CardIdentification (Data Dictionary Section 2.24) - 65 bytes
//   - DriverCardHolderIdentification (Data Dictionary Section 2.62) - 78 bytes
//
// Total binary size: 143 bytes (fixed)
//
// CardIdentification ASN.1 Specification (Data Dictionary Section 2.24):
//
//	CardIdentification ::= SEQUENCE {
//	    cardIssuingMemberState    NationNumeric,
//	    cardNumber                CardNumber,
//	    cardIssuingAuthorityName  Name,
//	    cardIssueDate             TimeReal,
//	    cardValidityBegin         TimeReal,
//	    cardExpiryDate            TimeReal
//	}
//
// DriverCardHolderIdentification ASN.1 Specification (Data Dictionary Section 2.62):
//
//	DriverCardHolderIdentification ::= SEQUENCE {
//	    cardHolderName                  HolderName,
//	    cardHolderBirthDate             Datef,
//	    cardHolderPreferredLanguage     Language
//	}
//
// Binary Layout (143 bytes total):
//   - cardIssuingMemberState: 1 byte (NationNumeric)
//   - driverIdentification: 14 bytes (IA5String for driver cards)
//   - cardReplacementIndex: 1 byte (IA5String)
//   - cardRenewalIndex: 1 byte (IA5String)
//   - cardIssuingAuthorityName: 36 bytes (1 byte code page + 35 bytes data)
//   - cardIssueDate: 4 bytes (TimeReal)
//   - cardValidityBegin: 4 bytes (TimeReal)
//   - cardExpiryDate: 4 bytes (TimeReal)
//   - cardHolderSurname: 36 bytes (1 byte code page + 35 bytes data)
//   - cardHolderFirstNames: 36 bytes (1 byte code page + 35 bytes data)
//   - cardHolderBirthDate: 4 bytes (Datef)
//   - cardHolderPreferredLanguage: 2 bytes (Language)
func (opts UnmarshalOptions) unmarshalDriverCardIdentification(data []byte) (*cardv1.DriverCardIdentification, error) {
	const (
		lenDriverCardIdentification = 143

		// CardIdentification part (65 bytes)
		idxIssuingMemberState   = 0
		lenIssuingMemberState   = 1
		idxDriverIdentification = 1
		lenDriverIdentification = 14
		idxReplacementIndex     = 15
		lenReplacementIndex     = 1
		idxRenewalIndex         = 16
		lenRenewalIndex         = 1
		idxAuthorityName        = 17
		lenAuthorityName        = 36
		idxIssueDate            = 53
		lenIssueDate            = 4
		idxValidityBegin        = 57
		lenValidityBegin        = 4
		idxExpiryDate           = 61
		lenExpiryDate           = 4

		// DriverCardHolderIdentification part (78 bytes)
		idxSurname    = 65
		lenSurname    = 36
		idxFirstNames = 101
		lenFirstNames = 36
		idxBirthDate  = 137
		lenBirthDate  = 4
		idxLanguage   = 141
		lenLanguage   = 2
	)

	if len(data) != lenDriverCardIdentification {
		return nil, fmt.Errorf("invalid data length for DriverCardIdentification: got %d, want %d", len(data), lenDriverCardIdentification)
	}

	var id cardv1.DriverCardIdentification

	// Parse CardIdentification part (65 bytes)

	// Nation (1 byte)
	if nation, err := dd.UnmarshalEnum[ddv1.NationNumeric](data[idxIssuingMemberState]); err == nil {
		id.SetCardIssuingMemberState(nation)
	} else {
		id.SetCardIssuingMemberState(ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED)
	}

	// DriverIdentification (14 + 1 + 1 = 16 bytes for driver cards)
	// Driver cards use the driverIdentification variant of CardNumber:
	//   - 14 bytes: driverIdentificationNumber (IA5String)
	//   - 1 byte: cardReplacementIndex (IA5String)
	//   - 1 byte: cardRenewalIndex (IA5String)
	driverID := &ddv1.DriverIdentification{}

	// Parse driver identification number (14 bytes)
	identificationNumber, err := opts.UnmarshalIa5StringValue(data[idxDriverIdentification : idxDriverIdentification+lenDriverIdentification])
	if err != nil {
		return nil, fmt.Errorf("failed to parse driver identification number: %w", err)
	}
	driverID.SetDriverIdentificationNumber(identificationNumber)

	// Parse replacement index (1 byte)
	replacementIndex, err := opts.UnmarshalIa5StringValue(data[idxReplacementIndex : idxReplacementIndex+lenReplacementIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card replacement index: %w", err)
	}
	driverID.SetCardReplacementIndex(replacementIndex)

	// Parse renewal index (1 byte)
	renewalIndex, err := opts.UnmarshalIa5StringValue(data[idxRenewalIndex : idxRenewalIndex+lenRenewalIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card renewal index: %w", err)
	}
	driverID.SetCardRenewalIndex(renewalIndex)

	id.SetDriverIdentification(driverID)

	// Authority name (36 bytes)
	authorityName, err := opts.UnmarshalStringValue(data[idxAuthorityName : idxAuthorityName+lenAuthorityName])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card issuing authority name: %w", err)
	}
	id.SetCardIssuingAuthorityName(authorityName)

	// Card issue date (4 bytes)
	cardIssueDate, err := opts.UnmarshalTimeReal(data[idxIssueDate : idxIssueDate+lenIssueDate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card issue date: %w", err)
	}
	id.SetCardIssueDate(cardIssueDate)

	// Card validity begin (4 bytes)
	cardValidityBegin, err := opts.UnmarshalTimeReal(data[idxValidityBegin : idxValidityBegin+lenValidityBegin])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card validity begin: %w", err)
	}
	id.SetCardValidityBegin(cardValidityBegin)

	// Card expiry date (4 bytes)
	cardExpiryDate, err := opts.UnmarshalTimeReal(data[idxExpiryDate : idxExpiryDate+lenExpiryDate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card expiry date: %w", err)
	}
	id.SetCardExpiryDate(cardExpiryDate)

	// Parse DriverCardHolderIdentification part (78 bytes)

	// Card holder surname (36 bytes)
	surname, err := opts.UnmarshalStringValue(data[idxSurname : idxSurname+lenSurname])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card holder surname: %w", err)
	}
	id.SetCardHolderSurname(surname)

	// Card holder first names (36 bytes)
	firstNames, err := opts.UnmarshalStringValue(data[idxFirstNames : idxFirstNames+lenFirstNames])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card holder first names: %w", err)
	}
	id.SetCardHolderFirstNames(firstNames)

	// Card holder birth date (4 bytes)
	birthDate, err := opts.UnmarshalDate(data[idxBirthDate : idxBirthDate+lenBirthDate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card holder birth date: %w", err)
	}
	id.SetCardHolderBirthDate(birthDate)

	// Card holder preferred language (2 bytes)
	preferredLanguage, err := opts.UnmarshalIa5StringValue(data[idxLanguage : idxLanguage+lenLanguage])
	if err != nil {
		return nil, fmt.Errorf("failed to parse card holder preferred language: %w", err)
	}
	id.SetCardHolderPreferredLanguage(preferredLanguage)

	return &id, nil
}

// MarshalDriverCardIdentification marshals the binary representation of DriverCardIdentification to bytes.
//
// The data type `DriverCardIdentification` combines two structures:
//   - CardIdentification (Data Dictionary Section 2.24) - 65 bytes
//   - DriverCardHolderIdentification (Data Dictionary Section 2.62) - 78 bytes
//
// Total binary size: 143 bytes (fixed)
//
// CardIdentification ASN.1 Specification (Data Dictionary Section 2.24):
//
//	CardIdentification ::= SEQUENCE {
//	    cardIssuingMemberState    NationNumeric,
//	    cardNumber                CardNumber,
//	    cardIssuingAuthorityName  Name,
//	    cardIssueDate             TimeReal,
//	    cardValidityBegin         TimeReal,
//	    cardExpiryDate            TimeReal
//	}
//
// DriverCardHolderIdentification ASN.1 Specification (Data Dictionary Section 2.62):
//
//	DriverCardHolderIdentification ::= SEQUENCE {
//	    cardHolderName                  HolderName,
//	    cardHolderBirthDate             Datef,
//	    cardHolderPreferredLanguage     Language
//	}
func (opts MarshalOptions) MarshalDriverCardIdentification(id *cardv1.DriverCardIdentification) ([]byte, error) {
	if id == nil {
		return nil, fmt.Errorf("driver card identification cannot be nil")
	}

	const lenDriverCardIdentification = 143
	dst := make([]byte, 0, lenDriverCardIdentification)

	// Marshal CardIdentification part (65 bytes)

	// Nation (1 byte)
	memberStateByte, err := dd.MarshalEnum(id.GetCardIssuingMemberState())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal member state: %w", err)
	}
	dst = append(dst, memberStateByte)

	// DriverIdentification (14 + 1 + 1 = 16 bytes)
	driverID := id.GetDriverIdentification()
	if driverID == nil {
		return nil, fmt.Errorf("driver identification cannot be nil")
	}

	// Driver identification number (14 bytes)
	identificationBytes, err := opts.MarshalIa5StringValue(driverID.GetDriverIdentificationNumber())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal driver identification number: %w", err)
	}
	dst = append(dst, identificationBytes...)

	// Card replacement index (1 byte)
	replacementBytes, err := opts.MarshalIa5StringValue(driverID.GetCardReplacementIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card replacement index: %w", err)
	}
	dst = append(dst, replacementBytes...)

	// Card renewal index (1 byte)
	renewalBytes, err := opts.MarshalIa5StringValue(driverID.GetCardRenewalIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card renewal index: %w", err)
	}
	dst = append(dst, renewalBytes...)

	// Authority name (36 bytes)
	authorityNameBytes, err := opts.MarshalStringValue(id.GetCardIssuingAuthorityName())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card issuing authority name: %w", err)
	}
	dst = append(dst, authorityNameBytes...)

	// Card issue date (4 bytes)
	issueDateBytes, err := opts.MarshalTimeReal(id.GetCardIssueDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card issue date: %w", err)
	}
	dst = append(dst, issueDateBytes...)

	// Card validity begin (4 bytes)
	validityBeginBytes, err := opts.MarshalTimeReal(id.GetCardValidityBegin())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card validity begin: %w", err)
	}
	dst = append(dst, validityBeginBytes...)

	// Card expiry date (4 bytes)
	expiryDateBytes, err := opts.MarshalTimeReal(id.GetCardExpiryDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card expiry date: %w", err)
	}
	dst = append(dst, expiryDateBytes...)

	// Marshal DriverCardHolderIdentification part (78 bytes)

	// Card holder surname (36 bytes)
	surnameBytes, err := opts.MarshalStringValue(id.GetCardHolderSurname())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card holder surname: %w", err)
	}
	dst = append(dst, surnameBytes...)

	// Card holder first names (36 bytes)
	firstNamesBytes, err := opts.MarshalStringValue(id.GetCardHolderFirstNames())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card holder first names: %w", err)
	}
	dst = append(dst, firstNamesBytes...)

	// Card holder birth date (4 bytes)
	birthDateBytes, err := opts.MarshalDate(id.GetCardHolderBirthDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card holder birth date: %w", err)
	}
	dst = append(dst, birthDateBytes...)

	// Card holder preferred language (2 bytes)
	languageBytes, err := opts.MarshalIa5StringValue(id.GetCardHolderPreferredLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card holder preferred language: %w", err)
	}
	dst = append(dst, languageBytes...)

	return dst, nil
}

// anonymizeDriverCardIdentification creates an anonymized copy of DriverCardIdentification,
// replacing all personally identifiable information with safe, deterministic test values while
// preserving the structure and validity for testing.
//
// Anonymization strategy:
//   - Names: Replaced with generic test names
//   - Driver identification: Replaced with test value
//   - Addresses: Replaced with generic test addresses
//   - Birth dates: Replaced with static test date (2000-01-01)
//   - Card dates: Replaced with static test dates (issue/validity: 2020-01-01, expiry: 2024-12-31)
//   - Countries: Preserved (structural info)
//   - Signatures: Cleared (will be invalid after anonymization anyway)
func (opts AnonymizeOptions) anonymizeDriverCardIdentification(id *cardv1.DriverCardIdentification) *cardv1.DriverCardIdentification {
	if id == nil {
		return nil
	}

	result := &cardv1.DriverCardIdentification{}

	// Preserve country (structural info)
	result.SetCardIssuingMemberState(id.GetCardIssuingMemberState())

	// Anonymize driver identification
	if driverID := id.GetDriverIdentification(); driverID != nil {
		result.SetDriverIdentification(dd.AnonymizeDriverIdentification(driverID))
	}

	// Anonymize issuing authority (ASCII-only to avoid encoding issues)
	authName := &ddv1.StringValue{}
	authName.SetValue("Transport and Communications Agency")
	authName.SetEncoding(ddv1.Encoding_ISO_8859_1)
	authName.SetLength(35) // Data length (not including code page byte)
	result.SetCardIssuingAuthorityName(authName)

	// Replace card dates with static test dates (valid 5-year period)
	// Issue/validity: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
	// Expiry: 2024-12-31 23:59:59 UTC (epoch: 1735689599)
	result.SetCardIssueDate(&timestamppb.Timestamp{Seconds: 1577836800})
	result.SetCardValidityBegin(&timestamppb.Timestamp{Seconds: 1577836800})
	result.SetCardExpiryDate(&timestamppb.Timestamp{Seconds: 1735689599})

	// Anonymize holder names (ASCII-only to avoid encoding issues)
	surname := &ddv1.StringValue{}
	surname.SetValue("Doe")
	surname.SetEncoding(ddv1.Encoding_ISO_8859_1)
	surname.SetLength(35)
	result.SetCardHolderSurname(surname)

	firstName := &ddv1.StringValue{}
	firstName.SetValue("John")
	firstName.SetEncoding(ddv1.Encoding_ISO_8859_1)
	firstName.SetLength(35)
	result.SetCardHolderFirstNames(firstName)

	// Replace birth date with static test date (2000-01-01)
	birthDate := &ddv1.Date{}
	birthDate.SetYear(2000)
	birthDate.SetMonth(1)
	birthDate.SetDay(1)
	// Regenerate raw_data for binary fidelity
	defOpts := dd.MarshalOptions{}
	if rawData, err := defOpts.MarshalDate(birthDate); err == nil {
		birthDate.SetRawData(rawData)
	}
	result.SetCardHolderBirthDate(birthDate)

	// Preserve language (not sensitive), but ensure it's always set with proper length
	language := id.GetCardHolderPreferredLanguage()
	if language == nil || !language.HasLength() {
		// Create a default language if missing (2 bytes for IA5String)
		language = &ddv1.Ia5StringValue{}
		language.SetValue("en")
		language.SetLength(2)
	}
	result.SetCardHolderPreferredLanguage(language)

	// Signature field left unset (nil) - TLV marshaller will omit the signature block

	return result
}
