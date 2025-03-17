# Era

CLI tool for managing querying and manipulating date times written in
[Go](https://go.dev)

[Mirrored on GitHub at `whiite/era`](https://github.com/whiite/era) from [GitLab at `monokuro/era`](https://gitlab.com/monokuro/era)

## Install

Installation currently requires you have `go` installed on your system

To install:

```bash
go install gitlab.com/monokuro/era@latest
```

## Usage

Use `era help` to see available commands and flags otherwise here are some examples:

```bash
# Prints the current time
era now

# Prints the current time as a unix timestamp
era now --formatter unix

# Prints the current date and time using Luxon and Moment formatting tokens in Tokyo's time zone
era now --timezone Asia/Tokyo --formatter luxon "h:mm d/L/yyyy"
era now --timezone Asia/Tokyo --formatter moment "h:mm D/M/Y"

# Prints the available supported tokens and descriptions for the strptime/strftime formatter
era tokens --formatter strftime

# Prints help for the specific CLI command
era help <sub command>
```

## Supported Formatters

More is planned to be added in the future as I come across them but these three
I've found I've used the most.

Correct escape sequences are supported for all

**Note: Not every token is implemented fully yet. Most work as intended; some locale
specific formats are hardcoded to the UK or English versions but all are covered by tests**

- [moment](https://momentjs.com)
  - Some locale specific formats may return slightly different strings to the real moment
- [luxon](https://moment.github.io/luxon/#/)
  - Missing support for tokens:
    - `ZZZZ` - full offset name e.g. `"Eastern Standard Time"`
    - `TTTT` - Localised 24 hour time with full time zone name
    - `f`, `ff`, `fff`, `ffff` - localised date and time
    - `F`, `FF`, `FFF`, `FFFF` - localised date and time with seconds
- [strftime](https://linux.die.net/man/3/strftime) (tokens used in a variety of languages including the `date` CLI)
  - Full compatibility via C FFI bindings to the `strftime` function
  - An alternative Go implementation (using `go:strftime` as the `formatter`)
    - Missing true locale support with some tokens such as `%Ec`, `%c` currently hardcoded to the UK representation
- [go](https://pkg.go.dev/time) (time package format)
  - Full support as this CLI tool is written in Go and uses the standard library time package

## Compatibility table

| Feature                | Go     | strftime/strptime | Go strftime/strptime | Luxon    | Moment |
| ---------------------- | ------ | ----------------- | -------------------- | -------- | ------ |
| Formatting with tokens | All ✅ | All ✅            | All ✅               | Most ⚠️  | All ✅ |
| Token descriptions     | All ✅ | All ✅            | All ✅               | All ✅⚠️ | All ✅ |
| Locale support         | N/A    | Yes ✅            | Some                 | Some⚠️   | Some   |
| Parsing tokens         | All ✅ | All ✅            | Some                 | No ❌⚠️  | No ❌  |

## Under consideration

- [ ] Better UX for locale names - (possibly using: https://pkg.go.dev/github.com/zlasd/tzloc)
- [ ] A "describe"/"explain" command for describing the tokens provided according to the specified formatter
