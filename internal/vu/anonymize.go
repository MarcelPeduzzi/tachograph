package vu

import (
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// AnonymizeOptions configures the anonymization of VU files.
type AnonymizeOptions struct {
	// PreserveDistanceAndTrips controls whether distance and trip data are preserved.
	PreserveDistanceAndTrips bool

	// PreserveTimestamps controls whether timestamps are preserved.
	PreserveTimestamps bool
}

// AnonymizeVehicleUnitFile creates an anonymized copy of a vehicle unit file.
func (opts AnonymizeOptions) AnonymizeVehicleUnitFile(file *vuv1.VehicleUnitFile) (*vuv1.VehicleUnitFile, error) {
	if file == nil {
		return nil, nil
	}

	// Clone the file to avoid mutating the input
	result := proto.Clone(file).(*vuv1.VehicleUnitFile)

	// Anonymize based on generation
	switch result.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		gen1 := result.GetGen1()
		if gen1 == nil {
			return result, nil
		}

		// Anonymize Overview
		if overview := gen1.GetOverview(); overview != nil {
			gen1.SetOverview(opts.anonymizeOverviewGen1(overview))
		}

		// Anonymize Activities (multiple transfers)
		var anonymizedActivities []*vuv1.ActivitiesGen1
		for _, activities := range gen1.GetActivities() {
			anonymizedActivities = append(anonymizedActivities, opts.anonymizeActivitiesGen1(activities))
		}
		gen1.SetActivities(anonymizedActivities)

		// Anonymize Events and Faults (multiple transfers)
		var anonymizedEventsAndFaults []*vuv1.EventsAndFaultsGen1
		for _, ef := range gen1.GetEventsAndFaults() {
			anonymizedEventsAndFaults = append(anonymizedEventsAndFaults, opts.anonymizeEventsAndFaultsGen1(ef))
		}
		gen1.SetEventsAndFaults(anonymizedEventsAndFaults)

		// Anonymize Detailed Speed (multiple transfers)
		var anonymizedDetailedSpeed []*vuv1.DetailedSpeedGen1
		for _, ds := range gen1.GetDetailedSpeed() {
			anonymizedDetailedSpeed = append(anonymizedDetailedSpeed, opts.anonymizeDetailedSpeedGen1(ds))
		}
		gen1.SetDetailedSpeed(anonymizedDetailedSpeed)

		// Anonymize Technical Data (multiple transfers)
		var anonymizedTechnicalData []*vuv1.TechnicalDataGen1
		for _, td := range gen1.GetTechnicalData() {
			anonymizedTechnicalData = append(anonymizedTechnicalData, opts.anonymizeTechnicalDataGen1(td))
		}
		gen1.SetTechnicalData(anonymizedTechnicalData)

	case ddv1.Generation_GENERATION_2:
		if result.GetVersion() == ddv1.Version_VERSION_2 {
			// Handle Gen2 V2
			gen2v2 := result.GetGen2V2()
			if gen2v2 == nil {
				return result, nil
			}

			// Anonymize Overview
			if overview := gen2v2.GetOverview(); overview != nil {
				gen2v2.SetOverview(opts.anonymizeOverviewGen2V2(overview))
			}

			// Anonymize Activities (multiple transfers)
			var anonymizedActivities []*vuv1.ActivitiesGen2V2
			for _, activities := range gen2v2.GetActivities() {
				anonymizedActivities = append(anonymizedActivities, opts.anonymizeActivitiesGen2V2(activities))
			}
			gen2v2.SetActivities(anonymizedActivities)

			// Anonymize Events and Faults (multiple transfers)
			var anonymizedEventsAndFaults []*vuv1.EventsAndFaultsGen2V2
			for _, ef := range gen2v2.GetEventsAndFaults() {
				anonymizedEventsAndFaults = append(anonymizedEventsAndFaults, opts.anonymizeEventsAndFaultsGen2V2(ef))
			}
			gen2v2.SetEventsAndFaults(anonymizedEventsAndFaults)

			// Anonymize Detailed Speed (multiple transfers)
			var anonymizedDetailedSpeed []*vuv1.DetailedSpeedGen2
			for _, ds := range gen2v2.GetDetailedSpeed() {
				anonymizedDetailedSpeed = append(anonymizedDetailedSpeed, opts.anonymizeDetailedSpeedGen2(ds))
			}
			gen2v2.SetDetailedSpeed(anonymizedDetailedSpeed)

			// Anonymize Technical Data (multiple transfers)
			var anonymizedTechnicalData []*vuv1.TechnicalDataGen2V2
			for _, td := range gen2v2.GetTechnicalData() {
				anonymizedTechnicalData = append(anonymizedTechnicalData, opts.anonymizeTechnicalDataGen2V2(td))
			}
			gen2v2.SetTechnicalData(anonymizedTechnicalData)

		} else {
			// Handle Gen2 V1
			gen2v1 := result.GetGen2V1()
			if gen2v1 == nil {
				return result, nil
			}

			// Anonymize Overview
			if overview := gen2v1.GetOverview(); overview != nil {
				gen2v1.SetOverview(opts.anonymizeOverviewGen2V1(overview))
			}

			// Anonymize Activities (multiple transfers)
			var anonymizedActivities []*vuv1.ActivitiesGen2V1
			for _, activities := range gen2v1.GetActivities() {
				anonymizedActivities = append(anonymizedActivities, opts.anonymizeActivitiesGen2V1(activities))
			}
			gen2v1.SetActivities(anonymizedActivities)

			// Anonymize Events and Faults (multiple transfers)
			var anonymizedEventsAndFaults []*vuv1.EventsAndFaultsGen2V1
			for _, ef := range gen2v1.GetEventsAndFaults() {
				anonymizedEventsAndFaults = append(anonymizedEventsAndFaults, opts.anonymizeEventsAndFaultsGen2V1(ef))
			}
			gen2v1.SetEventsAndFaults(anonymizedEventsAndFaults)

			// Anonymize Detailed Speed (multiple transfers)
			var anonymizedDetailedSpeed []*vuv1.DetailedSpeedGen2
			for _, ds := range gen2v1.GetDetailedSpeed() {
				anonymizedDetailedSpeed = append(anonymizedDetailedSpeed, opts.anonymizeDetailedSpeedGen2(ds))
			}
			gen2v1.SetDetailedSpeed(anonymizedDetailedSpeed)

			// Anonymize Technical Data (multiple transfers)
			var anonymizedTechnicalData []*vuv1.TechnicalDataGen2V1
			for _, td := range gen2v1.GetTechnicalData() {
				anonymizedTechnicalData = append(anonymizedTechnicalData, opts.anonymizeTechnicalDataGen2V1(td))
			}
			gen2v1.SetTechnicalData(anonymizedTechnicalData)
		}
	}

	return result, nil
}
