package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// Unparse converts a parsed File back into its raw representation.
// This is the inverse of Parse.
//
// Unparse takes a semantic File (with parsed data structures) and converts it
// back into a RawFile (with binary TLV/TV records). This is useful for:
// - Round-trip testing (Parse → Unparse → Parse should be identical)
// - Creating modified files (Parse → modify → Unparse → Marshal)
// - Generating test data from semantic structures
func Unparse(file *tachographv1.File) (*tachographv1.RawFile, error) {
	if file == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}

	var result tachographv1.RawFile

	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		rawCard, err := card.UnparseDriverCardFile(file.GetDriverCard())
		if err != nil {
			return nil, fmt.Errorf("failed to unparse driver card: %w", err)
		}
		result.SetType(tachographv1.RawFile_CARD)
		result.SetCard(rawCard)

	case tachographv1.File_VEHICLE_UNIT:
		rawVU, err := vu.UnparseVehicleUnitFile(file.GetVehicleUnit())
		if err != nil {
			return nil, fmt.Errorf("failed to unparse vehicle unit: %w", err)
		}
		result.SetType(tachographv1.RawFile_VEHICLE_UNIT)
		result.SetVehicleUnit(rawVU)

	default:
		return nil, fmt.Errorf("unsupported file type: %v", file.GetType())
	}

	return &result, nil
}
