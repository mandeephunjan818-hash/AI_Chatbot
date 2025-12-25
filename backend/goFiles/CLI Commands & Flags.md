# CLI Commands & Flags

The Go backend includes a Command Line Interface (CLI) for managing startup behavior and running diagnostics.

## Usage Syntax

```bash
go run . [flags]
```

## Available Flags

The following flags can be passed to the application binary or run command:

| Flag | Default | Definition |
|------|---------|------------|
| `-help` | false | Displays the help message, usage instructions, and exits. |
| `-test` | false | **Full Test Mode.** Runs the comprehensive database test suite (CRUD). Exits with code 0 (success) or 1 (failure). |
| `-quick` | false | **Quick Check.** Performs a simple ping to verify DB connectivity and exits. |
| `-skip-test` | false | **Skip Verification.** Bypasses the initial connection test before starting the server (Not recommended). |

## Usage Examples

### Start Server (Standard)

Checks database connection first, then starts HTTP server.

```bash
go run .
```

### Run Deep Diagnostics

Executes insert, read, and delete operations to ensure database write access.

```bash
go run . -test
```

