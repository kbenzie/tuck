package path

import (
	"fmt"
	"path/filepath"

	"github.com/gofrs/flock"
)

// AcquireLock attempts to acquire a file lock to ensure only one instance of
// tuck is running at a time. It returns a cleanup function that must be called
// to release the lock, or an error if the lock could not be acquired.
func AcquireLock() (func(), error) {
	lockFile := filepath.Join(StateDir, "tuck.lock")
	fileLock := flock.New(lockFile)

	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock on %s: %w", lockFile, err)
	}

	if !locked {
		return nil, fmt.Errorf("could not acquire lock (is tuck already running?)")
	}

	unlock := func() {
		fileLock.Unlock()
	}

	return unlock, nil
}
