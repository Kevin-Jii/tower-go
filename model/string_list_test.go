package model

import (
	"reflect"
	"testing"
)

func TestStringListDatabaseRoundTrip(t *testing.T) {
	want := StringList{"https://example.com/one.jpg", "https://example.com/two.jpg"}
	value, err := want.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}

	var got StringList
	if err := got.Scan(value); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}

func TestStringListScanNullReturnsEmptyList(t *testing.T) {
	var got StringList
	if err := got.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error = %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Fatalf("Scan(nil) = %#v, want non-nil empty list", got)
	}
}
