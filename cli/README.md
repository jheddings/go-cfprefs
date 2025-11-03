# cfprefs - command-line interface to CFPreferences

This CLI serves as a demonstration of the `go-cfprefs` module. It loosely follows the `defaults` command on macOS and provides support for JSON Pointer paths to work with complex nested preference structures.

## Installation

Build the CLI from source:

```bash
make build
```

The binary will be created in the `dist/` directory.

## Commands

### `read` - Read preference values

Read preference values from CFPreferences, with support for JSON Pointer paths.

#### Basic Usage

```bash
# Read all keys for an application
cfprefs read com.example.app

# Read a specific key
cfprefs read com.example.app username

# Read a nested field using JSON Pointer path
cfprefs read com.example.app config/server/port

# Read an array element
cfprefs read com.example.app items/0
```

### `write` - Write preference values

Write preference values to CFPreferences, with support for JSON Pointer paths.

#### Basic Usage

```bash
# Write a string value
cfprefs write com.example.app username "john_doe"

# Write with type specification
cfprefs write com.example.app maxConnections "10" --int
cfprefs write com.example.app isEnabled "true" --bool
cfprefs write com.example.app lastLogin "2024-01-15T10:30:00Z" --date

# Write a nested value using JSON Pointer path
cfprefs write com.example.app config/server/port 8080 --int

# Replace an array element
cfprefs write com.example.app items/0 "updated item"
```

#### Advanced Operators

The `write` commands supports additional operators for working with data structures.

***Array Append***

To append an element to the end of an array, use the `~]` operator.

```bash
cfprefs write com.example.app items/~] "last item"
```

***Array Prepend***

To insert an element at the beginning of an array, use the `~[` operator.

```bash
cfprefs write com.example.app items/~[ "first item"
```

#### Type Flags

- `--string` (default): Parse value as string
- `--int`: Parse value as integer
- `--float`: Parse value as float
- `--bool`: Parse value as boolean
- `--date`: Parse value as date (ISO 8601 format)

### `delete` - Delete preference keys

Delete preference keys from CFPreferences, with support for JSON Pointer paths.

#### Basic Usage

```bash
# Delete a specific key
cfprefs delete com.example.app username

# Delete a nested field using JSON Pointer path
cfprefs delete com.example.app config/server/port

# Delete an array element
cfprefs delete com.example.app items/0
```

## JSON Pointer Path Syntax

Preference keys can be specified as simple names or as [JSON Pointer](https://datatracker.ietf.org/doc/html/rfc6901) paths to access nested values. Use `/` to separate object keys and array indices.

For more details, see the [JSON Pointer RFC](https://datatracker.ietf.org/doc/html/rfc6901).

## Global Flags

- `-v, --verbose`: Increase verbosity in logging
- `-q, --quiet`: Only log errors and warnings
- `-y, --yes`: Assume 'yes' for confirmation prompts
