# GRC - Gmail Rule Converter

Generate Atom XML feeds from YAML filter configurations.

## Installation

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
sudo cp bin/grc /usr/local/bin/
```

Or use `go install`:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Usage

```bash
grc [options] <input.yaml> [output]
```

### Arguments

| Argument    | Description                                                  |
|-------------|--------------------------------------------------------------|
| `input.yaml` | Path to the YAML configuration file (required)              |
| `output`     | Output path for the XML file (default: `<input>.xml`)       |

### Options

| Option       | Description                            |
|--------------|----------------------------------------|
| `--force`    | Overwrite existing output file         |
| `--verbose`  | Enable detailed logging                |
| `--version`  | Show version information               |
| `--help`     | Show this help message                 |

## YAML Configuration

### Structure

```yaml
defaults:
  author:
    name: Default Author
    email: author@example.com
  category:
    term: default
    label: Default Category

filters:
  - name: Filter Name
    expression: subject:keyword
    category:
      term: category-slug
      label: Category Label
    author:
      name: Author Name
      email: author@example.com
    properties:
      - name: property-name
        value: property-value
```

### Fields

| Field        | Type   | Required | Description                          |
|--------------|--------|----------|--------------------------------------|
| `name`       | string | yes      | Filter name (becomes feed entry title) |
| `expression` | string | yes      | Search expression                    |
| `category`   | object | no       | Atom category (term, scheme, label)  |
| `author`     | object | no       | Atom author (name, email, uri)       |
| `properties` | array  | no       | Key-value pairs added to summary     |

### Defaults

The `defaults` block sets fallback values for `author` and `category` applied to any filter that does not define them explicitly.

## Examples

### Minimal

```yaml
filters:
  - name: Newsletter
    expression: subject:newsletter
```

### Full

See `cfg/example.yaml` for a complete example with newsletters, support tickets, internal notices, client communications, financial documents, and spam filters.

## Building from Source

```bash
make build        # Compile binary to bin/
make test         # Run all tests
make quality      # Format, vet, and test
make clean        # Remove build artifacts
```

Use `make build` or build directly:

```bash
go build -o bin/grc ./grc/cmd/grc
```

Embed version info:

```bash
go build -ldflags="-X github.com/carlosrabelo/grc/grc/internal/app.Version=$(git describe --tags --always --dirty) -X github.com/carlosrabelo/grc/grc/internal/app.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/grc ./grc/cmd/grc
```
