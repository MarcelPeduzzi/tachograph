package vu

import (
	"context"
	"errors"
	"fmt"

	"github.com/way-platform/tachograph-go/internal/cert"
	"github.com/way-platform/tachograph-go/internal/security"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// AuthenticateOptions configures the VU authentication process.
type AuthenticateOptions struct {
	// CertificateResolver is used to resolve CA certificates by their Certificate Authority Reference (CAR).
	CertificateResolver cert.Resolver
}

// AuthenticateRawVehicleUnitFile performs cryptographic authentication on all records
// in a raw vehicle unit file. For each record, it verifies the signature against the
// data and populates the authentication field.
//
// This function mutates the rawFile parameter by setting the authentication field
// on each record.
//
// Authentication is performed by:
// 1. Extracting the equipment certificate from the appropriate transfer type
// 2. Verifying the certificate chain against the European Root CA
// 3. Verifying the data signature using the equipment certificate
//
// If any record fails authentication, its authentication.status will be set to the
// appropriate error status, and this function will return an error after processing
// all records.
func (opts AuthenticateOptions) AuthenticateRawVehicleUnitFile(ctx context.Context, rawFile *vuv1.RawVehicleUnitFile) error {
	if rawFile == nil {
		return fmt.Errorf("rawFile cannot be nil")
	}
	if opts.CertificateResolver == nil {
		opts.CertificateResolver = cert.DefaultResolver()
	}

	var errs []error
	records := rawFile.GetRecords()

	for _, record := range records {
		if err := opts.authenticateRecord(ctx, record, records); err != nil {
			errs = append(errs, err)
			// Continue processing other records even if one fails
		}
	}

	return errors.Join(errs...)
}

// authenticateRecord authenticates a single VU record by verifying its signature
// against the data. This function mutates the record by setting its authentication field.
func (opts AuthenticateOptions) authenticateRecord(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record) error {
	generation := record.GetGeneration()
	transferType := record.GetType()

	// Initialize authentication result
	auth := &securityv1.Authentication{}
	record.SetAuthentication(auth)

	// Select authentication strategy based on generation
	switch generation {
	case ddv1.Generation_GENERATION_1:
		return opts.authenticateGen1Record(ctx, record, allRecords, auth)
	case ddv1.Generation_GENERATION_2:
		return opts.authenticateGen2Record(ctx, record, allRecords, auth)
	default:
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("unsupported generation for %v: %v", transferType, generation)
	}
}

// authenticateGen1Record authenticates a Generation 1 VU record using RSA signature recovery.
func (opts AuthenticateOptions) authenticateGen1Record(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record, auth *securityv1.Authentication) error {
	// Step 1: Find the Overview record to get the VU and MSCA certificates
	overviewRecord := opts.findOverviewRecord(allRecords)
	if overviewRecord == nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("Overview record not found for authentication")
	}

	// Step 2: Parse the Overview to extract certificates
	vuCert, mscaCert, err := opts.extractGen1Certificates(overviewRecord)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to extract Gen1 certificates: %w", err)
	}

	// Step 3: Verify certificate chain
	if err := opts.verifyGen1CertificateChain(ctx, vuCert, mscaCert, auth); err != nil {
		return err
	}

	// Step 4: Verify the data signature
	if err := opts.verifyGen1DataSignature(record, vuCert, auth); err != nil {
		return err
	}

	// Authentication succeeded
	auth.SetStatus(securityv1.Authentication_VERIFIED)
	auth.SetSignatureAlgorithm(securityv1.SignatureAlgorithm_SHA1_WITH_RSA_ENCRYPTION)

	return nil
}

// authenticateGen2Record authenticates a Generation 2 VU record using ECDSA signature verification.
func (opts AuthenticateOptions) authenticateGen2Record(ctx context.Context, record *vuv1.RawVehicleUnitFile_Record, allRecords []*vuv1.RawVehicleUnitFile_Record, auth *securityv1.Authentication) error {
	// TODO: Implement Gen2 authentication
	// 1. Find VU certificate from appropriate transfer type (e.g., VuCertificateGen2)
	// 2. Verify certificate chain: EUR Root -> MSCA -> VU Certificate
	// 3. Use ECDSA verification to verify data signature
	// 4. Populate auth with result

	auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
	return fmt.Errorf("Gen2 authentication not yet implemented")
}

// findOverviewRecord finds the first Gen1 Overview record in the file.
// This record contains the VU and MSCA certificates needed for authentication.
func (opts AuthenticateOptions) findOverviewRecord(records []*vuv1.RawVehicleUnitFile_Record) *vuv1.RawVehicleUnitFile_Record {
	for _, rec := range records {
		if rec.GetType() == vuv1.TransferType_OVERVIEW_GEN1 {
			return rec
		}
	}
	return nil
}

// extractGen1Certificates extracts the VU and MSCA certificates from a Gen1 Overview record.
func (opts AuthenticateOptions) extractGen1Certificates(overviewRecord *vuv1.RawVehicleUnitFile_Record) (*securityv1.RsaCertificate, *securityv1.RsaCertificate, error) {
	data := overviewRecord.GetData()

	// Gen1 Overview structure (from regulation):
	// - MemberStateCertificate: 194 bytes
	// - VuCertificate: 194 bytes
	// - VehicleIdentificationNumber: variable
	// - ...rest of data

	const lenRsaCertificate = 194
	const idxMscaCert = 0
	const idxVuCert = 194
	const minDataLen = 388 // Two certificates

	if len(data) < minDataLen {
		return nil, nil, fmt.Errorf("insufficient data for certificates: got %d, need at least %d", len(data), minDataLen)
	}

	// Extract MSCA certificate
	mscaCertData := data[idxMscaCert : idxMscaCert+lenRsaCertificate]
	mscaCert := &securityv1.RsaCertificate{}
	mscaCert.SetRawData(mscaCertData)

	// Extract VU certificate
	vuCertData := data[idxVuCert : idxVuCert+lenRsaCertificate]
	vuCert := &securityv1.RsaCertificate{}
	vuCert.SetRawData(vuCertData)

	return vuCert, mscaCert, nil
}

// verifyGen1CertificateChain verifies the certificate chain for Gen1:
// EUR Root -> MSCA -> VU
func (opts AuthenticateOptions) verifyGen1CertificateChain(ctx context.Context, vuCert *securityv1.RsaCertificate, mscaCert *securityv1.RsaCertificate, auth *securityv1.Authentication) error {
	// Step 1: Get EUR root certificate
	rootCert, err := opts.CertificateResolver.GetRootCertificate(ctx)
	if err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("failed to get root certificate: %w", err)
	}

	// Step 2: Verify MSCA certificate against EUR root
	if err := security.VerifyRsaCertificateWithRoot(mscaCert, rootCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("MSCA certificate verification failed: %w", err)
	}

	// Step 3: Verify VU certificate against MSCA
	if err := security.VerifyRsaCertificateWithCA(vuCert, mscaCert); err != nil {
		auth.SetStatus(securityv1.Authentication_CERTIFICATE_VERIFICATION_FAILED)
		return fmt.Errorf("VU certificate verification failed: %w", err)
	}

	// Step 4: Populate certificate info in auth
	// TODO: Extract certificate info (CHR, nation, validity dates) and populate
	// auth.signer_certificate and auth.root_certificate

	return nil
}

// verifyGen1DataSignature verifies the RSA signature on the data portion of a Gen1 record.
func (opts AuthenticateOptions) verifyGen1DataSignature(record *vuv1.RawVehicleUnitFile_Record, vuCert *securityv1.RsaCertificate, auth *securityv1.Authentication) error {
	allData := record.GetData()
	signature := record.GetSignature()

	// Gen1 uses RSA-1024 with 128-byte signatures
	const lenRsaSignature = 128
	if len(signature) != lenRsaSignature {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("invalid signature length for Gen1: got %d, want %d", len(signature), lenRsaSignature)
	}

	// For Overview (TREP 01), Technical Data (TREP 05), and other Gen1 records,
	// the signature is over the data AFTER the certificates.
	// Each RSA certificate is 194 bytes: MSCA cert + VU cert = 388 bytes total.
	var signedData []byte
	transferType := record.GetType()

	switch transferType {
	case vuv1.TransferType_OVERVIEW_GEN1:
		// Overview: Signature covers data from VehicleIdentificationNumber onwards (after 2 certificates)
		const lenCertificates = 388 // 194 bytes MSCA + 194 bytes VU
		if len(allData) <= lenCertificates {
			auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
			return fmt.Errorf("insufficient data for Overview signature verification: got %d, need > %d", len(allData), lenCertificates)
		}
		signedData = allData[lenCertificates:]

	case vuv1.TransferType_TECHNICAL_DATA_GEN1:
		// Technical Data: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_ACTIVITIES_GEN1:
		// Activities: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1:
		// Events/Faults: Signature covers all data (no certificates in this record)
		signedData = allData

	case vuv1.TransferType_DETAILED_SPEED_GEN1:
		// Detailed Speed: Signature covers all data (no certificates in this record)
		signedData = allData

	default:
		// For unknown transfer types, assume signature covers all data
		signedData = allData
	}

	// Verify the signature using PKCS#1 v1.5
	if err := security.VerifyRsaDataSignature(signedData, signature, vuCert); err != nil {
		auth.SetStatus(securityv1.Authentication_DATA_SIGNATURE_INVALID)
		return fmt.Errorf("data signature verification failed: %w", err)
	}

	return nil
}
