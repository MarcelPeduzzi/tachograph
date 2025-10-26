package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalSoftwareIdentification parses the SoftwareIdentification structure.
//
// See Data Dictionary, Section 2.225, `VuSoftwareIdentification`.
//
// ASN.1 Specification:
//
//	VuSoftwareIdentification ::= SEQUENCE {
//	    vuSoftwareVersion VuSoftwareVersion,           -- 4 bytes (IA5String, no code page)
//	    vuSoftInstallationDate VuSoftInstallationDate  -- 4 bytes (TimeReal)
//	}
//
// Binary Layout (fixed length: 8 bytes):
//   - VU Software Version (4 bytes): IA5String (no code page byte)
//   - VU Software Installation Date (4 bytes): TimeReal
func (opts UnmarshalOptions) UnmarshalSoftwareIdentification(data []byte) (*ddv1.SoftwareIdentification, error) {
	const lenSoftwareIdentification = 8

	if len(data) != lenSoftwareIdentification {
		return nil, fmt.Errorf(
			"invalid data length for SoftwareIdentification: got %d, want %d",
			len(data), lenSoftwareIdentification,
		)
	}

	softwareIdent := &ddv1.SoftwareIdentification{}

	const (
		idxSoftwareVersion  = 0
		lenSoftwareVersion  = 4
		idxInstallationDate = 4
		lenInstallationDate = 4
	)

	// Parse software version (4 bytes, IA5String)
	softwareVersion, err := opts.UnmarshalIa5StringValue(data[idxSoftwareVersion : idxSoftwareVersion+lenSoftwareVersion])
	if err != nil {
		return nil, fmt.Errorf("failed to parse software version: %w", err)
	}
	softwareIdent.SetSoftwareVersion(softwareVersion)

	// Parse installation date (4 bytes)
	installationDate, err := opts.UnmarshalTimeReal(data[idxInstallationDate : idxInstallationDate+lenInstallationDate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse software installation date: %w", err)
	}
	softwareIdent.SetSoftwareInstallationDate(installationDate)

	return softwareIdent, nil
}

// MarshalSoftwareIdentification marshals the SoftwareIdentification structure.
//
// See Data Dictionary, Section 2.225, `VuSoftwareIdentification`.
func (opts MarshalOptions) MarshalSoftwareIdentification(softwareIdent *ddv1.SoftwareIdentification) ([]byte, error) {
	if softwareIdent == nil {
		return nil, fmt.Errorf("softwareIdent cannot be nil")
	}

	const size = 8
	var canvas [size]byte

	offset := 0

	// Marshal software version (4 bytes, IA5String)
	softwareVersionBytes, err := opts.MarshalIa5StringValue(softwareIdent.GetSoftwareVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal software version: %w", err)
	}
	if len(softwareVersionBytes) != 4 {
		return nil, fmt.Errorf(
			"invalid software version length: got %d, want 4",
			len(softwareVersionBytes),
		)
	}
	copy(canvas[offset:offset+4], softwareVersionBytes)
	offset += 4

	// Marshal installation date (4 bytes)
	installationDateBytes, err := opts.MarshalTimeReal(softwareIdent.GetSoftwareInstallationDate())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal software installation date: %w", err)
	}
	if len(installationDateBytes) != 4 {
		return nil, fmt.Errorf(
			"invalid software installation date length: got %d, want 4",
			len(installationDateBytes),
		)
	}
	copy(canvas[offset:offset+4], installationDateBytes)
	offset += 4

	if offset != size {
		return nil, fmt.Errorf(
			"SoftwareIdentification marshalling size mismatch: wrote %d bytes, expected %d",
			offset, size,
		)
	}

	return canvas[:], nil
}
