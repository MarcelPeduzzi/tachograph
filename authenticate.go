package tachograph

import (
	"context"
	"fmt"

	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// AuthenticateOptions configures the signature authentication process.
type AuthenticateOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	// If nil, this defaults to using DefaultCertificateResolver.
	CertificateResolver CertificateResolver
}

// AuthenticateRawFile performs cryptographic authentication on a raw tachograph file,
// populating Authentication fields in the raw records.
// This function mutates the rawFile parameter.
func AuthenticateRawFile(ctx context.Context, rawFile *tachographv1.RawFile) error {
	return AuthenticateOptions{}.AuthenticateRawFile(ctx, rawFile)
}

// AuthenticateRawFile performs cryptographic authentication with custom options.
func (o AuthenticateOptions) AuthenticateRawFile(ctx context.Context, rawFile *tachographv1.RawFile) error {
	if rawFile == nil {
		return fmt.Errorf("rawFile cannot be nil")
	}
	if o.CertificateResolver == nil {
		o.CertificateResolver = DefaultCertificateResolver()
	}

	switch rawFile.GetType() {
	case tachographv1.RawFile_CARD:
		return o.authenticateCardFile(ctx, rawFile.GetCard())
	case tachographv1.RawFile_VEHICLE_UNIT:
		return o.authenticateVehicleUnitFile(ctx, rawFile.GetVehicleUnit())
	default:
		return fmt.Errorf("unsupported file type: %v", rawFile.GetType())
	}
}

// authenticateCardFile authenticates a raw card file.
// TODO: Implement card file authentication
func (o AuthenticateOptions) authenticateCardFile(ctx context.Context, rawCardFile interface{}) error {
	// TODO: Implement
	// 1. Iterate through records
	// 2. Find signature records
	// 3. Verify against data records
	// 4. Populate authentication field in corresponding data records
	return nil
}

// authenticateVehicleUnitFile authenticates a raw VU file.
// TODO: Implement VU file authentication
func (o AuthenticateOptions) authenticateVehicleUnitFile(ctx context.Context, rawVUFile interface{}) error {
	// TODO: Implement
	// 1. For each record
	// 2. Extract embedded signature from value bytes
	// 3. Verify signature
	// 4. Populate authentication field
	return nil
}
