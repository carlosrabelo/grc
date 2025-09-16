# GRC - GMail Rules Creator

This project provides a tool to generate XML files for configuring Gmail filters based on a YAML input file. It allows for easy management and customization of mail filters while adhering to a specified structure.

## Features
- Accepts a YAML input file containing author details, default values, and filter configurations.
- Supports a broad set of Gmail filter criteria (from, to, subject, query, attachments, etc.) and actions (archive, mark as read, star, forward, trash, labels and more).
- Automatically applies default values to boolean actions when optional parameters are omitted.
- Outputs a formatted XML file compatible with Gmail's filter configuration.

## Project Structure
```
project/
├── resources/
│   └── example.yaml  # Example YAML file for testing
├── main.go           # Main source code file
├── Makefile          # Build and execution commands
├── .gitignore        # Ignored files and directories
└── README.md         # Project documentation
```

### Available filter options

Each filter entry may include any combination of the following criteria:
- `from`, `to`, `subject`, `hasTheWord`, `doesNotHaveTheWord`, `list`, `query`, `hasAttachment`

Actions that can be applied:
- `label`, `smartLabel`, `forwardTo`
- Boolean flags: `shouldArchive`, `shouldMarkAsRead`, `shouldStar`, `shouldNeverSpam`, `shouldAlwaysMarkAsImportant`, `shouldNeverMarkAsImportant`, `shouldTrash`

All boolean flags inherit their defaults from the `default` section when omitted. Every filter must declare at least one condition and one action.

## Prerequisites
- Go (Golang) 1.16 or later

## Getting Started

### 1. Clone the repository
```bash
git clone <repository-url>
cd project
```

### 2. Prepare your YAML input file
Ensure that a YAML file with the required structure exists in the `resources/` directory. You can use the provided `example.yaml` as a reference:

```yaml
author:
  name: "John Doe"
  email: "john.doe@gmail.com"

default:
  shouldArchive: true
  shouldMarkAsRead: false
  shouldStar: false
  shouldNeverSpam: true
  shouldAlwaysMarkAsImportant: false
  shouldNeverMarkAsImportant: false
  shouldTrash: false

filters:
  - from: "example1@example.com"
    label: "Work"
  - to: "support@example.com"
    subject: "[Ticket]"
    hasAttachment: true
    label: "Support"
    shouldMarkAsRead: true
    shouldStar: true
  - query: "list:announcements.example.com"
    label: "Announcements"
    forwardTo: "archive@example.com"
    shouldArchive: false
    shouldAlwaysMarkAsImportant: true
```

### 3. Build the project
Run the following command to compile the binary:
```bash
make build
```

### 4. Generate XML
Use the provided `generate-sample` Makefile target to create the XML file for the bundled example:
```bash
make generate-sample
```
This will generate an XML file in the `resources/` directory based on the YAML input.

### 5. Clean up
To remove generated files, use:
```bash
make clean
```

## Usage
Run the compiled binary directly with a specified YAML file as input:
```bash
./grc rules.yaml
```
The output XML file will be saved in the same directory as the input YAML file, with the `.xml` extension.

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
