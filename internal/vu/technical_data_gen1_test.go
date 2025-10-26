package vu

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

func TestTechnicalData_Gen1(t *testing.T) {
	// Discover all matching hexdump files
	hexdumpFiles, err := findHexdumpFiles(vuv1.TransferType_TECHNICAL_DATA_GEN1)
	if err != nil {
		t.Fatalf("Failed to discover hexdump files: %v", err)
	}
	if len(hexdumpFiles) == 0 {
		t.Skip("No hexdump files found for TECHNICAL_DATA_GEN1")
	}

	// Run subtest for each discovered file
	for _, hexdumpPath := range hexdumpFiles {
		// Use relative path from testdata as subtest name
		relPath := strings.TrimPrefix(hexdumpPath, "testdata/records/")
		testName := strings.TrimSuffix(relPath, ".data.hexdump")

		t.Run(testName, func(t *testing.T) {
			// Read hexdump
			data, err := readHexdump(hexdumpPath)
			if err != nil {
				t.Fatalf("Failed to read hexdump: %v", err)
			}

			// Unmarshal
			technicalData, err := unmarshalTechnicalDataGen1(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if technicalData == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Golden JSON comparison
			goldenPath := goldenJSONPath(hexdumpPath)
			loadOrCreateGolden(t, technicalData, goldenPath)

			// Round-trip test - marshal
			marshalOpts := MarshalOptions{}
			marshaled, err := marshalOpts.MarshalTechnicalDataGen1(technicalData)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if diff := cmp.Diff(data, marshaled); diff != "" {
				t.Errorf("Binary round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
