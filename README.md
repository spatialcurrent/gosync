[![CircleCI](https://circleci.com/gh/spatialcurrent/gosync/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/gosync/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/gosync)](https://goreportcard.com/report/spatialcurrent/gosync)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/gosync?status.svg)](https://godoc.org/github.com/spatialcurrent/gosync) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/gosync/blob/master/LICENSE)

# gosync

## Description

**gosync** is a simple command line program for synchronizing a source and destination, including AWS S3 buckets.  **gosync** supports the following operating systems and architectures.

## Platforms

The following platforms are supported.  Pull requests to support other platforms are welcome!

| GOOS | GOARCH |
| ---- | ------ |
| darwin | amd64 |
| linux | amd64 |
| windows | amd64 |
| linux | arm64 |

## Releases

Find releases at [https://github.com/spatialcurrent/gosync/releases](https://github.com/spatialcurrent/gosync/releases).  See the **Building** section below to build from scratch.

**Darwin**

- `gosync_darwin_amd64` - CLI for Darwin on amd64 (includes `macOS` and `iOS` platforms)

**Linux**

- `gosync_linux_amd64` - CLI for Linux on amd64
- `gosync_linux_arm64` - CLI for Linux on arm64

**Windows**

- `gosync_windows_amd64.exe` - CLI for Windows on amd64

## Usage

See the usage below or the following examples.

```shell
gosync is a super simple utility to print environment variables in a custom format.
Supports the following formats: csv, bson, go, gob, json, properties, tags, toml, tsv, yaml.

Usage:
  gosync [-f FORMAT] [flags] [variable]...

Flags:
  -f, --format string   output format, one of: csv, bson, go, gob, json, properties, tags, toml, tsv, yaml (default "properties")
  -h, --help            help for gosync
  -0, --null            use a NUL byte to end each line instead of a newline character
  -p, --pretty          use pretty output format
  -r, --reversed        if output is sorted, sort in reverse alphabetical order
  -s, --sorted          sort output
```

# Examples

**Sync Local Directories**

```shell
gosync /path/to/source/dir /path/to/destination/directory
```

**Upload to AWS S3**

`gosync` can upload a series of files to an AWS S3 bucket with a given prefix.

```shell
gosync file://path/to/local/dir s3://bucket/key/prefix
```

**Download to AWS S3**

`gosync` can download a series of files from AWS S3 to a local directory.

```shell
gosync s3://bucket/key/prefix file://path/to/local/dir
```

## Building

Use `make help` to see help information for each target.

**CLI**

The `make build_cli` script is used to build executables for Linux and Windows.  Use `make install` for standard installation as a go executable.

## Testing

**CLI**

To run CLI testes use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```
TMPDIR="/usr/local/tmp" make test_cli
```

**Go**

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

## Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/gosync/blob/master/CONTRIBUTING.md) for how to get started.

## License

This work is distributed under the **MIT License**.  See **LICENSE** file.
