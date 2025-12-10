# Protobuf Code Generation for C#

This directory contains 157 `.proto` definition files that define the data structures for tachograph files.

## Current Status

The proto files have been kept in their original form but the generated Go code (`proto/gen/go/`) has been removed as part of the migration to .NET/C#.

## Generating C# Code

To generate C# code from these proto files, you need to:

### Option 1: Using buf (Recommended)

1. Install buf: https://buf.build/docs/installation
2. Update `buf.gen.yaml` to use C# plugin
3. Run: `buf generate`

However, the current proto files have dependencies on `buf.build/bufbuild/protovalidate` which need to be resolved first.

### Option 2: Using protoc directly

1. Install protoc (Protocol Buffers compiler)
2. Download buf/validate/validate.proto dependency
3. Run protoc for each .proto file:

```bash
protoc --proto_path=. \
       --proto_path=<path-to-buf-validate> \
       --csharp_out=../src/Tachograph/Proto/Generated \
       wayplatform/connect/tachograph/**/*.proto
```

### Option 3: Using Grpc.Tools NuGet package

The `Grpc.Tools` package is already installed in the `src/Tachograph` project. Once the buf/validate dependency is resolved, you can enable protobuf compilation by uncommenting the `<Protobuf>` ItemGroup in `src/Tachograph/Tachograph.csproj`.

## Dependencies

The proto files depend on:
- `buf.build/bufbuild/protovalidate:v1.0.0` - Proto validation rules

These dependencies are declared in `buf.yaml` but need to be made available for C# code generation.

## Migration Notes

- There are **157 proto definition files** in `wayplatform/connect/tachograph/`
- The original Go generated code was in `proto/gen/go/` (removed)
- C# generated code should go in `src/Tachograph/Proto/Generated/`
- The generated C# code will use the namespace `Tachograph.Proto.*`

## Next Steps

1. Resolve buf/validate dependencies
2. Set up automated proto code generation in CI/CD
3. Generate C# code from all 157 proto files
4. Integrate generated code into the Tachograph library
