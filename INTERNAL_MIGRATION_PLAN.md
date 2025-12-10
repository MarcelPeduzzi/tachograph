# Internal Files Migration Plan

This document outlines the plan for migrating the 202 Go files from the `internal/` directory to C#.

## Overview

The `internal/` directory contained complex binary parsing and processing logic organized into several packages:

- **internal/card** - 64 Go files - Card file processing
- **internal/vu** - 48 Go files - Vehicle Unit file processing  
- **internal/dd** - 76 Go files - Data Dictionary types
- **internal/cert** - 7 Go files - Certificate handling
- **internal/security** - 4 Go files - Cryptographic operations
- **internal/brainpool** - 3 Go files - Elliptic curve support

**Total: 202 implementation files** (~20,000+ lines of code)

## Migration Strategy

### Phase 1: Core Data Dictionary (Priority: HIGH)
Migrate the foundational parsing utilities in `internal/dd/`:

**Critical Files (Week 1-2):**
- [ ] `time.go` - TimeReal, BCDString, Date parsing
- [ ] `encoding.go` - String encoding/decoding with code pages
- [ ] `int24.go` - 24-bit integer support
- [ ] `bcd.go` - BCD (Binary Coded Decimal) encoding
- [ ] `bcd_string.go` - BCD string handling

**Core Types (Week 3-4):**
- [ ] `full_card_number.go` - Card number parsing
- [ ] `holder_name.go` - Name parsing with code pages
- [ ] `driver_identification.go` - Driver ID structures
- [ ] `geo_coordinates.go` - GPS coordinate parsing
- [ ] `vehicle_registration_identification.go` - VIN parsing

### Phase 2: Card File Processing (Priority: HIGH)
Migrate `internal/card/` for driver card support:

**Structure Parsing (Week 5-6):**
- [ ] `unmarshal.go` - TLV structure parsing
- [ ] `rawcardfile.go` - Raw card file handling
- [ ] `tlv.go` - Tag-Length-Value parser
- [ ] `marshal.go` - Card serialization

**Application Files (Week 7-8):**
- [ ] `application_identification.go` - App ID (Gen1)
- [ ] `application_identification_g2.go` - App ID (Gen2)
- [ ] `ic.go` / `icc.go` - IC/ICC data
- [ ] `identification.go` - Card identification

**Activity Data (Week 9-10):**
- [ ] `activity.go` - Driver activities
- [ ] `places.go` / `places_g2.go` - Place records
- [ ] `events.go` / `faults.go` - Events and faults
- [ ] `specific_conditions.go` - Specific conditions

### Phase 3: Vehicle Unit Processing (Priority: MEDIUM)
Migrate `internal/vu/` for VU file support:

**Core Parsing (Week 11-12):**
- [ ] `unmarshal.go` - TV structure parsing
- [ ] `trep.go` - TREP format handling
- [ ] `marshal.go` - VU serialization

**VU Data (Week 13-14):**
- [ ] `overview.go` / `overview_gen*.go` - Overview data
- [ ] `activities.go` / `activities_gen*.go` - Activity records
- [ ] `events_faults.go` / `events_faults_gen*.go` - Events/faults
- [ ] `technical_data.go` / `technical_data_gen*.go` - Technical data

### Phase 4: Security & Certificates (Priority: MEDIUM)
Migrate `internal/cert/` and `internal/security/`:

**Certificate Handling (Week 15):**
- [ ] `cert/embedded.go` - Embedded certificates
- [ ] `cert/certcache/*.go` - Certificate cache (Gen1, Gen2, root)
- [ ] `cert/chain.go` - Certificate chain validation
- [ ] `cert/resolver.go` - Certificate resolution

**Cryptographic Operations (Week 16):**
- [ ] `security/signature.go` - Signature verification
- [ ] `security/rsa.go` - RSA operations
- [ ] `security/ecdsa.go` - ECDSA operations
- [ ] `brainpool/*.go` - Brainpool elliptic curves

### Phase 5: Additional Features (Priority: LOW)
Complete remaining functionality:

**Week 17-18:**
- [ ] `card/anonymize.go` - Card anonymization
- [ ] `vu/anonymize.go` - VU anonymization
- [ ] `card/authenticate.go` - Card authentication
- [ ] `vu/authenticate.go` - VU authentication
- [ ] `card/parse.go` / `vu/parse.go` - Semantic parsing

## Implementation Guidelines

### Code Translation Patterns

#### Go to C# Type Mappings
- `[]byte` → `byte[]` or `Span<byte>` or `ReadOnlySpan<byte>`
- `uint8` → `byte`
- `uint16` → `ushort`
- `uint32` → `uint`
- `uint64` → `ulong`
- `int64` → `long`
- `string` → `string`
- `error` → exceptions or `bool` with `out` parameters

#### Binary Parsing
- Go's `binary.BigEndian.Uint32(data[offset:])` → C#'s `BinaryPrimitives.ReadUInt32BigEndian(data.AsSpan(offset))`
- Go's `bufio.Scanner` → C# custom parsing or `BinaryReader`
- Go's slicing `data[offset:offset+length]` → C# `data.AsSpan(offset, length)` or `data.Skip(offset).Take(length)`

#### Error Handling
- Go's `if err != nil { return nil, err }` → C# `throw new InvalidDataException()`
- Go's multiple return values → C# exceptions or `Try` pattern

### Testing Requirements

Each migrated file must have:
1. Unit tests matching the original Go tests
2. Golden file tests with real .DDD files
3. Round-trip tests (parse → marshal → parse)

### Documentation

Each C# class must include:
1. XML documentation comments
2. References to EU regulation sections
3. ASN.1 definitions (where applicable)
4. Examples of usage

## Current Status

✅ **Completed:**
- Placeholder classes for all internal packages
- Project structure and build system
- CLI application framework

⏳ **In Progress:**
- Protobuf C# code generation setup

❌ **Not Started:**
- Internal implementation files (0 of 202 migrated)

## Estimated Timeline

- **Protobuf Generation:** 1 week (setup and dependencies)
- **Phase 1 (DD Core):** 4 weeks
- **Phase 2 (Card Processing):** 6 weeks  
- **Phase 3 (VU Processing):** 4 weeks
- **Phase 4 (Security):** 2 weeks
- **Phase 5 (Additional):** 2 weeks
- **Testing & Refinement:** 2 weeks

**Total Estimated Time: ~21 weeks (5 months)** for a complete, production-ready migration.

## Recommended Approach

Given the scope, we recommend:

1. **Incremental Migration:** Start with Phase 1 to establish patterns
2. **Parallel Protobuf:** Set up protobuf generation in parallel
3. **Test-Driven:** Write tests before implementing each component
4. **Review Cycles:** Code review after each phase
5. **Documentation:** Document as you go, not at the end

## Alternative: Hybrid Approach

Consider keeping the library architecture in C# but implementing critical parsing logic in a more automated way:
- Use code generation tools to translate Go to C#
- Focus manual effort on complex algorithms and crypto
- Leverage existing .NET libraries where possible (e.g., BouncyCastle for Brainpool curves)
