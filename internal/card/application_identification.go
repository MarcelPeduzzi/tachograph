package card

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/dd"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalApplicationIdentification parses the binary data for an EF_ApplicationIdentification record (Gen1 format).
//
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition (Gen1):
//
//	ApplicationIdentification ::= SEQUENCE {
//	    typeOfTachographCardId    EquipmentType,
//	    cardStructureVersion      CardStructureVersion,
//	    noOfEventsPerType         INTEGER(0..255),
//	    noOfFaultsPerType         INTEGER(0..255),
//	    activityStructureLength   INTEGER(0..65535),
//	    noOfCardVehicleRecords    INTEGER(0..255),
//	    noOfCardPlaceRecords      INTEGER(0..255)
//	}
func (opts UnmarshalOptions) unmarshalApplicationIdentification(data []byte) (*cardv1.ApplicationIdentification, error) {
	const (
		lenEfApplicationIdentificationGen1 = 10 // Gen1: 1 + 2 + 1 + 1 + 2 + 2 + 1 = 10 bytes for driver cards
	)

	if len(data) != lenEfApplicationIdentificationGen1 {
		return nil, fmt.Errorf("invalid data length for Gen1 application identification: got %d bytes, want %d", len(data), lenEfApplicationIdentificationGen1)
	}

	target := &cardv1.ApplicationIdentification{}
	r := bytes.NewReader(data)

	// Read type of tachograph card ID (1 byte)
	var cardType byte
	if err := binary.Read(r, binary.BigEndian, &cardType); err != nil {
		return nil, fmt.Errorf("failed to read card type: %w", err)
	}
	// Convert raw card type to enum using protocol annotations
	if equipmentType, err := dd.UnmarshalEnum[ddv1.EquipmentType](cardType); err == nil {
		target.SetTypeOfTachographCardId(equipmentType)
	} else {
		return nil, fmt.Errorf("invalid equipment type: %w", err)
	}

	// Read card structure version (2 bytes)
	structureVersionBytes := make([]byte, 2)
	if _, err := r.Read(structureVersionBytes); err != nil {
		return nil, fmt.Errorf("failed to read card structure version: %w", err)
	}
	// Parse BCD structure version using centralized helper
	cardStructureVersion, err := opts.UnmarshalCardStructureVersion(structureVersionBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal card structure version: %w", err)
	}
	target.SetCardStructureVersion(cardStructureVersion)

	// For now, assume this is a driver card and create the driver data
	driver := &cardv1.ApplicationIdentification_Driver{}

	// Read events per type count (1 byte)
	var eventsPerType byte
	if err := binary.Read(r, binary.BigEndian, &eventsPerType); err != nil {
		return nil, fmt.Errorf("failed to read events per type count: %w", err)
	}
	driver.SetEventsPerTypeCount(int32(eventsPerType))

	// Read faults per type count (1 byte)
	var faultsPerType byte
	if err := binary.Read(r, binary.BigEndian, &faultsPerType); err != nil {
		return nil, fmt.Errorf("failed to read faults per type count: %w", err)
	}
	driver.SetFaultsPerTypeCount(int32(faultsPerType))

	// Read activity structure length (2 bytes)
	var activityLength uint16
	if err := binary.Read(r, binary.BigEndian, &activityLength); err != nil {
		return nil, fmt.Errorf("failed to read activity structure length: %w", err)
	}
	driver.SetActivityStructureLength(int32(activityLength))

	// Read card vehicle records count (2 bytes in Gen1)
	var vehicleRecords uint16
	if err := binary.Read(r, binary.BigEndian, &vehicleRecords); err != nil {
		return nil, fmt.Errorf("failed to read vehicle records count: %w", err)
	}
	driver.SetCardVehicleRecordsCount(int32(vehicleRecords))

	// Read card place records count (1 byte in Gen1)
	var placeRecords byte
	if err := binary.Read(r, binary.BigEndian, &placeRecords); err != nil {
		return nil, fmt.Errorf("failed to read place records count: %w", err)
	}
	driver.SetCardPlaceRecordsCount(int32(placeRecords))

	// Set the driver data and card type
	target.SetDriver(driver)
	target.SetCardType(cardv1.CardType_DRIVER_CARD)

	return target, nil
}

// MarshalCardApplicationIdentification marshals Gen1 application identification data.
//
// The data type `ApplicationIdentification` is specified in the Data Dictionary, Section 2.2.
//
// ASN.1 Definition (Gen1):
//
//	ApplicationIdentification ::= SEQUENCE {
//	    typeOfTachographCardId    EquipmentType,
//	    cardStructureVersion      CardStructureVersion,
//	    noOfEventsPerType         INTEGER(0..255),
//	    noOfFaultsPerType         INTEGER(0..255),
//	    activityStructureLength   INTEGER(0..65535),
//	    noOfCardVehicleRecords    INTEGER(0..255),
//	    noOfCardPlaceRecords      INTEGER(0..255),
//	    noOfCalibrationRecords    INTEGER(0..255)
//	}
func (opts MarshalOptions) MarshalCardApplicationIdentification(appId *cardv1.ApplicationIdentification) ([]byte, error) {
	if appId == nil {
		return nil, nil
	}

	var data []byte

	// Type of tachograph card ID (1 byte)
	if appId.HasTypeOfTachographCardId() {
		protocolValue, _ := dd.MarshalEnum(appId.GetTypeOfTachographCardId())
		data = append(data, protocolValue)
	} else {
		data = append(data, 0x00)
	}

	// Card structure version (2 bytes)
	structureVersion := appId.GetCardStructureVersion()
	if structureVersion != nil {
		// Marshal using centralized helper
		
		versionBytes, err := opts.MarshalCardStructureVersion(structureVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal card structure version: %w", err)
		}
		data = append(data, versionBytes...)
	} else {
		data = append(data, 0x00, 0x01) // Default version
	}

	// Get driver data for the specific fields
	var driver *cardv1.ApplicationIdentification_Driver
	switch appId.GetCardType() {
	case cardv1.CardType_DRIVER_CARD:
		driver = appId.GetDriver()
	}

	if driver == nil {
		// If no driver data, append zeros for all driver-specific fields (7 bytes)
		data = append(data, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // events, faults, activity length (2 bytes), vehicle records (2 bytes), place records
		return data, nil
	}

	// Events per type count (1 byte)
	if driver.HasEventsPerTypeCount() {
		data = append(data, byte(driver.GetEventsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Faults per type count (1 byte)
	if driver.HasFaultsPerTypeCount() {
		data = append(data, byte(driver.GetFaultsPerTypeCount()))
	} else {
		data = append(data, 0x00)
	}

	// Activity structure length (2 bytes)
	if driver.HasActivityStructureLength() {
		activityLength := make([]byte, 2)
		binary.BigEndian.PutUint16(activityLength, uint16(driver.GetActivityStructureLength()))
		data = append(data, activityLength...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card vehicle records count (2 bytes in Gen1)
	if driver.HasCardVehicleRecordsCount() {
		vehicleRecords := make([]byte, 2)
		binary.BigEndian.PutUint16(vehicleRecords, uint16(driver.GetCardVehicleRecordsCount()))
		data = append(data, vehicleRecords...)
	} else {
		data = append(data, 0x00, 0x00)
	}

	// Card place records count (1 byte in Gen1)
	if driver.HasCardPlaceRecordsCount() {
		data = append(data, byte(driver.GetCardPlaceRecordsCount()))
	} else {
		data = append(data, 0x00)
	}

	return data, nil
}
