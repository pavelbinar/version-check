# version-check

version-check is a simple tool for comparing version numbers. It's designed to help developers ensure version consistency across different components of their projects.

## Features

- Compare semantic version numbers
- Support for custom version formats
- Command-line interface for easy integration

## Installation

To install version-check, make sure you have Go installed on your system. Then run:

```bash
go install github.com/pavelbinar/version-check@latest
```

## Usage

To use version-check, run the following command:

```bash
version-check
```

This will check the versions of all the tools listed in the config.yaml file and print the results.

## Configuration

The config.yaml file is used to specify the versions of the tools you want to check. The file should be in the following format:

```yaml
tools:
  - name: "ToolName"
    command: "command to check version"
    expect: "expected version"
```

Replace ToolName, command, and expect with the appropriate values for the tool you want to check.

### Example Config

```yaml
tools:
  - name: "Node"
    command: "node --version"
    expect: "20.17.0"
  - name: "Go"
    command: "go version"
    expect: "1.23.1"
```
