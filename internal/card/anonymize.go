package card

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
)

// AnonymizeDriverCardFile creates an anonymized copy of a driver card file.
//
// Anonymization replaces personally identifiable information (PII) with test values
// while preserving the structural integrity of the file for testing purposes.
//
// Current Implementation (Phase 1):
//
// For now, this function returns a deep clone of the file. Full PII anonymization
// logic (replacing names, addresses, license numbers, timestamps, odometer data)
// is deferred to Phase 2, pending business requirements for specific anonymization rules.
//
// Future Implementation (Phase 2):
//
// When fully implemented, this function will:
//   - Replace driver names and personal information with test values
//   - Replace addresses and company names with test values
//   - Optionally shift timestamps (if PreserveTimestamps is false)
//   - Optionally anonymize odometer/distance data (if PreserveDistanceAndTrips is false)
//   - Anonymize driving licence information
//   - Preserve file structure and size for testing compatibility
//
// The options provided (o.PreserveDistanceAndTrips, o.PreserveTimestamps) are available
// for use by the future implementation.
func (opts AnonymizeOptions) AnonymizeDriverCardFile(file *cardv1.DriverCardFile) (*cardv1.DriverCardFile, error) {
	if file == nil {
		return nil, fmt.Errorf("driver card file cannot be nil")
	}

	// Phase 1: Return a deep clone. This preserves structure for round-trip testing.
	// Phase 2: Implement actual PII anonymization using opts.
	return proto.Clone(file).(*cardv1.DriverCardFile), nil
}
