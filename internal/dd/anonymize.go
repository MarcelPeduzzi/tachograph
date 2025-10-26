package dd

import (
	"time"
)

// AnonymizeOptions configures anonymization behavior for DD-level helpers.
type AnonymizeOptions struct {
	PreserveDistanceAndTrips bool
	PreserveTimestamps       bool
	TimestampEpoch           time.Time // Base epoch for relative timestamp shifts
}

// DefaultTimestampEpoch is the default epoch for timestamp anonymization (2020-01-01 00:00:00 UTC).
var DefaultTimestampEpoch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
