# GRC — Gmail Rule Converter

Generate Atom XML feeds or Gmail import XML from YAML filter configs, and convert
Gmail Takeout XML back to YAML.

## Features

- **Atom XML** — Generate Atom feeds with categories, authors, and properties
- **Gmail XML** — Generate Gmail-compatible import XML (`--format gmail`)
- **Import** — Convert Google Takeout XML back to YAML (`--import`)
- **Expression parser** — Parse Gmail search syntax (`from:`, `subject:`, `has:attachment`, …)
- **Search builder** — Build `mail.google.com/#search/…` URLs from criteria
- **Round-trip** — YAML → Gmail XML → YAML without data loss

## Install

```bash
make build
sudo cp bin/grc /usr/local/bin/
```

Or:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Usage

```
grc [options] <input> [output]
grc --import <gmail.xml> [output]
```

| Flag              | Description                                         |
|-------------------|-----------------------------------------------------|
| `--format atom`   | Atom XML output (default)                           |
| `--format gmail`  | Gmail import XML output                             |
| `--import`        | Convert Gmail XML → YAML                            |
| `--force`         | Overwrite existing output file                      |
| `--verbose`       | Enable detailed logging                             |
| `--version`       | Show version                                        |
| `--help`          | Show help                                           |

### Examples

```bash
# Atom feed (default)
grc cfg/example.yaml

# Gmail import XML
grc --format gmail cfg/example-gmail.yaml

# Custom output path
grc --format gmail cfg/example-gmail.yaml my-filters.xml

# Import Gmail Takeout
grc --import takeout.xml converted.yaml
```

## YAML — Atom format

```yaml
defaults:
  author:
    name: Default Author
    email: author@example.com
  category:
    term: default
    label: Default Category

filters:
  - name: Newsletter
    expression: subject:newsletter
    category:
      term: newsletter
      label: Newsletters
    properties:
      - name: priority
        value: high
```

See `cfg/example.yaml`.

## YAML — Gmail format

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@example.com
    action:
      label: Newsletters
      markRead: true

  - name: Spam
    criteria:
      from: spam@example.com
      hasAttachment: true
    action:
      alwaysSpam: true
```

Criteria fields: `from`, `to`, `subject`, `hasTheWord`, `hasAttachment`,
`is`, `in`, `cc`, `bcc`, `deliveredto`, `list`, `filename`,
`after`, `before`, `older_than`, `newer_than`,
`size`, `larger`, `smaller`, `label`, `category`.

Action fields: `archive`, `markRead`, `star`, `delete`, `neverSpam`,
`alwaysSpam`, `forward`, `label`, `smartLabel`, `moveToTrash`.

See `cfg/example-gmail.yaml` for 16 sample filters.

## Build

```bash
make build        # bin/grc
make test         # all tests
make quality      # fmt + vet + lint + test
make clean        # remove artifacts
```
