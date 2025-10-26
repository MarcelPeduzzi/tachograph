package vu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/way-platform/tachograph-go/internal/hexdump"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// readHexdump reads and parses a hexdump file into binary data.
func readHexdump(path string) ([]byte, error) {
	hexdumpData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read hexdump file: %w", err)
	}

	data, err := hexdump.Unmarshal(hexdumpData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal hexdump: %w", err)
	}

	return data, nil
}

// goldenJSONPath converts a hexdump file path to its corresponding golden JSON path.
// Example: "testdata/records/000-14___FMS-379_.../000-OVERVIEW_GEN1.data.hexdump"
//
//	-> "testdata/records/000-14___FMS-379_.../000-OVERVIEW_GEN1.golden.json"
func goldenJSONPath(hexdumpPath string) string {
	return strings.TrimSuffix(hexdumpPath, ".data.hexdump") + ".golden.json"
}

// loadOrCreateGolden loads golden JSON for comparison or creates it in update mode.
// In update mode (-update flag), it writes the golden JSON file.
// In normal mode, it loads and compares against the existing golden JSON.
func loadOrCreateGolden(t *testing.T, message proto.Message, goldenPath string) {
	t.Helper()

	// Convert message to JSON with stable formatting
	jsonBytes, err := protojson.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf to JSON: %v", err)
	}

	// Reformat with json.Indent for stable output
	var indented []byte
	if len(jsonBytes) > 0 {
		var buf bytes.Buffer
		if err := json.Indent(&buf, jsonBytes, "", "  "); err != nil {
			t.Fatalf("Failed to format JSON: %v", err)
		}
		indented = buf.Bytes()
	}

	if *update {
		// Update mode: write golden file
		if err := os.WriteFile(goldenPath, indented, 0o644); err != nil {
			t.Fatalf("Failed to write golden JSON: %v", err)
		}
		t.Logf("Updated golden file: %s", goldenPath)
	} else {
		// Normal mode: compare with existing golden file
		expectedJSON, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("Failed to read golden JSON (run with -update to create): %v", err)
		}

		if diff := cmp.Diff(string(expectedJSON), string(indented)); diff != "" {
			t.Errorf("Golden JSON mismatch (-want +got):\n%s", diff)
		}
	}
}

// findHexdumpFiles discovers all hexdump files matching the specified transfer type.
// It searches recursively in testdata/records/ and returns absolute paths to matching files.
// Returns an error if the search fails, but returns an empty slice if no files match (caller should validate).
//
// VU hexdump files use the naming pattern:
// "NNN-<TRANSFER_TYPE>.data.hexdump" (e.g., "000-OVERVIEW_GEN1.data.hexdump")
// where the transfer type name already encodes the generation.
//
// Signature files are separate: "NNN-<TRANSFER_TYPE>.signature.hexdump"
func findHexdumpFiles(transferType vuv1.TransferType) ([]string, error) {
	// Convert enum to its string name
	transferTypeName := transferType.String()

	// Build the pattern we're looking for (data files only)
	// Example: "*-OVERVIEW_GEN1.data.hexdump"
	pattern := fmt.Sprintf("-%s.data.hexdump", transferTypeName)

	var matches []string

	// Walk the testdata/records/ directory
	recordsDir := "testdata/records"
	err := filepath.WalkDir(recordsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if filename matches our pattern
		if strings.HasSuffix(path, pattern) {
			matches = append(matches, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", recordsDir, err)
	}

	return matches, nil
}
