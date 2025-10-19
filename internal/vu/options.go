package vu

import (
	"github.com/way-platform/tachograph-go/internal/dd"
)

// UnmarshalOptions provides context for parsing binary VU data.
//
// This struct embeds dd.UnmarshalOptions to inherit generation/version context,
// and extends it with VU-specific unmarshal configuration.
//
// The zero value (UnmarshalOptions{}) is valid and represents the default
// parsing behavior for vehicle unit data.
type UnmarshalOptions struct {
	// Embed dd.UnmarshalOptions to inherit generation/version context.
	dd.UnmarshalOptions

	// Strict controls how the parser handles unrecognized transfer types.
	//
	// If true (default), the parser will return an error on any unrecognized
	// transfer types or tags.
	// If false, the parser will skip over unrecognized transfers and continue parsing.
	Strict bool
}

// ParseOptions configures the parsing of raw VU files into semantic structures.
type ParseOptions struct {
	// PreserveRawData controls whether raw byte slices are stored in
	// the raw_data field of parsed protobuf messages.
	PreserveRawData bool
}

// MarshalOptions configures the marshaling of VU files into binary format.
type MarshalOptions struct {
	// Embed dd.MarshalOptions to inherit marshaling configuration.
	dd.MarshalOptions
}

// AnonymizeOptions configures the anonymization of VU files.
type AnonymizeOptions struct {
	// PreserveDistanceAndTrips controls whether distance and trip data are preserved.
	PreserveDistanceAndTrips bool

	// PreserveTimestamps controls whether timestamps are preserved.
	PreserveTimestamps bool
}
