# pkg/shared

## Overview

`pkg/shared` hosts artifacts that are intentionally shared across services in the Hangout-Planner repository. It is the canonical place for small, stable contracts and helper types that multiple services consume. Typical contents include:

- Protocol Buffer definitions and generated code used for gRPC communication
- Shared enums, small value objects, and typed constants
- Lightweight utilities and types that are implementation-agnostic

## Goals

- Provide a minimal, well-documented surface area for cross-service contracts
- Avoid embedding business logic or service-specific code
- Keep generated artifacts clearly separated from source files

## Contents and structure

- `pkg/shared/proto/` — source `.proto` files organized by logical package (for example `file/`)
- `pkg/shared/proto/gen/` — generated Go code produced by `protoc` (kept separate from sources)
- `pkg/shared/enums/`, `pkg/shared/types/` — small Go packages with shared enums and helper types

## Protobuf and gRPC generation

Prerequisites

- `protoc` (protocol buffer compiler)
- `protoc-gen-go` and `protoc-gen-go-grpc` Go plugins

Install the Go plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Generate code (project root or `pkg/shared`):

```bash
cd pkg/shared
make proto
```

The `Makefile` target runs `protoc` and writes generated files into `pkg/shared/proto/gen` using `paths=source_relative` so generated files are organized next to their source package layout.
