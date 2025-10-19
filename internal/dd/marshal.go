package dd

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
