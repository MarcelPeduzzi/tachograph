package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// AnonymizeVehicleUnitFile creates an anonymized copy of a vehicle unit file.
//
// Anonymization replaces personally identifiable information (PII) with test values
// while preserving the structural integrity of the file for testing purposes.
//
// Current Implementation (Phase 1):
//
// For now, this function returns a deep clone of the file. Full PII anonymization
// logic (replacing vehicle identification, certificate data, timestamp shifting)
// is deferred to Phase 2, pending business requirements for specific anonymization rules.
//
// Future Implementation (Phase 2):
//
// When fully implemented, this function will:
//   - Replace vehicle identification numbers (VIN) with test values
//   - Replace vehicle registration information with test values
//   - Optionally shift timestamps (if PreserveTimestamps is false)
//   - Optionally anonymize odometer/distance data (if PreserveDistanceAndTrips is false)
//   - Anonymize certificate and signature data
//   - Preserve file structure and record counts for testing compatibility
//
// The options provided (o.PreserveDistanceAndTrips, o.PreserveTimestamps) are available
// for use by the future implementation.
func (opts AnonymizeOptions) AnonymizeVehicleUnitFile(file *vuv1.VehicleUnitFile) (*vuv1.VehicleUnitFile, error) {
	if file == nil {
		return nil, fmt.Errorf("vehicle unit file cannot be nil")
	}

	// Phase 1: Return a deep clone. This preserves structure for round-trip testing.
	// Phase 2: Implement actual PII anonymization using opts.
	return proto.Clone(file).(*vuv1.VehicleUnitFile), nil
}
