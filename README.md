# Deduplication Tool (CASD)

## Description

CLI tool to process files in chunks, store unique chunks, and save a 'formula' (sequence of hashes) that represents the file.

How it works overview
- Reads the file in 32KB chunks.
- Calculates the hash of each chunk.
- Stores each unique chunk.
- Saves a `Formula` with the hashes and returns an identifier (hash of the sequence).
- Shows progress in the terminal and, at the end, prints `identifier: <hash>`.

## Install
Clone from repository
```bash
git clone https://github.com/informeai/casd
```

## Build and Running

To run without building

```bash
go run ./cmd/main.go deduplicate --file /path/to/file
```

To compile and run the binary:

```bash
go build -o casd ./cmd
./casd deduplicate -f /path/to/file
```

Help command
```bash
casd --help
```