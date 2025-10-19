package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/dd"
	"github.com/way-platform/tachograph-go/internal/vu"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// Marshal serializes a parsed tachograph file into binary format with default options.
//
// This is a convenience function that uses default options:
// - UseRawData: true (use raw data painting for perfect round-trip fidelity)
//
// For custom options, use MarshalOptions directly:
//
//	opts := MarshalOptions{UseRawData: false}
//	data, err := opts.Marshal(file)
func Marshal(file *tachographv1.File) ([]byte, error) {
	opts := MarshalOptions{
		UseRawData: true,
	}
	return opts.Marshal(file)
}

// MarshalOptions configures the marshaling process for tachograph files.
type MarshalOptions struct {
	// UseRawData controls whether the marshaler uses raw_data fields from
	// parsed messages to reconstruct the file.
	//
	// If true (default), the marshaler will use the raw_data fields when
	// available, applying the "raw data painting" strategy to ensure perfect
	// binary round-tripping while validating semantic field correctness.
	//
	// If false, the marshaler will always encode from semantic fields,
	// ignoring any raw_data fields. This is useful when semantic fields
	// have been modified and you want to generate new binary data.
	UseRawData bool
}

// Marshal serializes a parsed tachograph file into its binary representation.
//
// Use NewMarshalOptions() to create options with sensible defaults, or
// construct MarshalOptions directly to customize behavior.
func (o MarshalOptions) Marshal(file *tachographv1.File) ([]byte, error) {
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		cardOpts := card.MarshalOptions{
			MarshalOptions: dd.MarshalOptions{
				UseRawData: o.UseRawData,
			},
		}
		return cardOpts.MarshalDriverCardFile(file.GetDriverCard())
	case tachographv1.File_VEHICLE_UNIT:
		vuOpts := vu.MarshalOptions{
			MarshalOptions: dd.MarshalOptions{
				UseRawData: o.UseRawData,
			},
		}
		return vuOpts.MarshalVehicleUnitFile(file.GetVehicleUnit())
	default:
		return nil, fmt.Errorf("unsupported file type for marshaling: %v", file.GetType())
	}
}
