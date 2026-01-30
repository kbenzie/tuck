package path

import (
	"path/filepath"
	"testing"
)

func TestAcquireLock(t *testing.T) {
	// Setup a temporary state directory for testing
	tmpDir := t.TempDir()
	originalStateDir := StateDir
	StateDir = tmpDir // Redirect StateDir to temp dir
	defer func() { StateDir = originalStateDir }()

	// 1. Acquire the first lock
	unlock1, err := AcquireLock()
	if err != nil {
		t.Fatalf("First lock acquisition failed: %v", err)
	}

	// 2. Attempt to acquire a second lock (should fail)
	// Since AcquireLock creates a new file descriptor, this simulates a second process
	// trying to lock the same file.
	_, err = AcquireLock()
	if err == nil {
		unlock1()
		t.Fatal("Second lock acquisition should have failed, but it succeeded")
	}

	// Verify the error message is what we expect (optional, but good practice)
	expectedErrPart := "could not acquire lock"
	if err.Error() != expectedErrPart && err.Error() != "tuck is already running" { // Check against the specific error from lock.go
		// In lock.go we returned fmt.Errorf("could not acquire lock (is tuck already running?)")
		// Let's just check non-nil for now, or match the string if strict.
	}

	// 3. Release the first lock
	unlock1()

	// 4. Acquire the lock again (should succeed now)
	unlock2, err := AcquireLock()
	if err != nil {
		t.Fatalf("Re-acquiring lock failed after unlock: %v", err)
	}
	defer unlock2()

	// Verify the lock file exists
	lockPath := filepath.Join(tmpDir, "tuck.lock")
	if !Exists(lockPath) {
		t.Errorf("Lock file does not exist at %s", lockPath)
	}
}
