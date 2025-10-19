package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalPreviousVehicleInfo unmarshals a Generation 1 PreviousVehicleInfo (19 bytes).
//
// The data type `PreviousVehicleInfo` is specified in the Data Dictionary, Section 2.118.
//
// ASN.1 Definition (Gen1):
//
//	PreviousVehicleInfo ::= SEQUENCE {
//	    vehicleRegistrationIdentification VehicleRegistrationIdentification,
//	    cardWithdrawalTime                TimeReal
//	}
func (opts UnmarshalOptions) UnmarshalPreviousVehicleInfo(data []byte) (*ddv1.PreviousVehicleInfo, error) {
	const (
		idxVehicleReg          = 0
		idxCardWithdrawalTime  = 15
		lenPreviousVehicleInfo = 19 // Fixed size for Gen1
	)

	if len(data) != lenPreviousVehicleInfo {
		return nil, fmt.Errorf("invalid data length for Gen1 PreviousVehicleInfo: got %d, want %d", len(data), lenPreviousVehicleInfo)
	}

	result := &ddv1.PreviousVehicleInfo{}
	if opts.PreserveRawData {
		result.SetRawData(data)
	}
	// Parse vehicleRegistrationIdentification (15 bytes)
	vehicleReg, err := opts.UnmarshalVehicleRegistration(data[idxVehicleReg : idxVehicleReg+15])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle registration: %w", err)
	}
	result.SetVehicleRegistration(vehicleReg)

	// Parse cardWithdrawalTime (TimeReal - 4 bytes)
	withdrawalTime, err := opts.UnmarshalTimeReal(data[idxCardWithdrawalTime : idxCardWithdrawalTime+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal card withdrawal time: %w", err)
	}
	result.SetCardWithdrawalTime(withdrawalTime)

	return result, nil
}

// MarshalPreviousVehicleInfo marshals a Generation 1 PreviousVehicleInfo (19 bytes) to bytes.
func (opts MarshalOptions) MarshalPreviousVehicleInfo(info *ddv1.PreviousVehicleInfo) ([]byte, error) {
	const lenPreviousVehicleInfo = 19 // Fixed size for Gen1

	// Use raw data painting strategy if available
	var canvas [lenPreviousVehicleInfo]byte
	if rawData := info.GetRawData(); len(rawData) > 0 {
		if len(rawData) != lenPreviousVehicleInfo {
			return nil, fmt.Errorf("invalid raw_data length for PreviousVehicleInfo: got %d, want %d", len(rawData), lenPreviousVehicleInfo)
		}
		copy(canvas[:], rawData)
	}

	// Paint semantic values over the canvas
	// Vehicle registration (15 bytes)
	vehicleRegBytes, err := opts.MarshalVehicleRegistration(info.GetVehicleRegistration())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vehicle registration: %w", err)
	}
	copy(canvas[0:15], vehicleRegBytes)

	// Card withdrawal time (4 bytes)
	timeBytes, err := opts.MarshalTimeReal(info.GetCardWithdrawalTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card withdrawal time: %w", err)
	}
	copy(canvas[15:19], timeBytes)

	return canvas[:], nil
}
