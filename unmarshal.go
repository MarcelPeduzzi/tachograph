package tachograph

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// UnmarshalRawFile performs the first parsing pass on a tachograph file,
// identifying record boundaries and preserving raw byte values.
// The returned RawFile is suitable for signature authentication via AuthenticateRawFile.
func UnmarshalRawFile(data []byte) (*tachographv1.RawFile, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for tachograph file: %w", io.ErrUnexpectedEOF)
	}

	var rawFile tachographv1.RawFile

	switch {
	// Vehicle unit file (starts with TREP prefix 0x76).
	case data[0] == 0x76:
		vuRaw, err := vu.UnmarshalRawVehicleUnitFile(data)
		if err != nil {
			return nil, err
		}
		rawFile.SetType(tachographv1.RawFile_VEHICLE_UNIT)
		rawFile.SetVehicleUnit(vuRaw)

	// Card file (starts with EF_ICC prefix 0x0002).
	case binary.BigEndian.Uint16(data[0:2]) == 0x0002:
		cardRaw, err := card.UnmarshalRawCardFile(data)
		if err != nil {
			return nil, err
		}
		rawFile.SetType(tachographv1.RawFile_CARD)
		rawFile.SetCard(cardRaw)

	default:
		return nil, errors.New("unknown or unsupported file type")
	}

	return &rawFile, nil
}

// ParseRawFile performs the second parsing pass, converting raw records
// into semantic data structures. If the raw file has been authenticated
// (via AuthenticateRawFile), the authentication results are propagated to
// the parsed messages.
func ParseRawFile(rawFile *tachographv1.RawFile) (*tachographv1.File, error) {
	var file tachographv1.File

	switch rawFile.GetType() {
	case tachographv1.RawFile_CARD:
		cardType := card.InferFileType(rawFile.GetCard())
		switch cardType {
		case cardv1.CardType_DRIVER_CARD:
			driverCard, err := card.ParseRawDriverCardFile(rawFile.GetCard())
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver card: %w", err)
			}
			file.SetType(tachographv1.File_DRIVER_CARD)
			file.SetDriverCard(driverCard)
		default:
			return nil, fmt.Errorf("unsupported card type: %v", cardType)
		}

	case tachographv1.RawFile_VEHICLE_UNIT:
		vuFile, err := vu.ParseRawVehicleUnitFile(rawFile.GetVehicleUnit())
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
