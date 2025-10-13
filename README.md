# GRC - Gmail Rules Creator

A command-line tool to generate XML files for configuring Gmail filters from YAML configuration files. Simplifies the management and customization of Gmail mail filters with validation and default value support.

## Features

- YAML Configuration: Define filters using clean, readable YAML syntax
- Comprehensive Criteria: Support for all Gmail filter criteria (from, to, subject, query, attachments, etc.)
- Rich Actions: Full range of Gmail actions (archive, mark as read, star, forward, trash, labels, smart labels)
- Default Values: Automatically applies default boolean action values when omitted
- Validation: Ensures author details, at least one filter, and that each filter has criteria and actions
- XML Generation: Outputs properly formatted XML compatible with Gmail's filter import
- Verbose Logging: Optional detailed logging for debugging and monitoring

## Project Structure
```
grc/
├── core/             # Core Go module and source code
│   ├── cmd/
│   │   └── grc/      # Main CLI application entry point
│   ├── internal/
│   │   ├── app/      # Application logic and CLI handling
│   │   └── rules/    # Core filtering logic and XML generation
│   ├── go.mod        # Go module definition
│   ├── go.sum        # Go module checksums
│   └── Makefile      # Core build automation
├── resources/
│   └── example.yaml  # Example YAML configuration
├── scripts/          # Installation and utility scripts
├── bin/              # Generated binaries (with .gitkeep)
├── Makefile          # Root build automation (delegates to core/)
├── README.md         # Project documentation
└── README-PT.md      # Documentação em português
```

## Filter Configuration

### Criteria (Conditions)
Each filter must include at least one of these criteria:
- `from` - Match sender email address
- `to` - Match recipient email address
- `subject` - Match email subject line
- `hasTheWord` - Match emails containing specific words
- `doesNotHaveTheWord` - Match emails NOT containing specific words
- `list` - Match mailing list emails
- `query` - Use Gmail search query syntax
- `hasAttachment` - Match emails with/without attachments

### Actions
Each filter must include at least one action:
- `label` - Apply a label to matching emails
- `smartLabel` - Apply Gmail smart labels (Important, Spam, etc.)
- `forwardTo` - Forward matching emails to another address
- Boolean actions:
  - `shouldArchive` - Skip inbox (archive)
  - `shouldMarkAsRead` - Mark as read automatically
  - `shouldStar` - Add star to matching emails
  - `shouldNeverSpam` - Never send to spam
  - `shouldAlwaysMarkAsImportant` - Always mark as important
  - `shouldNeverMarkAsImportant` - Never mark as important
  - `shouldTrash` - Delete matching emails

Boolean actions inherit defaults from the `default` section when not specified.

## Prerequisites
- Go 1.22 or later

## Installation

### Option 1: Build from Source
```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

### Option 2: Install Locally
```bash
make install
```
This installs the `grc` binary to your local bin directory (`$HOME/.local/bin` for users, `/usr/local/bin` for root) after building it automatically.



## Usage

### Basic Usage
```bash
grc [options] <yaml_file>
```

### Options
- `-output <file>` - Specify output XML file path (default: same as input with .xml extension)
- `-verbose` - Enable detailed logging output

### Example YAML Configuration
```yaml
author:
  name: "John Doe"
  email: "john.doe@corp.com"

default:
  shouldArchive: true
  shouldMarkAsRead: false
  shouldStar: false
  shouldNeverSpam: true
  shouldAlwaysMarkAsImportant: false
  shouldNeverMarkAsImportant: false
  shouldTrash: false

filters:
  - from: "info@newsletter.shopee.com.br"
    label: "@SaneLater"
  - to: "support@example.com"
    subject: "[Ticket]"
    hasAttachment: true
    label: "@Support"
    shouldMarkAsRead: true
    shouldStar: true
  - query: "list:announcements.example.com"
    label: "@Announcements"
    forwardTo: "archive@example.com"
    shouldArchive: false
    shouldAlwaysMarkAsImportant: true
```

### Examples
```bash
# Generate XML from YAML config
grc resources/example.yaml

# Specify custom output file
grc -output my-filters.xml resources/example.yaml

# Enable verbose logging
grc -verbose resources/example.yaml
```

## Development

### Available Make Targets
```bash
# Root level (delegates to core/)
make help          # Show available targets (default)
make build         # Build the binary
make test          # Run test suite
make run           # Run the application
make clean         # Remove build artifacts
make install       # Install binary locally
make uninstall     # Remove installed binary

# Direct core/ access
cd core && make help  # Show core-specific targets
cd core && make lint   # Run linter (golangci-lint)
```

### Working with the Core Module
The project uses a core/ directory structure with its own Go module:
- Module: `github.com/carlosrabelo/grc/core`
- All Go development happens within the core/ directory
- The root Makefile delegates all Go-related commands to core/Makefile
- Running `make` without arguments shows help by default

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
