package card

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// unmarshalCardDownload unmarshals card download data from a card EF.
//
// The data type `LastCardDownload` is specified in the Data Dictionary, Section 2.89.
//
// ASN.1 Definition:
//
//	LastCardDownload ::= TimeReal
func (opts UnmarshalOptions) unmarshalCardDownload(data []byte) (*cardv1.CardDownloadDriver, error) {
	const (
		lenCardDownloadDriver = 4 // TimeReal size
	)

	if len(data) < lenCardDownloadDriver {
		return nil, fmt.Errorf("insufficient data for card download")
	}

	var target cardv1.CardDownloadDriver

	// Read timestamp (4 bytes)
	timestamp, err := opts.UnmarshalTimeReal(data[:lenCardDownloadDriver])
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	target.SetTimestamp(timestamp)

	return &target, nil
}

// MarshalCardDownload marshals card download data to bytes.
//
// The data type `LastCardDownload` is specified in the Data Dictionary, Section 2.89.
//
// ASN.1 Definition:
//
//	LastCardDownload ::= TimeReal
func (opts MarshalOptions) MarshalCardDownload(lastDownload *cardv1.CardDownloadDriver) ([]byte, error) {
	if lastDownload == nil {
		return nil, nil
	}

	// Timestamp (4 bytes)
	
	timestampBytes, err := opts.MarshalTimeReal(lastDownload.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal timestamp: %w", err)
	}

	return timestampBytes, nil
}
