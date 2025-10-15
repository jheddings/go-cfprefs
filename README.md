# go-cfprefs

Go module wrapper for the `CFPreferences` API's in macOS.

## Features

- Read and write macOS preferences (plist files) using native Go types
- Automatic type conversion between Go types and CoreFoundation types
- Support for all common data types: strings, numbers, booleans, dates, arrays, dictionaries, and binary data

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

// The value is returned as the appropriate Go type
if showPathbar, ok := value.(bool); ok {
    fmt.Printf("Show Pathbar: %v\n", showPathbar)
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
```

### Deleting Preferences

```go
err := cfprefs.Delete("com.example.app", "username")
```

### Checking if a Key Exists

```go
exists, err := cfprefs.Exists("com.example.app", "username")
```

## Using Keypaths

TODO - document usage

## Command-Line Interface

There is a basic CLI that acts as a demonstration of this module, as well as used for testing.

## License

See [LICENSE](LICENSE) for details.
