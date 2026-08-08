package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitializationMarkerLifecycle(t *testing.T) {
	markerPath := filepath.Join(t.TempDir(), ".seed_executed")
	if initializationComplete(markerPath, seedDataVersion) {
		t.Fatal("missing marker must not be treated as complete")
	}

	if err := markInitializationComplete(markerPath, seedDataVersion); err != nil {
		t.Fatalf("mark initialization complete: %v", err)
	}
	if !initializationComplete(markerPath, seedDataVersion) {
		t.Fatal("current marker version must be treated as complete")
	}
	if initializationComplete(markerPath, "3") {
		t.Fatal("older marker version must not skip a newer initialization")
	}
}

func TestLegacyInitializationMarkerOnlySupportsVersionOne(t *testing.T) {
	markerPath := filepath.Join(t.TempDir(), ".seed_executed")
	if err := os.WriteFile(markerPath, []byte(legacyMarkerValue), 0644); err != nil {
		t.Fatalf("write legacy marker: %v", err)
	}
	if !initializationComplete(markerPath, "1") {
		t.Fatal("legacy executed marker should remain compatible with version 1")
	}
	if initializationComplete(markerPath, seedDataVersion) {
		t.Fatal("legacy marker must not skip newer seed data")
	}
}
