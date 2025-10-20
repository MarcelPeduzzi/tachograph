# Tachograph Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/way-platform/tachograph-go)](https://pkg.go.dev/github.com/way-platform/tachograph-go)
[![GoReportCard](https://goreportcard.com/badge/github.com/way-platform/tachograph-go)](https://goreportcard.com/report/github.com/way-platform/tachograph-go)
[![CI](https://github.com/way-platform/tachograph-go/actions/workflows/release.yaml/badge.svg)](https://github.com/way-platform/tachograph-go/actions/workflows/release.yaml)

A Go SDK and CLI tool for working with Tachograph data (.DDD files).

## Specification

This SDK implements parsing of downloaded tachograph data, according to [the requirements for the construction, testing, installation, operation and repair of tachographs and their components](https://eur-lex.europa.eu/eli/reg_impl/2016/799/oj/eng).

## Features

This library provides a comprehensive set of tools for working with Tachograph data, from raw binary parsing to anonymization and serialization.

```go
// 1. Unmarshal the raw .DDD file content into a RawFile object.
rawFile, err := tachograph.Unmarshal(dddBytes)
if err != nil {
    log.Fatalf("failed to unmarshal: %v", err)
}

// 2. Authenticate the signatures within the RawFile.
authenticatedRawFile, err := tachograph.Authenticate(context.Background(), rawFile)
if err != nil {
    log.Fatalf("failed to authenticate: %v", err)
}

// 3. Parse the authenticated RawFile into a File.
parsedFile, err := tachograph.Parse(authenticatedRawFile)
if err != nil {
    log.Fatalf("failed to parse: %v", err)
}

// 4. Anonymize personal data in the parsed File.
anonymizedFile, err := tachograph.Anonymize(parsedFile)
if err != nil {
    log.Fatalf("failed to anonymize: %v", err)
}

// 5. Marshal the file back into the binary .DDD format.
marshalledBytes, err := tachograph.Marshal(anonymizedFile)
if err != nil {
    log.Fatalf("failed to marshal: %v", err)
}
```

### Unmarshalling

[`tachograph.Unmarshal`](https://pkg.go.dev/github.com/way-platform/tachograph-go#Unmarshal) parses raw `.DDD` byte arrays into a [`tachographv1.RawFile`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#RawFile) object.

### Authenticating

[`tachograph.Authenticate`](https://pkg.go.dev/github.com/way-platform/tachograph-go#Authenticate) cryptographically verifies the signatures within a [`tachographv1.RawFile`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#RawFile).

### Parsing

[`tachograph.Parse`](https://pkg.go.dev/github.com/way-platform/tachograph-go#Parse) turns a [`tachographv1.RawFile`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#RawFile) into a [`tachographv1.File`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#File), with meaningful, high-level structures like driver activities and events, making the data easy to analyze.

### Anonymizing

[`tachograph.Anonymize`](https://pkg.go.dev/github.com/way-platform/tachograph-go#Anonymize) removes or obscures personal data from a [`tachographv1.File`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#File) - making the data usable for unit testing

### Marshalling

[`tachograph.Marshal`](https://pkg.go.dev/github.com/way-platform/tachograph-go#Marshal) serializes a [`tachographv1.File`](https://pkg.go.dev/github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1#File) back into the binary `.DDD` format.

## Alternatives

This SDK draws inspiration from other tachograph SDKs, including:

- [traconiq/tachoparser](https://github.com/traconiq/tachoparser)
- [jugglingcats/tachograph-reader](https://github.com/jugglingcats/tachograph-reader)
