package card

import (
	"context"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/cert"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// AuthenticateOptions configures the card authentication process.
type AuthenticateOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	CertificateResolver cert.Resolver
}

// AuthenticateRawCardFile performs cryptographic authentication on all signed Elementary Files
// in a raw card file. For each signed EF, it verifies the signature in the following signature
// block and populates the authentication field.
//
// This function mutates the rawFile parameter by setting the authentication field
// on signed EF records.
//
// Authentication is performed by:
// 1. Extracting the card certificate from the appropriate EF (e.g., EF_Card_Certificate)
// 2. Verifying the certificate chain against the European Root CA
// 3. Verifying the EF data signature using the card certificate
//
// If any EF fails authentication, its authentication.status will be set to the
// appropriate error status, and this function will return an error after processing
// all EFs.
func (opts AuthenticateOptions) AuthenticateRawCardFile(ctx context.Context, rawFile *cardv1.RawCardFile) error {
	if rawFile == nil {
		return fmt.Errorf("rawFile cannot be nil")
	}
	if opts.CertificateResolver == nil {
		opts.CertificateResolver = cert.DefaultResolver()
	}

	// Infer card generation and type from the file structure
	generation := opts.inferGeneration(rawFile)
	cardType := opts.inferCardType(rawFile)

	// Select authentication strategy based on generation
	switch generation {
	case ddv1.Generation_GENERATION_1:
		return opts.authenticateGen1CardFile(ctx, rawFile, cardType)
	case ddv1.Generation_GENERATION_2:
		return opts.authenticateGen2CardFile(ctx, rawFile, cardType)
	default:
		return fmt.Errorf("unable to determine card generation")
	}
}

// authenticateGen1CardFile authenticates a Generation 1 card file using RSA signature recovery.
func (opts AuthenticateOptions) authenticateGen1CardFile(ctx context.Context, rawFile *cardv1.RawCardFile, cardType cardv1.CardType) error {
	// TODO: Implement Gen1 card authentication
	// 1. Find card certificate (EF_Card_Certificate)
	// 2. Verify certificate chain: EUR Root -> MSCA -> Card Certificate
	// 3. For each signed EF:
	//    a. Find corresponding signature record
	//    b. Use RSA signature recovery (ISO/IEC 9796-2) to verify EF data
	//    c. Populate authentication field on the EF record

	return fmt.Errorf("Gen1 card authentication not yet implemented")
}

// authenticateGen2CardFile authenticates a Generation 2 card file using ECDSA signature verification.
func (opts AuthenticateOptions) authenticateGen2CardFile(ctx context.Context, rawFile *cardv1.RawCardFile, cardType cardv1.CardType) error {
	// TODO: Implement Gen2 card authentication
	// 1. Find card certificate (EF_Card_Certificate)
	// 2. Verify certificate chain: EUR Root -> MSCA -> Card Certificate
	// 3. For each signed EF:
	//    a. Find corresponding signature record
	//    b. Use ECDSA verification to verify EF data
	//    c. Populate authentication field on the EF record

	return fmt.Errorf("Gen2 card authentication not yet implemented")
}

// inferGeneration determines the card generation from the raw file structure.
func (opts AuthenticateOptions) inferGeneration(rawFile *cardv1.RawCardFile) ddv1.Generation {
	// Look for generation indicators in the file structure
	for _, record := range rawFile.GetRecords() {
		if record.GetTag()&0x80 != 0 {
			// Bit 7 set indicates Generation 2
			return ddv1.Generation_GENERATION_2
		}
	}
	return ddv1.Generation_GENERATION_1
}

// inferCardType determines the card type from the raw file structure.
func (opts AuthenticateOptions) inferCardType(rawFile *cardv1.RawCardFile) cardv1.CardType {
	// TODO: Implement proper card type inference
	// For now, default to driver card
	return cardv1.CardType_DRIVER_CARD
}
