package internal

import (
	"os"
	"path/filepath"
	"syscall"
)

// runLockPath is a fixed, single lock file for the whole tool. Cron only
// ever runs one instance of this binary on a given host at a time, so a
// single global lock is enough to stop overlapping runs from racing each
// other on approvals/merges - even across different repos/branches passed
// on the command line.
var runLockPath = filepath.Join(os.TempDir(), "github-autoapproval.lock")

type RunLock struct {
	file *os.File
}

// AcquireRunLock takes an exclusive, non-blocking lock for the duration of
// one run of this tool. If another instance already holds it, ok is false
// and the caller should skip this run rather than wait, since cron will
// simply retry on its next tick.
func AcquireRunLock() (lock *RunLock, ok bool) {
	f, err := os.OpenFile(runLockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, false
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, false
	}

	return &RunLock{file: f}, true
}

func (l *RunLock) Release() {
	syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	l.file.Close()
}
