package tachograph

import (
	"context"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/card"
	"github.com/way-platform/tachograph-go/internal/vu"
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

	// Convert top-level options to internal options
	cardOpts := card.AuthenticateOptions{
		CertificateResolver: o.CertificateResolver,
	}
	vuOpts := vu.AuthenticateOptions{
		CertificateResolver: o.CertificateResolver,
	}

	switch rawFile.GetType() {
	case tachographv1.RawFile_CARD:
		return cardOpts.AuthenticateRawCardFile(ctx, rawFile.GetCard())
	case tachographv1.RawFile_VEHICLE_UNIT:
		return vuOpts.AuthenticateRawVehicleUnitFile(ctx, rawFile.GetVehicleUnit())
	default:
		return fmt.Errorf("unsupported file type: %v", rawFile.GetType())
	}
}
