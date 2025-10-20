# Log: Anonymization and Marshalling Principles

This document establishes clear principles for the anonymization and marshalling of tachograph files to ensure that anonymized data is secure, compliant, and useful for testing.

## The Problem

The current marshalling process defaults to a "raw data painting" strategy (`UseRawData: true`), where it uses the original binary data (`raw_data` fields in the protobuf) as a canvas and paints semantic changes over it.

When anonymizing a file, the process is:
1. Parse a binary file into a protobuf message, preserving the original bytes in `raw_data` fields.
2. Clone the message and change semantic fields (e.g., names, dates) to anonymous values. The `raw_data` fields remain untouched copies of the original data.
3. Marshal the anonymized message. The marshaller uses the unmodified `raw_data` as a base, effectively leaking original data into the "anonymized" output.

Additionally, digital signatures are not handled correctly, remaining in the file even though they are invalidated by the anonymization process.

## Principles

To solve this, we will adopt the following principles:

### 1. Anonymization Logic is Self-Contained

- The `Anonymize` function is solely responsible for producing a clean, fully anonymized protobuf message.
- The marshalling logic **must not** contain special cases for anonymized data. It should be robust enough to handle data variations (like missing signatures) by default.

### 2. Anonymized Data Must Not Contain Original Raw Data

- The anonymization process **must** recursively traverse the protobuf message and strip all `raw_data` fields, setting them to empty byte slices.
- This ensures that a subsequent `Marshal` call cannot use raw data painting and is forced to serialize the file based purely on the anonymized semantic fields.

### 3. Anonymized Data Must Not Contain Signatures or Certificates

- The anonymization process **must** clear all fields containing digital signatures or certificates by setting them to `nil`.
- This prevents invalid signatures from being carried over into the anonymized file.

### 4. Marshalling Logic Must Gracefully Handle Missing Signatures

- The marshalling functions (`Append*`) must be updated to handle `nil` signature/certificate fields gracefully.
- **For TLV (Tag-Length-Value) structures:** If a signature field is `nil`, the marshaller will omit the entire TLV record for that signature.
- **For Fixed-Width structures:** If a signature field is `nil`, the marshaller will write a block of zero-bytes corresponding to the signature's fixed size and offset.

## Implementation Plan

### Phase 1: Update Anonymization Logic

1.  **Modify `internal/card/anonymize.go`:**
    - In `AnonymizeDriverCardFile`, after cloning, explicitly set certificate and signature fields to `nil`.
    - Add a recursive helper function to traverse the `cardv1.DriverCardFile` message and clear all `raw_data` fields.

2.  **Modify `internal/vu/anonymize.go`:**
    - In `AnonymizeVehicleUnitFile`, after cloning, explicitly set certificate and signature fields to `nil`.
    - Add a recursive helper function to traverse the `vuv1.VehicleUnitFile` message and clear all `raw_data` fields.

### Phase 2: Update Marshalling Logic

1.  **Review and update all `Append*` functions** related to signatures and certificates in `internal/dd`, `internal/card`, and `internal/vu`.
2.  Ensure each function checks for `nil` and implements the correct behavior (omit TLV record or write zero-bytes) as defined in the principles above.
3.  Verify that the `UseRawData: false` path in the marshaller is fully functional and correctly serializes all data structures from semantic fields alone. This is a prerequisite for the entire strategy to work.
