package vu

import (
	"github.com/way-platform/tachograph-go/internal/dd"
)

// ParseOptions configures the parsing of raw VU files into semantic structures.
type ParseOptions struct {
	// PreserveRawData controls whether raw byte slices are stored in
	// the raw_data field of parsed protobuf messages.
	PreserveRawData bool
}

// unmarshal returns UnmarshalOptions configured from ParseOptions.
func (o ParseOptions) unmarshal() UnmarshalOptions {
	return UnmarshalOptions{
		UnmarshalOptions: dd.UnmarshalOptions{
			PreserveRawData: o.PreserveRawData,
		},
	}
}
