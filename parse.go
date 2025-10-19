package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

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

// Parse performs the second parsing pass, converting raw records into semantic
// data structures. If the raw file has been authenticated (via
// AuthenticateOptions.Authenticate), the authentication results are propagated
// to the parsed messages.
//
// The zero value of ParseOptions preserves raw data for round-trip fidelity.
func (o ParseOptions) Parse(rawFile *tachographv1.RawFile) (*tachographv1.File, error) {
	// Apply defaults
	if o == (ParseOptions{}) {
		o.PreserveRawData = true
	}

	var file tachographv1.File

	switch rawFile.GetType() {
	case tachographv1.RawFile_CARD:
		cardType := card.InferFileType(rawFile.GetCard())
		switch cardType {
		case cardv1.CardType_DRIVER_CARD:
			cardOpts := card.ParseOptions{
				PreserveRawData: o.PreserveRawData,
			}
			driverCard, err := cardOpts.ParseRawDriverCardFile(rawFile.GetCard())
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver card: %w", err)
			}
			file.SetType(tachographv1.File_DRIVER_CARD)
			file.SetDriverCard(driverCard)
		default:
			return nil, fmt.Errorf("unsupported card type: %v", cardType)
		}

	case tachographv1.RawFile_VEHICLE_UNIT:
		vuOpts := vu.ParseOptions{
			PreserveRawData: o.PreserveRawData,
		}
		vuFile, err := vuOpts.ParseRawVehicleUnitFile(rawFile.GetVehicleUnit())
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
