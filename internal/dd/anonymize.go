package dd

// AnonymizeOptions configures the anonymization of data dictionary types.
type AnonymizeOptions struct {
	// PreserveTimestamps controls whether timestamps are preserved.
	//
	// If true, timestamps are preserved in their original form.
	// If false (default), timestamps are shifted to a fixed epoch (2020-01-01 00:00:00 UTC)
	// to obscure the exact time of events while maintaining relative ordering.
	PreserveTimestamps bool

	// PreserveDistanceAndTrips controls whether distance and trip data are preserved.
	//
	// If true, odometer readings and distance values are preserved in their original form.
	// If false (default), distance data is rounded or anonymized to obscure exact values.
	PreserveDistanceAndTrips bool
}
