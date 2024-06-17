# GRC — Gmail Rule Converter

Generates Atom XML feeds or Gmail import XML from YAML filter configs, and converts Gmail Takeout XML back to YAML.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://go.dev/)

## Highlights

- Generate Atom feeds with categories, authors, and properties
- Generate Gmail-compatible import XML via `grc gmail`
- Convert Google Takeout XML back to YAML via `grc import`
- Parse Gmail search syntax (`from:`, `subject:`, `has:attachment`, …)
- Build `mail.google.com/#search/…` URLs from criteria
- Round-trip YAML → Gmail XML → YAML without data loss

## Prerequisites

- **Go 1.22+** — required to build from source; [download](https://go.dev/dl/)

## Installation

### Build from Source

```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

Install to `~/.local/bin` (default), or system-wide to `/usr/local/bin` (sudo only for the copy):

```bash
make install
make install-system
make uninstall
make uninstall-system
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
| `atom`   | `.xml`             | `docs/example.yaml` → `example.xml`   |
| `gmail`  | `.gmail.xml`       | `docs/example.yaml` → `example.gmail.xml` |
| `import` | `.yaml`            | `takeout.xml` → `takeout.yaml`       |

### Examples

```bash
# Atom feed
grc atom docs/example.yaml

# Gmail import XML
grc gmail docs/example-gmail.yaml

# Custom output path (positional)
grc gmail docs/example-gmail.yaml my-filters.xml

# Same, via --output flag
grc gmail --output my-filters.xml docs/example-gmail.yaml

# Import Gmail Takeout
grc import takeout.xml converted.yaml

# Overwrite existing output
grc gmail --force docs/example-gmail.yaml
```

## Configuration

### Atom format

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

See `docs/example.yaml`.

### Gmail format

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

See `docs/example-gmail.yaml` for 16 sample filters.

## Project Layout

```
grc/cmd/grc/         # Go entry point
grc/internal/        # Private packages
grc/pkg/types/       # Shared public types
bin/                 # Compiled binaries (git-ignored)
.make/               # Build and install scripts
docs/                # Examples and guides
```

## Development

```bash
make build           # Compile binary to bin/grc
make test            # Run all tests
make quality         # Format, vet, and test
make install         # Install binary to ~/.local/bin
make install-system  # Install binary to /usr/local/bin
make uninstall       # Remove from ~/.local/bin
make uninstall-system  # Remove from /usr/local/bin
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/description`
3. Commit with Conventional Commits: `git commit -m "feat: add X"`
4. Push and open a pull request

## License

This project is licensed under the GNU General Public License v3.0 — see [LICENSE](LICENSE) for details.
