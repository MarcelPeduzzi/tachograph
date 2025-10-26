package vu

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// UnparseVehicleUnitFile converts a parsed VehicleUnitFile back into its raw TV representation.
// This is the inverse of ParseRawVehicleUnitFile.
func UnparseVehicleUnitFile(file *vuv1.VehicleUnitFile) (*vuv1.RawVehicleUnitFile, error) {
	if file == nil {
		return nil, fmt.Errorf("vehicle unit file cannot be nil")
	}

	var records []*vuv1.RawVehicleUnitFile_Record
	marshalOpts := MarshalOptions{}

	// Helper to create a raw record from transfer value
	appendRecord := func(transferType vuv1.TransferType, transferValue []byte) error {
		if transferValue == nil {
			return nil
		}

		// Calculate signature size for this transfer type
		_, sigSize, err := sizeOfTransferValue(transferValue, transferType)
		if err != nil {
			return fmt.Errorf("failed to determine signature size: %w", err)
		}

		// Create record with complete transfer value
		record := &vuv1.RawVehicleUnitFile_Record{}
		record.SetType(transferType)
		record.SetGeneration(file.GetGeneration())
		record.SetValue(transferValue)          // Store complete value directly
		record.SetSignatureSize(int32(sigSize)) // Store signature size for efficient splitting
		records = append(records, record)

		return nil
	}

	switch file.GetGeneration() {
	case ddv1.Generation_GENERATION_1:
		gen1 := file.GetGen1()
		if gen1 == nil {
			return nil, fmt.Errorf("Gen1 data is nil")
		}

		// Unparse Overview (TREP 01)
		if overview := gen1.GetOverview(); overview != nil {
			transferValue, err := marshalOpts.MarshalOverviewGen1(overview)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Overview Gen1: %w", err)
			}
			if err := appendRecord(vuv1.TransferType_OVERVIEW_GEN1, transferValue); err != nil {
				return nil, err
			}
		}

		// Unparse Activities (TREP 02) - multiple transfers
		for i, activities := range gen1.GetActivities() {
			transferValue, err := marshalOpts.MarshalActivitiesGen1(activities)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal Activities Gen1 [%d]: %w", i, err)
			}
			if err := appendRecord(vuv1.TransferType_ACTIVITIES_GEN1, transferValue); err != nil {
				return nil, err
			}
		}

		// Unparse Events and Faults (TREP 03) - multiple transfers
		for i, eventsAndFaults := range gen1.GetEventsAndFaults() {
			transferValue, err := marshalOpts.MarshalEventsAndFaultsGen1(eventsAndFaults)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen1 [%d]: %w", i, err)
			}
			if err := appendRecord(vuv1.TransferType_EVENTS_AND_FAULTS_GEN1, transferValue); err != nil {
				return nil, err
			}
		}

		// Unparse Detailed Speed (TREP 04) - multiple transfers
		for i, detailedSpeed := range gen1.GetDetailedSpeed() {
			transferValue, err := marshalOpts.MarshalDetailedSpeedGen1(detailedSpeed)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen1 [%d]: %w", i, err)
			}
			if err := appendRecord(vuv1.TransferType_DETAILED_SPEED_GEN1, transferValue); err != nil {
				return nil, err
			}
		}

		// Unparse Technical Data (TREP 05) - multiple transfers
		for i, technicalData := range gen1.GetTechnicalData() {
			transferValue, err := marshalOpts.MarshalTechnicalDataGen1(technicalData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal TechnicalData Gen1 [%d]: %w", i, err)
			}
			if err := appendRecord(vuv1.TransferType_TECHNICAL_DATA_GEN1, transferValue); err != nil {
				return nil, err
			}
		}

	case ddv1.Generation_GENERATION_2:
		if file.GetVersion() == ddv1.Version_VERSION_2 {
			// Handle Gen2 V2
			gen2v2 := file.GetGen2V2()
			if gen2v2 == nil {
				return nil, fmt.Errorf("Gen2V2 data is nil")
			}

			// Unparse Overview (TREP 31)
			if overview := gen2v2.GetOverview(); overview != nil {
				transferValue, err := marshalOpts.MarshalOverviewGen2V2(overview)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Overview Gen2V2: %w", err)
				}
				if err := appendRecord(vuv1.TransferType_OVERVIEW_GEN2_V2, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Activities (TREP 32) - multiple transfers
			for i, activities := range gen2v2.GetActivities() {
				transferValue, err := marshalOpts.MarshalActivitiesGen2V2(activities)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Activities Gen2V2 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_ACTIVITIES_GEN2_V2, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Events and Faults (TREP 33) - multiple transfers
			for i, eventsAndFaults := range gen2v2.GetEventsAndFaults() {
				transferValue, err := marshalOpts.MarshalEventsAndFaultsGen2V2(eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen2V2 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V2, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Detailed Speed (TREP 34) - multiple transfers
			for i, detailedSpeed := range gen2v2.GetDetailedSpeed() {
				transferValue, err := marshalOpts.MarshalDetailedSpeedGen2(detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen2V2 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_DETAILED_SPEED_GEN2, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Technical Data (TREP 35) - multiple transfers
			for i, technicalData := range gen2v2.GetTechnicalData() {
				transferValue, err := marshalOpts.MarshalTechnicalDataGen2V2(technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal TechnicalData Gen2V2 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_TECHNICAL_DATA_GEN2_V2, transferValue); err != nil {
					return nil, err
				}
			}

		} else {
			// Handle Gen2 V1
			gen2v1 := file.GetGen2V1()
			if gen2v1 == nil {
				return nil, fmt.Errorf("Gen2V1 data is nil")
			}

			// Unparse Overview (TREP 11)
			if overview := gen2v1.GetOverview(); overview != nil {
				transferValue, err := marshalOpts.MarshalOverviewGen2V1(overview)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Overview Gen2V1: %w", err)
				}
				if err := appendRecord(vuv1.TransferType_OVERVIEW_GEN2_V1, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Activities (TREP 12) - multiple transfers
			for i, activities := range gen2v1.GetActivities() {
				transferValue, err := marshalOpts.MarshalActivitiesGen2V1(activities)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal Activities Gen2V1 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_ACTIVITIES_GEN2_V1, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Events and Faults (TREP 13) - multiple transfers
			for i, eventsAndFaults := range gen2v1.GetEventsAndFaults() {
				transferValue, err := marshalOpts.MarshalEventsAndFaultsGen2V1(eventsAndFaults)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal EventsAndFaults Gen2V1 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Detailed Speed (TREP 14) - multiple transfers
			for i, detailedSpeed := range gen2v1.GetDetailedSpeed() {
				transferValue, err := marshalOpts.MarshalDetailedSpeedGen2(detailedSpeed)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal DetailedSpeed Gen2V1 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_DETAILED_SPEED_GEN2, transferValue); err != nil {
					return nil, err
				}
			}

			// Unparse Technical Data (TREP 15) - multiple transfers
			for i, technicalData := range gen2v1.GetTechnicalData() {
				transferValue, err := marshalOpts.MarshalTechnicalDataGen2V1(technicalData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal TechnicalData Gen2V1 [%d]: %w", i, err)
				}
				if err := appendRecord(vuv1.TransferType_TECHNICAL_DATA_GEN2_V1, transferValue); err != nil {
					return nil, err
				}
			}
		}

	default:
		return nil, fmt.Errorf("unsupported generation: %v", file.GetGeneration())
	}

	rawFile := &vuv1.RawVehicleUnitFile{}
	rawFile.SetRecords(records)
	return rawFile, nil
}
