package cfprefs

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestConcurrentRead tests concurrent read operations
func TestConcurrentRead(t *testing.T) {
	// Setup test data
	testData := map[string]any{
		"string": "test value",
		"number": int64(42),
		"bool":   true,
		"array":  []any{"a", "b", "c"},
		"object": map[string]any{
			"nested": "value",
			"count":  int64(10),
		},
	}

	key := "test-concurrent-read"
	cleanup := setupTest(t, testAppID, key, testData)
	defer cleanup()

	t.Run("concurrent reads of same key", func(t *testing.T) {
		const goroutines = 50
		const iterations = 100

		var wg sync.WaitGroup
		wg.Add(goroutines)

		errors := make(chan error, goroutines*iterations)

		for range goroutines {
			go func() {
				defer wg.Done()
				for range iterations {
					value, err := Get(testAppID, key)
					if err != nil {
						errors <- fmt.Errorf("read error: %v", err)
						continue
					}

					// Verify we got the expected data
					if mapValue, ok := value.(map[string]any); ok {
						if mapValue["string"] != "test value" {
							errors <- fmt.Errorf("unexpected value: %v", mapValue["string"])
						}
					} else {
						errors <- fmt.Errorf("value is not a map: %T", value)
					}
				}
			}()
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			t.Error(err)
		}
	})

	t.Run("concurrent reads of different paths", func(t *testing.T) {
		paths := []string{
			key + "/string",
			key + "/number",
			key + "/bool",
			key + "/array/0",
			key + "/object/nested",
		}

		const goroutinesPerPath = 10
		var wg sync.WaitGroup
		wg.Add(len(paths) * goroutinesPerPath)

		errors := make(chan error, len(paths)*goroutinesPerPath*100)

		for _, path := range paths {
			for i := 0; i < goroutinesPerPath; i++ {
				go func(p string) {
					defer wg.Done()
					for j := 0; j < 100; j++ {
						_, err := Get(testAppID, p)
						if err != nil {
							errors <- fmt.Errorf("read error on %s: %v", p, err)
						}
					}
				}(path)
			}
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			t.Error(err)
		}
	})
}

// TestConcurrentWrite tests concurrent write operations
func TestConcurrentWrite(t *testing.T) {
	t.Run("concurrent writes to different keys", func(t *testing.T) {
		const goroutines = 20
		const keysPerGoroutine = 10

		var wg sync.WaitGroup
		wg.Add(goroutines)

		errors := make(chan error, goroutines*keysPerGoroutine)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < keysPerGoroutine; j++ {
					key := fmt.Sprintf("concurrent-write-%d-%d", id, j)
					value := fmt.Sprintf("value-%d-%d", id, j)

					err := Set(testAppID, key, value)
					if err != nil {
						errors <- fmt.Errorf("write error: %v", err)
						continue
					}

					// Clean up
					defer Delete(testAppID, key)

					// Verify the write
					got, err := GetStr(testAppID, key)
					if err != nil {
						errors <- fmt.Errorf("read after write error: %v", err)
					} else if got != value {
						errors <- fmt.Errorf("expected %s, got %s", value, got)
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			t.Error(err)
		}
	})

	t.Run("concurrent updates to same key", func(t *testing.T) {
		key := "concurrent-update-test"
		const goroutines = 10
		const iterations = 50

		// Initialize the key
		cleanup := setupTest(t, testAppID, key, int64(0))
		defer cleanup()

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					// Each goroutine writes its ID
					_ = Set(testAppID, key, int64(id))

					// Small random delay to increase contention
					time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))
				}
			}(i)
		}

		wg.Wait()

		// The final value should be from one of the goroutines (0-9)
		finalValue, err := GetInt(testAppID, key)
		testutil.AssertNoError(t, err, "get final value")

		if finalValue < 0 || finalValue >= goroutines {
			t.Errorf("unexpected final value: %d", finalValue)
		}
	})
}

// TestConcurrentMixed tests mixed read/write operations
func TestConcurrentMixed(t *testing.T) {
	t.Run("reads and writes on different keys", func(t *testing.T) {
		// Setup some initial data to read
		readKeys := make([]string, 10)
		for i := range readKeys {
			readKeys[i] = fmt.Sprintf("mixed-read-%d", i)
			cleanup := setupTest(t, testAppID, readKeys[i], fmt.Sprintf("value-%d", i))
			defer cleanup()
		}

		const readers = 20
		const writers = 10
		const operations = 50

		var wg sync.WaitGroup
		wg.Add(readers + writers)

		errors := make(chan error, (readers+writers)*operations)

		// Start readers
		for i := 0; i < readers; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < operations; j++ {
					key := readKeys[rand.Intn(len(readKeys))]
					_, err := Get(testAppID, key)
					if err != nil {
						errors <- fmt.Errorf("read error: %v", err)
					}
				}
			}()
		}

		// Start writers
		for i := 0; i < writers; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operations; j++ {
					key := fmt.Sprintf("mixed-write-%d-%d", id, j)
					err := Set(testAppID, key, fmt.Sprintf("value-%d-%d", id, j))
					if err != nil {
						errors <- fmt.Errorf("write error: %v", err)
					}
					// Clean up
					defer Delete(testAppID, key)
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			t.Error(err)
		}
	})
}

// TestConcurrentDelete tests concurrent delete operations
func TestConcurrentDelete(t *testing.T) {
	t.Run("concurrent deletes of different keys", func(t *testing.T) {
		const numKeys = 100
		const goroutines = 10

		// Setup keys to delete
		baseKey := "concurrent-delete"
		for i := 0; i < numKeys; i++ {
			key := fmt.Sprintf("%s-%d", baseKey, i)
			err := Set(testAppID, key, fmt.Sprintf("value-%d", i))
			testutil.AssertNoError(t, err, "setup key")
		}
		defer func() {
			// Clean up any remaining keys
			for i := 0; i < numKeys; i++ {
				Delete(testAppID, fmt.Sprintf("%s-%d", baseKey, i))
			}
		}()

		// Divide keys among goroutines
		keysPerGoroutine := numKeys / goroutines

		var wg sync.WaitGroup
		wg.Add(goroutines)

		errors := make(chan error, numKeys)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				start := id * keysPerGoroutine
				end := start + keysPerGoroutine
				if id == goroutines-1 {
					end = numKeys
				}

				for j := start; j < end; j++ {
					key := fmt.Sprintf("%s-%d", baseKey, j)
					err := Delete(testAppID, key)
					if err != nil {
						errors <- fmt.Errorf("delete error on %s: %v", key, err)
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			t.Error(err)
		}

		// Verify all keys were deleted
		for i := 0; i < numKeys; i++ {
			key := fmt.Sprintf("%s-%d", baseKey, i)
			exists, err := Exists(testAppID, key)
			testutil.AssertNoError(t, err, "check exists")
			if exists {
				t.Errorf("key %s should have been deleted", key)
			}
		}
	})
}
