package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// Parse performs semantic parsing on raw tachograph records with default options.
// If the raw file has been authenticated (via Authenticate), the authentication
// results are propagated to the parsed messages.
//
// This is a convenience function that uses default options:
// - PreserveRawData: true (store raw bytes for round-tripping)
//
// For custom options, use ParseOptions directly:
//
//	opts := ParseOptions{PreserveRawData: false}
//	file, err := opts.Parse(rawFile)
func Parse(rawFile *tachographv1.RawFile) (*tachographv1.File, error) {
	opts := ParseOptions{
		PreserveRawData: true,
	}
	return opts.Parse(rawFile)
}

// ParseOptions configures the parsing process for converting raw tachograph
// files into semantic data structures.
type ParseOptions struct {
	// PreserveRawData controls whether raw byte slices are stored in
	// the raw_data field of parsed protobuf messages.
	//
	// If true (default), the raw byte slice for each parsed element
	// will be stored in the raw_data field, enabling perfect binary
	// round-tripping via Marshal.
	//
	// If false, raw_data fields will be left empty, reducing memory usage
	// but preventing exact binary reconstruction.
	PreserveRawData bool
}

// card returns card.ParseOptions configured from ParseOptions.
func (o ParseOptions) card() card.ParseOptions {
	return card.ParseOptions{
		PreserveRawData: o.PreserveRawData,
	}
}

// vu returns vu.ParseOptions configured from ParseOptions.
func (o ParseOptions) vu() vu.ParseOptions {
	return vu.ParseOptions{
		PreserveRawData: o.PreserveRawData,
	}
}

// Parse performs the second parsing pass, converting raw records into semantic
// data structures. If the raw file has been authenticated (via
// AuthenticateOptions.Authenticate), the authentication results are propagated
// to the parsed messages.
func (o ParseOptions) Parse(rawFile *tachographv1.RawFile) (*tachographv1.File, error) {
	var file tachographv1.File

	switch rawFile.GetType() {
	case tachographv1.RawFile_CARD:
		cardType := card.InferFileType(rawFile.GetCard())
		switch cardType {
		case cardv1.CardType_DRIVER_CARD:
			driverCard, err := o.card().ParseRawDriverCardFile(rawFile.GetCard())
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver card: %w", err)
			}
			file.SetType(tachographv1.File_DRIVER_CARD)
			file.SetDriverCard(driverCard)
		default:
			return nil, fmt.Errorf("unsupported card type: %v", cardType)
		}

	case tachographv1.RawFile_VEHICLE_UNIT:
		vuFile, err := o.vu().ParseRawVehicleUnitFile(rawFile.GetVehicleUnit())
		if err != nil {
			return nil, err
		}
		file.SetType(tachographv1.File_VEHICLE_UNIT)
		file.SetVehicleUnit(vuFile)

	default:
		return nil, fmt.Errorf("unknown raw file type: %v", rawFile.GetType())
	}

	return &file, nil
}
