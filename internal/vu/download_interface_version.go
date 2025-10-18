package vu

import (
	"fmt"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// ===== sizeOf Functions =====

// sizeOfDownloadInterfaceVersion returns the size of the DownloadInterfaceVersion value.
// This is a fixed 2-byte structure with no signature.
//
// Binary Layout (2 bytes total):
//   - generation: 1 byte
//   - version: 1 byte
//
// See Appendix 7, Section 2.2.6.1.
func sizeOfDownloadInterfaceVersion(data []byte, transferType vuv1.TransferType) (totalSize, signatureSize int, err error) {
	const lenDownloadInterfaceVersion = 2
	if len(data) < lenDownloadInterfaceVersion {
		return 0, 0, fmt.Errorf("insufficient data for DownloadInterfaceVersion: need %d, have %d", lenDownloadInterfaceVersion, len(data))
	}
	// No signature for DownloadInterfaceVersion
	return lenDownloadInterfaceVersion, 0, nil
}

// ===== Unmarshal Functions =====

// unmarshalDownloadInterfaceVersion parses the download interface version from VU data.
// It accepts the complete value (without the tag) and populates the raw_data field.
//
// The data type `DownloadInterfaceVersion` is specified in Appendix 7, Section 2.2.6.1.
//
// Binary Layout (2 bytes total):
//   - generation: 1 byte
//   - version: 1 byte
//
// ASN.1 Definition:
//
//	DownloadInterfaceVersion ::= OCTET STRING (SIZE (2))
