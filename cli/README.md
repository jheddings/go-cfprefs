# cfprefs - command-line interface to CFPreferences

This CLI serves as a demonstration of the `go-cfprefs` module. It loosely follows the `defaults` command on macOS and provides powerful JSONPath query capabilities for working with complex nested preference structures.

## Installation

Build the CLI from source:

```bash
make build
```

The binary will be created in the `dist/` directory.

## Commands

### `read` - Read preference values

Read preference values from CFPreferences, with support for keypaths and JSONPath queries.

#### Basic Usage

```bash
# Read all keys for an application
cfprefs read com.example.app

# Read a specific key
cfprefs read com.example.app username

# Read a nested value using keypath
cfprefs read com.example.app settings/display/brightness
```

#### JSONPath Query Usage

```bash
# Read a nested field using JSONPath
cfprefs read com.example.app userData --query '$.user.name'

# Read all items from an array
cfprefs read com.example.app data --query '$.items[*]'

# Read filtered array items
cfprefs read com.example.app data --query '$.items[?(@.active == true)]'

# Read a specific array element
cfprefs read com.example.app data --query '$.items[0]'
```

### `write` - Write preference values

Write preference values to CFPreferences, with support for keypaths and JSONPath queries.

#### Basic Usage

```bash
# Write a string value
cfprefs write com.example.app username "john_doe"

# Write a nested value using keypath
cfprefs write com.example.app settings/display/brightness "75"

# Write with type specification
cfprefs write com.example.app maxConnections "10" --int
cfprefs write com.example.app isEnabled "true" --bool
cfprefs write com.example.app lastLogin "2024-01-15T10:30:00Z" --date
```

#### JSONPath Query Usage

```bash
# Set a nested field using JSONPath
cfprefs write com.example.app userData --query '$.user.name' "John Doe"

# Set an array element
cfprefs write com.example.app data --query '$.items[0]' "new item"

# Append to an array
cfprefs write com.example.app data --query '$.items[]' "appended item"

# Set a deeply nested value
cfprefs write com.example.app config --query '$.database.host' "localhost"
```

#### Type Flags

- `--string` (default): Parse value as string
- `--int`: Parse value as integer
- `--float`: Parse value as float
- `--bool`: Parse value as boolean
- `--date`: Parse value as date (ISO 8601 format)

### `delete` - Delete preference keys

Delete preference keys from CFPreferences, with support for keypaths and JSONPath queries.

#### Basic Usage

```bash
# Delete a specific key
cfprefs delete com.example.app username

# Delete a nested key using keypath
cfprefs delete com.example.app settings/display/brightness
```

#### JSONPath Query Usage

```bash
# Delete a nested field using JSONPath
cfprefs delete com.example.app userData --query '$.user.name'

# Delete an array element
cfprefs delete com.example.app data --query '$.items[0]'

# Delete a deeply nested field
cfprefs delete com.example.app config --query '$.database.host'
```

## JSONPath Query Syntax

The `--query` flag accepts JSONPath expressions for precise data manipulation:

### Basic Syntax

- `$` - Root object
- `$.field` - Access object field
- `$.array[0]` - Access array element by index
- `$.array[*]` - Access all array elements
- `$.field.subfield` - Nested field access

### Advanced Queries

- `$.items[?(@.active == true)]` - Filter array items by condition
- `$.items[?(@.count > 5)]` - Filter by numeric comparison
- `$.items[?(@.name =~ /^prefix/)]` - Filter by regex pattern

### Examples

```bash
# Get all active items
cfprefs read com.example.app data --query '$.items[?(@.active == true)]'

# Get the first item's name
cfprefs read com.example.app data --query '$.items[0].name'

# Set a value in a filtered location
cfprefs write com.example.app data --query '$.items[?(@.id == "item1")].status' "updated"

# Delete a specific array element
cfprefs delete com.example.app data --query '$.items[?(@.id == "old_item")]'
```

## Global Flags

- `-v, --verbose`: Increase verbosity in logging
- `-q, --quiet`: Only log errors and warnings
- `-y, --yes`: Assume 'yes' for confirmation prompts

## Examples

### Working with Application Settings

```bash
# Read application theme
cfprefs read com.example.app settings --query '$.theme'

# Set user preferences
cfprefs write com.example.app userPrefs --query '$.notifications.email' --bool true
cfprefs write com.example.app userPrefs --query '$.notifications.push' --bool false

# Update nested configuration
cfprefs write com.example.app config --query '$.database.connectionTimeout' --int 30
cfprefs write com.example.app config --query '$.database.retries' --int 3
```

### Managing Complex Data Structures

```bash
# Read all user profiles
cfprefs read com.example.app profiles --query '$.users[*]'

# Add a new user to the array
cfprefs write com.example.app profiles --query '$.users[]' '{"id": "user123", "name": "John Doe"}'

# Update a specific user's settings
cfprefs write com.example.app profiles --query '$.users[?(@.id == "user123")].settings.theme' "dark"

# Delete an inactive user
cfprefs delete com.example.app profiles --query '$.users[?(@.active == false)]'
```
