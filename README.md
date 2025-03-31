# memcached-checker

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/aethiopicuschan/memcached-checker)](https://goreportcard.com/report/github.com/aethiopicuschan/memcached-checker)
[![CI](https://github.com/aethiopicuschan/memcached-checker/actions/workflows/ci.yaml/badge.svg)](https://github.com/aethiopicuschan/memcached-checker/actions/workflows/ci.yaml)

Tools to verify the operation of memcached and its compatible servers.

## Installation

```bash
go install github.com/aethiopicuschan/memcached-checker@latest
```

## Usage

```bash
# Benchmark
❯ memcached-checker benchmark -h
Run benchmarks for Set, Get, and Del.

Usage:
  memcached-checker benchmark [flags]

Flags:
  -a, --address string    Address of the memcached server (default "127.0.0.1:11211")
  -c, --concurrency int   Number of concurrent workers (default 10)
  -h, --help              help for benchmark
  -r, --requests int      Request count per worker (default 10000)

# Check functions
❯ memcached-checker check -h
Check that it operates correctly as both a memcached and a compatible server.

Usage:
  memcached-checker check [flags]

Flags:
  -a, --address string   Address of the memcached server (default "127.0.0.1:11211")
  -f, --flush            Flush the memcached server before running the check
  -h, --help             help for check
```

## Checklist

Commands to check with the `check` command.

### Supported

- [x] Ping
- [x] Flush
- [x] Set
- [x] Add
- [x] Replace
- [x] Get
- [x] Gets
- [x] Append
- [x] Prepend
- [x] Increment
- [x] Decrement
- [x] Touch

### TODO

- [ ] CAS
- [ ] Stats
- [ ] Version
- [ ] Quit
- [ ] Verbosity
