package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalGNSSPlaceRecord unmarshals a GNSSPlaceRecord from binary data.
//
// The data type `GNSSPlaceRecord` is specified in the Data Dictionary, Section 2.80.
//
// ASN.1 Definition:
//
//	GNSSPlaceRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates
//	}
//
// Binary Layout (11 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//
// This format is used consistently across all contexts (VU downloads, card files, etc.).
func (opts UnmarshalOptions) UnmarshalGNSSPlaceRecord(data []byte) (*ddv1.GNSSPlaceRecord, error) {
	const (
		lenGNSSPlaceRecord = 11
		idxTimestamp       = 0
		idxAccuracy        = 4
		idxGeoCoords       = 5
	)

	if len(data) != lenGNSSPlaceRecord {
		return nil, fmt.Errorf("invalid data length for GNSSPlaceRecord: got %d, want %d", len(data), lenGNSSPlaceRecord)
	}

	record := &ddv1.GNSSPlaceRecord{}

	// Parse timestamp (TimeReal - 4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[idxTimestamp : idxTimestamp+4])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}
	record.SetTimestamp(timestamp)

	// Parse gnssAccuracy (1 byte)
	accuracy := int32(data[idxAccuracy])
	record.SetGnssAccuracy(accuracy)

	// Parse geoCoordinates (6 bytes)
	geoCoords, err := opts.UnmarshalGeoCoordinates(data[idxGeoCoords : idxGeoCoords+6])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo coordinates: %w", err)
	}
	record.SetGeoCoordinates(geoCoords)

	return record, nil
}

// MarshalGNSSPlaceRecord marshals a GNSSPlaceRecord to bytes.
//
// The data type `GNSSPlaceRecord` is specified in the Data Dictionary, Section 2.80.
//
// ASN.1 Definition:
//
//	GNSSPlaceRecord ::= SEQUENCE {
//	    timeStamp TimeReal,
//	    gnssAccuracy GNSSAccuracy,
//	    geoCoordinates GeoCoordinates
//	}
//
// Binary Layout (11 bytes):
//
//	Offset 0: timeStamp (4 bytes)
//	Offset 4: gnssAccuracy (1 byte)
//	Offset 5: geoCoordinates (6 bytes)
//
// This format is used consistently across all contexts (VU downloads, card files, etc.).
func (opts MarshalOptions) MarshalGNSSPlaceRecord(record *ddv1.GNSSPlaceRecord) ([]byte, error) {
	if record == nil {
		// Return 11 zero bytes if no GNSS data
		return make([]byte, 11), nil
	}

	var dst []byte

	// Marshal timestamp (TimeReal - 4 bytes)
	timeBytes, err := opts.MarshalTimeReal(record.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal timestamp: %w", err)
	}
	dst = append(dst, timeBytes...)

	// Marshal gnssAccuracy (1 byte)
	accuracy := record.GetGnssAccuracy()
	if accuracy < 0 || accuracy > 255 {
		return nil, fmt.Errorf("invalid GNSS accuracy: %d (must be 0-255)", accuracy)
	}
	dst = append(dst, byte(accuracy))

	// Marshal geoCoordinates (6 bytes)
	geoBytes, err := opts.MarshalGeoCoordinates(record.GetGeoCoordinates())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal geo coordinates: %w", err)
	}
	dst = append(dst, geoBytes...)

	return dst, nil
}

// AnonymizeGNSSPlaceRecord creates an anonymized copy of GNSSPlaceRecord,
// replacing GNSS coordinates with a fixed, safe location (Helsinki, Finland)
// while preserving the timestamp and accuracy.
//
// Note: Timestamp normalization happens at the EF level (PlacesG2), not here.
//
// Helsinki coordinates: 60.17°N, 24.94°E
func AnonymizeGNSSPlaceRecord(record *ddv1.GNSSPlaceRecord) *ddv1.GNSSPlaceRecord {
	if record == nil {
		return nil
	}

	result := &ddv1.GNSSPlaceRecord{}

	// Preserve timestamp (will be normalized at EF level)
	result.SetTimestamp(record.GetTimestamp())

	// Preserve accuracy (structural information)
	result.SetGnssAccuracy(record.GetGnssAccuracy())

	// Replace coordinates with Helsinki, Finland
	// Helsinki: approximately 60°10'N, 24°56'E
	// Encoded as ±DDMM.M * 10 (latitude) and ±DDDMM.M * 10 (longitude)
	helsinkiGeo := &ddv1.GeoCoordinates{}
	helsinkiGeo.SetLatitude(60100)  // 60°10.0'N
	helsinkiGeo.SetLongitude(24560) // 24°56.0'E
	result.SetGeoCoordinates(helsinkiGeo)

	return result
}
