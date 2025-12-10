# Tachograph

[![CI](https://github.com/MarcelPeduzzi/tachograph/actions/workflows/release.yaml/badge.svg)](https://github.com/MarcelPeduzzi/tachograph/actions/workflows/release.yaml)

A .NET 10 SDK for working with Tachograph data (.DDD files).

## Specification

This SDK implements parsing of downloaded tachograph data, according to [the requirements for the construction, testing, installation, operation and repair of tachographs and their components](https://eur-lex.europa.eu/eli/reg_impl/2016/799/oj/eng).

## Features

This library provides a comprehensive set of tools for working with Tachograph data, from raw binary parsing to anonymization and serialization.

```csharp
// 1. Unmarshal the raw .DDD file content into a RawFile object.
var rawFile = Tachograph.Unmarshal(dddBytes);

// 2. Authenticate the signatures within the RawFile.
var authenticatedRawFile = await Tachograph.AuthenticateAsync(rawFile);

// 3. Parse the authenticated RawFile into a File.
var parsedFile = Tachograph.Parse(authenticatedRawFile);

// 4. Anonymize personal data in the parsed File.
var anonymizedFile = Tachograph.Anonymize(parsedFile);

// 5. Marshal the file back into the binary .DDD format.
var marshalledBytes = Tachograph.Marshal(anonymizedFile);
```

### Unmarshalling

`Tachograph.Unmarshal` parses raw `.DDD` byte arrays into a `RawFile` object.

```csharp
var rawFile = Tachograph.Unmarshal(dddBytes);
```

Or with custom options:

```csharp
var opts = new UnmarshalOptions { Strict = false };
var rawFile = opts.Unmarshal(dddBytes);
```

### Authenticating

`Tachograph.AuthenticateAsync` cryptographically verifies the signatures within a `RawFile`.

```csharp
var authenticatedRawFile = await Tachograph.AuthenticateAsync(rawFile);
```

Or with custom options:

```csharp
var opts = new AuthenticateOptions { Mutate = true };
var authenticatedRawFile = await opts.AuthenticateAsync(rawFile);
```

### Parsing

`Tachograph.Parse` turns a `RawFile` into a `File`, with meaningful, high-level structures like driver activities and events, making the data easy to analyze.

```csharp
var parsedFile = Tachograph.Parse(authenticatedRawFile);
```

Or with custom options:

```csharp
var opts = new ParseOptions { PreserveRawData = false };
var parsedFile = opts.Parse(authenticatedRawFile);
```

### Anonymizing

`Tachograph.Anonymize` removes or obscures personal data from a `File` - making the data usable for unit testing.

```csharp
var anonymizedFile = Tachograph.Anonymize(parsedFile);
```

Or with custom options:

```csharp
var opts = new AnonymizeOptions 
{ 
    PreserveTimestamps = true,
    PreserveDistanceAndTrips = true
};
var anonymizedFile = opts.Anonymize(parsedFile);
```

### Marshalling

`Tachograph.Marshal` serializes a `File` back into the binary `.DDD` format.

```csharp
var marshalledBytes = Tachograph.Marshal(anonymizedFile);
```

Or with custom options:

```csharp
var opts = new MarshalOptions { UseRawData = false };
var marshalledBytes = opts.Marshal(anonymizedFile);
```

## Installation

```bash
dotnet add package Tachograph
```

## Building

```bash
dotnet build
```

## Testing

```bash
dotnet test
```

## Requirements

- .NET 10.0 or later

## Alternatives

This SDK draws inspiration from other tachograph SDKs, including:

- [traconiq/tachoparser](https://github.com/traconiq/tachoparser)
- [jugglingcats/tachograph-reader](https://github.com/jugglingcats/tachograph-reader)
- [way-platform/tachograph-go](https://github.com/way-platform/tachograph-go)
