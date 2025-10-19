package card

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestSpecificConditionsRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestSpecificConditionsRoundTrip(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/specific_conditions.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	unmarshalOpts := UnmarshalOptions{}
	sc1, err := unmarshalOpts.unmarshalSpecificConditions(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	marshalOpts := MarshalOptions{}
	marshaled, err := marshalOpts.MarshalCardSpecificConditions(sc1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	sc2, err := unmarshalOpts.unmarshalSpecificConditions(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	if diff := cmp.Diff(sc1, sc2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestSpecificConditionsAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestSpecificConditionsAnonymization -update -v  # regenerate
func TestSpecificConditionsAnonymization(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/specific_conditions.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	unmarshalOpts := UnmarshalOptions{}
	sc, err := unmarshalOpts.unmarshalSpecificConditions(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	anonymizeOpts := AnonymizeOptions{}
	anonymized := anonymizeOpts.anonymizeSpecificConditions(sc)

	marshalOpts := MarshalOptions{}
	anonymizedData, err := marshalOpts.MarshalCardSpecificConditions(anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/specific_conditions.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write specific_conditions.b64: %v", err)
		}

		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/specific_conditions.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write specific_conditions.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		expectedB64, err := os.ReadFile("testdata/specific_conditions.b64")
		if err != nil {
			t.Fatalf("Failed to read expected specific_conditions.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		expectedJSON, err := os.ReadFile("testdata/specific_conditions.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.SpecificConditions
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	if anonymized == nil {
		t.Fatal("Anonymized SpecificConditions is nil")
	}

	// Verify record count is preserved
	if len(anonymized.GetRecords()) != len(sc.GetRecords()) {
		t.Errorf("Record count changed: got %d, want %d",
			len(anonymized.GetRecords()), len(sc.GetRecords()))
	}
}

// AnonymizeSpecificConditions creates an anonymized copy of SpecificConditions.
// Specific conditions are categorical (not personally identifiable), but we anonymize timestamps.
