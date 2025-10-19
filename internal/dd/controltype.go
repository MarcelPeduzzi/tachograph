package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalControlType unmarshals a control type from a byte slice
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//	ControlType ::= OCTET STRING (SIZE(1))
func (opts UnmarshalOptions) UnmarshalControlType(data []byte) (*ddv1.ControlType, error) {
	const lenControlType = 1
	if len(data) != lenControlType {
		return nil, fmt.Errorf("invalid data length for ControlType: got %d, want %d", len(data), lenControlType)
	}

	output := &ddv1.ControlType{}
	if opts.PreserveRawData {
		output.SetRawData(data[:lenControlType])
	}
	b := data[0]
	output.SetCardDownloading((b & 0x80) != 0)
	output.SetVuDownloading((b & 0x40) != 0)
	output.SetPrinting((b & 0x20) != 0)
	output.SetDisplay((b & 0x10) != 0)
	output.SetCalibrationChecking((b & 0x08) != 0)

	return output, nil
}

// MarshalControlType marshals a ControlType as a single byte bitmask.
//
// The data type `ControlType` is specified in the Data Dictionary, Section 2.53.
//
// ASN.1 Definition:
//
//	ControlType ::= OCTET STRING (SIZE(1))
//
// Binary Layout (1 byte):
//   - Bit 7: card downloading
//   - Bit 6: VU downloading
//   - Bit 5: printing
//   - Bit 4: display
//   - Bit 3: calibration checking (Gen2+)
//   - Bits 2-0: Reserved (RFU)
func (opts MarshalOptions) MarshalControlType(controlType *ddv1.ControlType) ([]byte, error) {
	const lenControlType = 1
	var canvas [lenControlType]byte
	if controlType.HasRawData() {
		if len(controlType.GetRawData()) != lenControlType {
			return nil, fmt.Errorf(
				"invalid raw_data length for ControlType: got %d, want %d",
				len(controlType.GetRawData()), lenControlType,
			)
		}
		copy(canvas[:], controlType.GetRawData())
	}
	// Clear known bits, then set based on semantic fields (raw data painting strategy)
	canvas[0] &= 0x07 // Clear bits 7-3 (known bits), preserve bits 2-0 (reserved)
	if controlType.GetCardDownloading() {
		canvas[0] |= 0x80
	}
	if controlType.GetVuDownloading() {
		canvas[0] |= 0x40
	}
	if controlType.GetPrinting() {
		canvas[0] |= 0x20
	}
	if controlType.GetDisplay() {
		canvas[0] |= 0x10
	}
	if controlType.GetCalibrationChecking() {
		canvas[0] |= 0x08
	}
	return canvas[:], nil
}
