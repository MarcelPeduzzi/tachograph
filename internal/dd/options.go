package dd

// UnmarshalOptions provides context for parsing binary tachograph data.
//
// This struct follows the pattern used in protojson.UnmarshalOptions and other
// Go standard library packages, where unmarshal functions are methods on the
// options struct.
//
// See also: tachograph.UnmarshalOptions for the public API definition.
type UnmarshalOptions struct {
	// PreserveRawData controls whether raw byte slices are stored in
	// the raw_data field of parsed protobuf messages.
	//
	// If true, the raw byte slice for each parsed element will be stored
	// in the raw_data field, enabling perfect binary round-tripping via
	// the raw data painting strategy in Marshal.
	//
	// If false, raw_data fields will be left empty, reducing memory usage
	// but preventing exact binary reconstruction.
	PreserveRawData bool
}

// MarshalOptions provides context for marshaling binary tachograph data.
type MarshalOptions struct {
	// UseRawData controls whether the raw_data field is used during marshaling.
	//
	// If true (default), the "raw data painting" strategy is used: the raw_data
	// field serves as a canvas, and semantic fields are painted over it at their
	// designated byte offsets. This preserves unknown bits and reserved fields.
	//
	// If false, data is always encoded from semantic fields, ignoring raw_data.
	//
	// NOTE: This option is currently always treated as true. Full implementation
	// of the UseRawData=false behavior is deferred to a future phase.
	UseRawData bool
}
