package card

import (
	"encoding/binary"

	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// MarshalOptions configures the marshaling of card files into binary format.
type MarshalOptions struct {
	// Embed dd.MarshalOptions to inherit marshaling configuration.
	dd.MarshalOptions
}

// MarshalRawCardFile serializes a RawCardFile into binary format.
func (opts MarshalOptions) MarshalRawCardFile(file *cardv1.RawCardFile) ([]byte, error) {
	var result []byte
	for _, record := range file.GetRecords() {
		// Write tag (FID + appendix)
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetTag()>>8))
		result = append(result, byte(record.GetTag()&0xFF))
		// Write length
		result = binary.BigEndian.AppendUint16(result, uint16(record.GetLength()))
		// Write value
		result = append(result, record.GetValue()...)
	}
	return result, nil
}
