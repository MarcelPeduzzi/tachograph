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
