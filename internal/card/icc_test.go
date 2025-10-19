package card

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestIccRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestIccRoundTrip(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/icc.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// First unmarshal
	unmarshalOpts := UnmarshalOptions{}
	icc1, err := unmarshalOpts.unmarshalIcc(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	// Marshal
	marshalOpts := MarshalOptions{}
	marshaled, err := marshalOpts.MarshalIcc(icc1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify binary equality
	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	// Second unmarshal
	icc2, err := unmarshalOpts.unmarshalIcc(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify structural equality
	if diff := cmp.Diff(icc1, icc2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestIccAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestIccAnonymization -update -v  # regenerate
func TestIccAnonymization(t *testing.T) {
	// Read test data
	b64Data, err := os.ReadFile("testdata/icc.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	// Unmarshal
	unmarshalOpts := UnmarshalOptions{}
	icc, err := unmarshalOpts.unmarshalIcc(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Anonymize
	anonymizeOpts := AnonymizeOptions{}
	anonymized := anonymizeOpts.anonymizeIcc(icc)

	// Marshal anonymized data
	marshalOpts := MarshalOptions{}
	anonymizedData, err := marshalOpts.MarshalIcc(anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		// Write anonymized binary
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/icc.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write icc.b64: %v", err)
		}

		// Write golden JSON with stable formatting
		// First convert to JSON using protojson
		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		// Then reformat with json.Indent for stable, deterministic output
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/icc.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write icc.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		// Assert binary matches
		expectedB64, err := os.ReadFile("testdata/icc.b64")
		if err != nil {
			t.Fatalf("Failed to read expected icc.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		// Assert JSON matches
		expectedJSON, err := os.ReadFile("testdata/icc.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.Icc
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	// Always: structural assertions on anonymized data
	if anonymized == nil {
		t.Fatal("Anonymized ICC is nil")
	}

	// Verify clock stop mode is set
	if anonymized.GetClockStop() == ddv1.ClockStopMode_CLOCK_STOP_MODE_UNSPECIFIED {
		t.Error("Clock stop mode should be set")
	}

	// Verify extended serial number
	esn := anonymized.GetCardExtendedSerialNumber()
	if esn == nil {
		t.Fatal("Extended serial number is nil")
	}
	if esn.GetSerialNumber() == 0 {
		t.Error("Serial number should be non-zero")
	}
	if esn.GetType() == ddv1.EquipmentType_EQUIPMENT_TYPE_UNSPECIFIED {
		t.Error("Equipment type should be set")
	}

	// Verify approval number
	approval := anonymized.GetCardApprovalNumber()
	if approval == nil {
		t.Fatal("Card approval number is nil")
	}
	if approval.GetValue() == "" {
		t.Error("Card approval number value should not be empty")
	}
}

// AnonymizeIcc creates an anonymized copy of ICC data, replacing sensitive identifiers
// with static, deterministic test values.
