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

// TestControlActivityDataRoundTrip verifies binary fidelity (unmarshal → marshal → unmarshal)
func TestControlActivityDataRoundTrip(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/control_activity.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	unmarshalOpts := UnmarshalOptions{}
	ca1, err := unmarshalOpts.unmarshalControlActivityData(data)
	if err != nil {
		t.Fatalf("First unmarshal failed: %v", err)
	}

	opts := MarshalOptions{}
	marshaled, err := opts.MarshalCardControlActivityData(ca1)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if diff := cmp.Diff(data, marshaled); diff != "" {
		t.Errorf("Binary mismatch after marshal (-want +got):\n%s", diff)
	}

	ca2, err := unmarshalOpts.unmarshalControlActivityData(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	if diff := cmp.Diff(ca1, ca2, protocmp.Transform()); diff != "" {
		t.Errorf("Structural mismatch after round-trip (-want +got):\n%s", diff)
	}
}

// TestControlActivityDataAnonymization is a golden file test with deterministic anonymization
//
//	go test -run TestControlActivityDataAnonymization -update -v  # regenerate
func TestControlActivityDataAnonymization(t *testing.T) {
	b64Data, err := os.ReadFile("testdata/control_activity.b64")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	unmarshalOpts := UnmarshalOptions{}
	ca, err := unmarshalOpts.unmarshalControlActivityData(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	anonymizeOpts := AnonymizeOptions{}
	anonymized := anonymizeOpts.anonymizeControlActivityData(ca)

	marshalOpts := MarshalOptions{}
	anonymizedData, err := marshalOpts.MarshalCardControlActivityData(anonymized)
	if err != nil {
		t.Fatalf("Failed to marshal anonymized data: %v", err)
	}

	if *update {
		anonymizedB64 := base64.StdEncoding.EncodeToString(anonymizedData)
		if err := os.WriteFile("testdata/control_activity.b64", []byte(anonymizedB64), 0o644); err != nil {
			t.Fatalf("Failed to write control_activity.b64: %v", err)
		}

		jsonBytes, err := protojson.Marshal(anonymized)
		if err != nil {
			t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
		}
		var stableJSON bytes.Buffer
		if err := json.Indent(&stableJSON, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		if err := os.WriteFile("testdata/control_activity.golden.json", stableJSON.Bytes(), 0o644); err != nil {
			t.Fatalf("Failed to write control_activity.golden.json: %v", err)
		}

		t.Log("Updated golden files")
	} else {
		expectedB64, err := os.ReadFile("testdata/control_activity.b64")
		if err != nil {
			t.Fatalf("Failed to read expected control_activity.b64: %v", err)
		}
		expectedData, err := base64.StdEncoding.DecodeString(string(expectedB64))
		if err != nil {
			t.Fatalf("Failed to decode expected base64: %v", err)
		}
		if diff := cmp.Diff(expectedData, anonymizedData); diff != "" {
			t.Errorf("Binary mismatch (-want +got):\n%s", diff)
		}

		expectedJSON, err := os.ReadFile("testdata/control_activity.golden.json")
		if err != nil {
			t.Fatalf("Failed to read expected JSON: %v", err)
		}
		var expected cardv1.ControlActivityData
		if err := protojson.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if diff := cmp.Diff(&expected, anonymized, protocmp.Transform()); diff != "" {
			t.Errorf("JSON mismatch (-want +got):\n%s", diff)
		}
	}

	if anonymized == nil {
		t.Fatal("Anonymized ControlActivityData is nil")
	}

	// Verify valid status is preserved
	if anonymized.GetValid() != ca.GetValid() {
		t.Errorf("Valid status changed: got %v, want %v", anonymized.GetValid(), ca.GetValid())
	}

	if anonymized.GetValid() {
		// Verify control type is preserved (categorical)
		if anonymized.GetControlType() == nil {
			t.Error("Control type should not be nil")
		}

		// Verify card number is anonymized
		if cardNum := anonymized.GetControlCardNumber(); cardNum != nil {
			if fcn := cardNum.GetFullCardNumber(); fcn != nil {
				// Check driver identification (for driver cards)
				if driverID := fcn.GetDriverIdentification(); driverID != nil {
					if idStr := driverID.GetDriverIdentificationNumber(); idStr != nil && idStr.GetValue() == "" {
						t.Error("Driver identification should not be empty")
					}
				}
			}
		}

		// Verify vehicle registration is FINLAND
		if vehicleReg := anonymized.GetControlVehicleRegistration(); vehicleReg != nil {
			if vehicleReg.GetNation() != ddv1.NationNumeric_FINLAND {
				t.Errorf("Vehicle nation = %v, want FINLAND", vehicleReg.GetNation())
			}
		}
	}
}

// AnonymizeControlActivityData creates an anonymized copy of ControlActivityData.
