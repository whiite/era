# Era

CLI tool for managing querying and manipulating date times written in
[Go](https://go.dev)

## Install

Currently not published to a Go registry (WIP) nor are there prebuilt binaries
so instead clone the repo and install using the following:

```bash
cd era
go install .
```

## Usage

Use `era help` to see available commands and flags otherwise here are some examples:

```bash
# Prints the current time
era now

# Prints the current time as a unix timestamp
era now --format unix
```

## Under consideration

- [ ] Better UX for locale names - (possibly using: https://pkg.go.dev/github.com/zlasd/tzloc)
