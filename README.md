# go-cfprefs

Go module wrapper for the `CFPreferences` API's in macOS.

## Features

- Read and write macOS preferences using native Go types
- Automatic type conversion between Go types and CoreFoundation types
- Support for all common data types: strings, numbers, booleans, dates, arrays, dictionaries, and binary data
    - Use JSON Pointer paths to access nested structures

## Installation

```bash
go get github.com/jheddings/go-cfprefs
```

## Usage

### Reading Preferences

```go
import "github.com/jheddings/go-cfprefs"

// Read a preference value
value, err := cfprefs.Get("com.apple.finder", "ShowPathbar")
if err != nil {
    log.Fatal(err)
}

// Read a nested value using JSON Pointer path
value, err = cfprefs.Get("com.example.app", "config/server/port")
if err != nil {
    log.Fatal(err)
}
```

### Writing Preferences

The `Set` function accepts any native Go type and automatically converts it to the appropriate CoreFoundation type:

```go
import "github.com/jheddings/go-cfprefs"

// Write a string
err := cfprefs.Set("com.example.app", "username", "john_doe")

// Write a number
err = cfprefs.Set("com.example.app", "count", 42)

// Write a boolean
err = cfprefs.Set("com.example.app", "enabled", true)

// Write a float
err = cfprefs.Set("com.example.app", "pi", 3.14159)

// Write a date
err = cfprefs.Set("com.example.app", "lastAccess", time.Now())

// Write an array
err = cfprefs.Set("com.example.app", "items", []any{"apple", "banana", "cherry"})

// Write a dictionary
err = cfprefs.Set("com.example.app", "config", map[string]any{
    "theme": "dark",
    "fontSize": 14,
    "autoSave": true,
})

// Write a nested value using JSON Pointer path
err = cfprefs.Set("com.example.app", "config/server/port", 8080)
```

### Deleting Preferences

```go
// Delete a top-level key
err := cfprefs.Delete("com.example.app", "username")

// Delete a nested value using JSON Pointer path
err = cfprefs.Delete("com.example.app", "config/server/port")
```

### Checking if a Key Exists

```go
// Check for a top-level key
exists, err := cfprefs.Exists("com.example.app", "username")

// Check for a nested value using JSON Pointer path
exists, err = cfprefs.Exists("com.example.app", "config/server/port")
```

## Path Syntax

Preference keys can be specified as simple names or as [JSON Pointer](https://datatracker.ietf.org/doc/html/rfc6901) paths to access nested values. Use `/` to separate object keys and array indices.

### Examples

- `settings` — top-level key
- `user/name` — nested object field
- `items/0` — first element of an array
- `config/database/host` — deeply nested field

For more details, see the [JSON Pointer RFC](https://datatracker.ietf.org/doc/html/rfc6901).

## Command-Line Interface

There is a [basic CLI](cli/README.md) that acts as a demonstration of this module, as well as used for testing.

## License

This project is licensed under the terms of the MIT license. See [LICENSE](LICENSE) for details.
