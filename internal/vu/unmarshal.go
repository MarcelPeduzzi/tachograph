package vu

import (
	"github.com/way-platform/tachograph-go/internal/dd"
)

// UnmarshalOptions provides context for parsing binary VU data.
//
// This struct embeds dd.UnmarshalOptions to inherit all data dictionary unmarshal methods,
// and extends it with VU-specific unmarshal methods.
//
// See also: tachograph.UnmarshalOptions for the public API definition.
type UnmarshalOptions struct {
	// Embed dd.UnmarshalOptions to inherit all data dictionary unmarshal methods.
	// This allows vu.UnmarshalOptions to be used wherever dd.UnmarshalOptions is needed.
	//
	// Inherited fields:
	// - PreserveRawData: controls whether raw byte slices are stored
	dd.UnmarshalOptions

	// Strict controls how the parser handles unrecognized transfer types.
	//
	// If true (default), the parser will return an error on any unrecognized
	// transfer types or tags.
	// If false, the parser will skip over unrecognized transfers and continue parsing.
	Strict bool
}
