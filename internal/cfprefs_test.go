package internal

import (
	"testing"
	"time"
)

func TestBasicGetSet(t *testing.T) {
	dateTime := time.Now().Format(time.RFC3339)
	t.Log(dateTime)

	err := Set("com.jheddings.cfprefs.testing", "basic-test", dateTime)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "basic-test")

	if err != nil {
		t.Fatal(err)
	}

	if value != dateTime {
		t.Fatal("value does not match")
	}

	err = Delete("com.jheddings.cfprefs.testing", "basic-test")

	if err != nil {
		t.Fatal(err)
	}

	exists, err := Exists("com.jheddings.cfprefs.testing", "basic-test")

	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatal("key should not exist")
	}
}
