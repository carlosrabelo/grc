# GRC — Gmail Rule Converter Guide

## Overview

GRC converts YAML filter definitions into Atom XML feeds.

## Installation

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
```

## Command-line Interface

```
grc <command> [flags] <input> [output]
```

### Commands

| Command   | Description                     |
|-----------|---------------------------------|
| `atom`    | Generate Atom XML feed from YAML |
| `help`    | Show help                       |
| `version` | Show version                    |

### Common flags

| Flag             | Description                    |
|------------------|--------------------------------|
| `--output PATH`  | Explicit output path           |
| `--force`        | Overwrite an existing file     |
| `--verbose`      | Enable verbose logging         |
| `-h, --help`     | Show help                      |

### Default output paths

| Command | Default extension | Example                              |
|---------|-------------------|--------------------------------------|
| `atom`  | `.xml`            | `docs/example.yaml` → `example.xml`  |

## YAML Configuration

```yaml
defaults:
  author:
    name: Default Author
    email: author@example.com
  category:
    term: email
    label: Email Category

filters:
  - name: Newsletter
    expression: subject:(newsletter OR weekly OR digest)
    category:
      term: newsletter
      label: Newsletters
```

### Atom fields

| Field        | Type    | Required | Description                         |
|--------------|---------|----------|-------------------------------------|
| `name`       | string  | yes      | Filter name (becomes entry title)   |
| `expression` | string  | yes      | Gmail search expression             |
| `category`   | object  | no       | Atom category (term, scheme, label) |
| `author`     | object  | no       | Atom author (name, email, uri)      |
| `properties` | array   | no       | Key-value pairs in entry summary    |

## Usage examples

```bash
grc atom docs/example.yaml
```

See `docs/example.yaml` for a complete Atom config.

## Development

```bash
make build        # Compile to bin/grc
make test         # Run all tests
make quality      # Format, vet, test
```
