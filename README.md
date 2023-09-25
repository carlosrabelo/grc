# GRC — Gmail Rule Converter

Generate Atom XML feeds or Gmail import XML from YAML filter configs, and convert
Gmail Takeout XML back to YAML.

## Features

- **Atom XML** — Generate Atom feeds with categories, authors, and properties
- **Gmail XML** — Generate Gmail-compatible import XML via `grc gmail`
- **Import** — Convert Google Takeout XML back to YAML via `grc import`
- **Expression parser** — Parse Gmail search syntax (`from:`, `subject:`, `has:attachment`, …)
- **Search builder** — Build `mail.google.com/#search/…` URLs from criteria
- **Round-trip** — YAML → Gmail XML → YAML without data loss

## Installation

### Build from Source

```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

Install to `~/.local/bin` (no root required):

```bash
make install
```

### Using Go Install

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Usage

```
grc <command> [flags] <input> [output]

Commands:
  atom      Generate Atom XML feed from YAML
  gmail     Generate Gmail import XML from YAML
  import    Import Gmail filter XML back to YAML
  help      Show help (use "grc help <command>" for details)
  version   Show version

Common flags:
  --output PATH   Explicit output path
  --force         Overwrite existing output file
  --verbose       Enable verbose logging
  -h, --help      Show help for the current command
```

### Default output paths

| Command  | Default extension  | Example                              |
|----------|--------------------|--------------------------------------|
| `atom`   | `.xml`             | `cfg/example.yaml` → `example.xml`   |
| `gmail`  | `.gmail.xml`       | `cfg/example.yaml` → `example.gmail.xml` |
| `import` | `.yaml`            | `takeout.xml` → `takeout.yaml`       |

### Examples

```bash
# Atom feed
grc atom cfg/example.yaml

# Gmail import XML
grc gmail cfg/example-gmail.yaml

# Custom output path (positional)
grc gmail cfg/example-gmail.yaml my-filters.xml

# Same, via --output flag
grc gmail --output my-filters.xml cfg/example-gmail.yaml

# Import Gmail Takeout
grc import takeout.xml converted.yaml

# Overwrite existing output
grc gmail --force cfg/example-gmail.yaml
```

## YAML — Atom format

### Filter fields

| Field        | Type    | Description                            |
|--------------|---------|----------------------------------------|
| `name`       | string  | Filter name (becomes entry title)      |
| `expression` | string  | Gmail search expression                |
| `criteria`   | object  | Structured Gmail criteria              |
| `action`     | object  | **Required.** Must include `label`     |
| `category`   | object  | Atom category (term, scheme, label)    |
| `author`     | object  | Atom author (name, email, uri)         |
| `properties` | array   | Key-value pairs in entry summary       |

**Important**: every filter must specify `action.label` — the destination folder (Gmail label).
Filters are sorted alphabetically by this label in the generated output.

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

## Development

```bash
make build      # Compile binary to bin/grc
make test       # Run all tests
make quality    # fmt + vet + lint + test
make install    # Install to ~/.local/bin
make clean      # Remove build artifacts
```

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.
