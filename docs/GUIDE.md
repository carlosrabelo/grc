# GRC — Gmail Rule Converter Guide

## Overview

GRC converts YAML filter definitions into:

- **Atom XML** feeds (for RSS readers, dashboards)
- **Gmail import XML** (compatible with Google Workspace / Gmail filters import)
- **YAML config** from Google Takeout Gmail filter XML (`--import`)

## Installation

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
sudo cp bin/grc /usr/local/bin/
```

Via `go install`:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Command-line Interface

```
grc <command> [flags] <input> [output]
```

GRC is organised around subcommands. Each command takes a positional input
file and an optional positional output path.

### Commands

| Command  | Description                                                |
|----------|------------------------------------------------------------|
| `atom`   | Generate Atom XML feed from YAML                           |
| `gmail`  | Generate Gmail-compatible import XML from YAML             |
| `import` | Convert Gmail filter XML (Google Takeout) back to YAML     |
| `help`   | Show help (`grc help <command>` for command-specific help) |
| `version`| Show version and build information                         |

### Common flags (available on `atom`, `gmail`, `import`)

| Flag             | Description                                            |
|------------------|--------------------------------------------------------|
| `--output PATH`  | Explicit output path (overrides the positional output) |
| `--force`        | Overwrite an existing output file                      |
| `--verbose`      | Enable verbose logging                                 |
| `-h, --help`     | Show help for the current command                      |

### Default output paths

| Command  | Default extension | Example                                |
|----------|-------------------|----------------------------------------|
| `atom`   | `.xml`            | `cfg/example.yaml` → `example.xml`     |
| `gmail`  | `.gmail.xml`      | `cfg/example.yaml` → `example.gmail.xml` |
| `import` | `.yaml`           | `takeout.xml` → `takeout.yaml`         |

### Exit codes

| Code | Meaning       |
|------|---------------|
| 0    | Success       |
| 1    | Runtime error |
| 2    | Usage error   |

## YAML Configuration

Two configuration formats are supported:

### 1. Atom format (expression-based)

Uses Gmail search-like expressions and Atom-specific metadata.

```yaml
defaults:
  author:
    name: Default Author
    email: author@example.com
  category:
    term: email
    scheme: http://example.com/categories
    label: Email Category

filters:
  - name: Newsletter
    expression: subject:(newsletter OR weekly OR digest)
    category:
      term: newsletter
      label: Newsletters
    author:
      name: Newsletter Bot
      email: newsletter@example.com
    properties:
      - name: priority
        value: high
```

#### Atom fields

| Field        | Type    | Required | Description                           |
|--------------|---------|----------|---------------------------------------|
| `name`       | string  | yes      | Filter name (becomes entry title)     |
| `expression` | string  | yes*     | Gmail search expression               |
| `criteria`   | object  | yes*     | Structured Gmail criteria (alt to expression) |
| `action`     | object  | no       | Gmail actions to apply                |
| `category`   | object  | no       | Atom category (term, scheme, label)   |
| `author`     | object  | no       | Atom author (name, email, uri)        |
| `properties` | array   | no       | Key-value pairs in entry summary      |

\* At least one of `expression` or `criteria` is required.

### 2. Gmail format (criteria + action)

Uses structured fields that map directly to Gmail filter properties.

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@example.com
    action:
      label: Newsletters
      markRead: true
```

#### Criteria fields

| Field          | Type    | Description                          |
|----------------|---------|--------------------------------------|
| `from`         | string  | Sender email                         |
| `to`           | string  | Recipient email                      |
| `subject`      | string  | Subject line                         |
| `hasTheWord`   | string  | Free-text search terms               |
| `hasAttachment` | bool   | Has attachment                       |
| `is`           | string  | `unread`, `starred`, `snoozed`, `important`, `chat` |
| `in`           | string  | `inbox`, `trash`, `spam`, `anywhere`, `sent`, `drafts` |
| `cc`           | string  | CC recipient                         |
| `bcc`          | string  | BCC recipient                        |
| `deliveredto`  | string  | Delivered-to address                 |
| `list`         | string  | Mailing list address                 |
| `filename`     | string  | Attachment filename                  |
| `after`        | string  | Date after (e.g., `2024/01/01`)      |
| `before`       | string  | Date before                          |
| `older_than`   | string  | Relative age (e.g., `30d`, `2y`)     |
| `newer_than`   | string  | Relative age (e.g., `7d`)            |
| `size`         | string  | Exact size (e.g., `10M`)             |
| `larger`       | string  | Larger than (e.g., `5M`)             |
| `smaller`      | string  | Smaller than (e.g., `1M`)            |
| `label`        | string  | Has label                            |
| `category`     | string  | Gmail category (`primary`, `social`, `promotions`, `updates`, `forums`) |

#### Action fields

| Field         | Type    | Description                      |
|---------------|---------|----------------------------------|
| `archive`     | bool    | Skip inbox                       |
| `markRead`    | bool    | Mark as read                     |
| `star`        | bool    | Star the message                 |
| `delete`      | bool    | Send to trash                    |
| `neverSpam`   | bool    | Never send to spam               |
| `alwaysSpam`  | bool    | Always mark as spam              |
| `forward`     | string  | Forward to email address         |
| `label`       | string  | Apply label                      |
| `smartLabel`  | string  | Apply smart label category       |
| `moveToTrash` | bool    | Move to trash                    |

## Usage examples

### Generate Atom feed

```bash
grc atom cfg/example.yaml
# Produces example.xml
```

### Generate Gmail import XML

```bash
grc gmail cfg/example-gmail.yaml
# Produces example.gmail.xml
```

### Import Gmail filters from Google Takeout

```bash
grc import takeout-filters.xml
# Produces takeout-filters.yaml
```

### Force overwrite

```bash
grc gmail --force cfg/example-gmail.yaml -o output.gmail.xml
```

### Verbose mode

```bash
grc gmail --verbose cfg/example-gmail.yaml
```

### Full round-trip

```bash
# YAML → Gmail XML
grc gmail cfg/example-gmail.yaml filters.gmail.xml

# Gmail XML → YAML
grc import filters.gmail.xml
```

## Expression parser

When using `expression` (Atom format), GRC parses Gmail search syntax:

```
from:newsletter@example.com subject:(weekly OR digest) has:attachment
```

Supported operators:

| Operator            | Maps to                      |
|---------------------|------------------------------|
| `from:`             | Criteria.From                |
| `to:`               | Criteria.To                  |
| `subject:`          | Criteria.Subject             |
| `has:`              | `has:attachment` → HasAttachment=true |
| `is:`               | Criteria.Is                  |
| `in:`               | Criteria.In                  |
| `cc:`               | Criteria.Cc                  |
| `bcc:`              | Criteria.Bcc                 |
| `deliveredto:`      | Criteria.DeliveredTo         |
| `list:`             | Criteria.List                |
| `filename:`         | Criteria.Filename            |
| `after:`            | Criteria.After               |
| `before:`           | Criteria.Before              |
| `older_than:`       | Criteria.OlderThan           |
| `newer_than:`       | Criteria.NewerThan           |
| `size:`             | Criteria.Size                |
| `larger:`           | Criteria.Larger              |
| `smaller:`          | Criteria.Smaller             |
| `label:`            | Criteria.Label               |
| `category:`         | Criteria.Category            |
| Unrecognized tokens | Criteria.HasTheWord          |

## Search query builder

GRC can build `mail.google.com` search URLs from criteria:

- `query.BuildSearchString(criteria)` → `from:user@example.com subject:hello`
- `query.BuildSearchURL(criteria)` → `https://mail.google.com/mail/u/0/#search/from%3Auser%40example.com%20subject%3Ahello`

## Examples

- `cfg/example.yaml` — Atom-oriented config with 6 expression-based filters
- `cfg/example-gmail.yaml` — Gmail-oriented config with 16 filters covering all criteria and action types

## Development

```bash
make build        # Compile to bin/grc
make test         # Run all tests (12 packages)
make quality      # Format, vet, lint, test
make clean        # Remove bin/ and .make/ stamps

# Lint (requires golangci-lint)
make lint
```
