package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalMonthYear parses BCD-encoded month and year data.
//
// The data type `monthYear` is specified in the Data Dictionary, Section 2.72
// as part of ExtendedSerialNumber.
//
// ASN.1 Definition:
//
//	monthYear BCDString(SIZE(2))
//
// Binary Layout (2 bytes):
//   - BCD-encoded MMYY format (2 bytes)
func (opts UnmarshalOptions) UnmarshalMonthYear(data []byte) (*ddv1.MonthYear, error) {
	const lenMonthYear = 2

	if len(data) != lenMonthYear {
		return nil, fmt.Errorf("invalid data length for MonthYear: got %d, want %d", len(data), lenMonthYear)
	}

	monthYear := &ddv1.MonthYear{}
	if opts.PreserveRawData {
		monthYear.SetRawData(data[:lenMonthYear])
	}

	// Decode BCD month/year as 4-digit number MMYY
	monthYearInt, err := decodeBCD(data[:lenMonthYear])
	if err == nil && monthYearInt > 0 {
		month := int32(monthYearInt / 100)
		year := int32(monthYearInt % 100)

		// Convert 2-digit year to 4-digit (assuming 20xx for years 00-99)
		if year >= 0 && year <= 99 {
			year += 2000
		}

		monthYear.SetMonth(month)
		monthYear.SetYear(year)
	}

	return monthYear, nil
}

// MarshalMonthYear marshals a 2-byte BCD-encoded month/year value.
//
// The data type `MonthYear` is specified in the Data Dictionary.
//
// Binary Layout (2 bytes):
//   - Month (1 byte): BCD-encoded MM
//   - Year (1 byte): BCD-encoded YY (last 2 digits)
func (opts MarshalOptions) MarshalMonthYear(my *ddv1.MonthYear) ([]byte, error) {
	const lenMonthYear = 2
	var canvas [lenMonthYear]byte
	if my.HasRawData() {
		if len(my.GetRawData()) != lenMonthYear {
			return nil, fmt.Errorf(
				"invalid raw_data length for MonthYear: got %d, want %d",
				len(my.GetRawData()), lenMonthYear,
			)
		}
		copy(canvas[:], my.GetRawData())
	}
	month := int(my.GetMonth())
	year := int(my.GetYear())
	canvas[0] = byte((month/10)%10<<4 | month%10)
	canvas[1] = byte((year/10)%10<<4 | year%10)
	return canvas[:], nil
}
