package tachograph

import (
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// AnonymizeOptions configures the anonymization process for tachograph files.
type AnonymizeOptions struct {
	// PreserveDistanceAndTrips controls whether distance and trip data are preserved.
	//
	// If true, odometer readings and distance values are preserved in their original form.
	// If false (default), distance data is rounded or anonymized to obscure exact values.
	PreserveDistanceAndTrips bool

	// PreserveTimestamps controls whether timestamps are preserved.
	//
	// If true, timestamps are preserved in their original form.
	// If false (default), timestamps are shifted to a fixed epoch (2020-01-01 00:00:00 UTC)
	// to obscure the exact time of events while maintaining relative ordering.
	PreserveTimestamps bool
}

// Anonymize creates an anonymized copy of a parsed tachograph file.
//
// Anonymization replaces personally identifiable information (PII) with test values
// while preserving the structural integrity of the file for testing purposes.
//
// The zero value of AnonymizeOptions anonymizes both timestamps and distances.
func (o AnonymizeOptions) Anonymize(file *tachographv1.File) (*tachographv1.File, error) {
	if file == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}

	var result tachographv1.File
	result.SetType(file.GetType())

	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		cardOpts := card.AnonymizeOptions{
			PreserveDistanceAndTrips: o.PreserveDistanceAndTrips,
			PreserveTimestamps:       o.PreserveTimestamps,
		}
		anonymizedCard, err := cardOpts.AnonymizeDriverCardFile(file.GetDriverCard())
		if err != nil {
			return nil, fmt.Errorf("failed to anonymize driver card: %w", err)
		}
		result.SetDriverCard(anonymizedCard)

	case tachographv1.File_VEHICLE_UNIT:
		vuOpts := vu.AnonymizeOptions{
			PreserveDistanceAndTrips: o.PreserveDistanceAndTrips,
			PreserveTimestamps:       o.PreserveTimestamps,
		}
		anonymizedVU, err := vuOpts.AnonymizeVehicleUnitFile(file.GetVehicleUnit())
		if err != nil {
			return nil, fmt.Errorf("failed to anonymize vehicle unit: %w", err)
		}
		result.SetVehicleUnit(anonymizedVU)

	default:
		return nil, fmt.Errorf("unsupported file type for anonymization: %v", file.GetType())
	}

	return &result, nil
}
