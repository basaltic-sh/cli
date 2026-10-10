package config

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// Every read-modify-write shares this OS lock, including token rotation. Keep
// the lock file in place: deleting it would let callers lock different inodes.
// The OS releases ownership if a command exits or is killed.
func lockCredentials(ctx context.Context, path string) (func(), error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	for {
		if err := ctx.Err(); err != nil {
			f.Close()
			return nil, err
		}
		locked, err := tryCredentialLock(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		if locked {
			return func() { f.Close() }, nil
		}
		select {
		case <-ctx.Done():
			f.Close()
			return nil, ctx.Err()
		case <-time.After(25 * time.Millisecond):
		}
	}
}

func updateCredentials(ctx context.Context, path string, update func(*credentialsFile) error) error {
	unlock, err := lockCredentials(ctx, path)
	if err != nil {
		return err
	}
	defer unlock()
	file := loadCredentials(path)
	if err := update(file); err != nil {
		return err
	}
	return saveCredentials(path, file)
}
