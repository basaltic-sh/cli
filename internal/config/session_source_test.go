package config

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSessionRefreshProcess(t *testing.T) {
	if os.Getenv("BASALTIC_REFRESH_TEST_PROCESS") != "1" {
		return
	}
	source := &SessionTokenSource{Profile: "shared", Refresh: func(_ context.Context, _ *http.Client, endpoint, refresh string) (string, string, time.Time, error) {
		f, err := os.OpenFile(os.Getenv("BASALTIC_REFRESH_TEST_MARKER"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return "", "", time.Time{}, ErrSessionExpired
		}
		f.Close()
		time.Sleep(100 * time.Millisecond)
		return "renewed", "rotated", time.Now().Add(time.Hour), nil
	}}
	if _, err := source.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConcurrentSessionRefreshProcesses(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.yaml")
	t.Setenv(EnvCredentialsFile, path)
	if err := StoreSession("shared", "old", "initial", "https://example.test/token", time.Now().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := StoreSession("other", "keep", "keep-refresh", "https://other.test/token", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := updateCredentials(context.Background(), path, func(f *credentialsFile) error {
		f.store("service", "key", "cached", time.Now().Add(time.Hour))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.Command(os.Args[0], "-test.run=^TestSessionRefreshProcess$")
			cmd.Env = append(os.Environ(), "BASALTIC_REFRESH_TEST_PROCESS=1", "BASALTIC_REFRESH_TEST_MARKER="+path+".refreshed")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Errorf("child refresh: %s (%v)", out, err)
			}
		}()
	}
	wg.Wait()
	f := loadCredentials(path)
	if f.Sessions["shared"].AccessToken != "renewed" || f.Sessions["shared"].RefreshToken != "rotated" {
		t.Fatal("renewed session lost")
	}
	if f.Sessions["other"].AccessToken != "keep" || f.Tokens["service"].AccessToken != "cached" {
		t.Fatal("unrelated credentials lost")
	}
}

func TestSessionInvalidationAndTransientFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.yaml")
	t.Setenv(EnvCredentialsFile, path)
	if err := StoreSession("user", "old", "refresh", "endpoint", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	stale := &SessionTokenSource{Profile: "user"}
	if _, err := stale.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := StoreSession("user", "new", "rotated", "endpoint", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	stale.Invalidate()
	if access, _, _, _, ok := LookupSession("user"); !ok || access != "new" {
		t.Fatal("stale failure erased new login")
	}
	current := &SessionTokenSource{Profile: "user"}
	if _, err := current.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	current.Invalidate()
	if _, _, _, _, ok := LookupSession("user"); ok {
		t.Fatal("genuinely revoked login retained")
	}
	transient := errors.New("temporary network failure")
	for _, failure := range []error{transient, ErrSessionExpired} {
		if err := StoreSession("user", "old", "refresh", "endpoint", time.Now().Add(-time.Hour)); err != nil {
			t.Fatal(err)
		}
		source := &SessionTokenSource{Profile: "user", Refresh: func(context.Context, *http.Client, string, string) (string, string, time.Time, error) {
			return "", "", time.Time{}, failure
		}}
		if _, err := source.Token(context.Background()); !errors.Is(err, failure) {
			t.Fatalf("wrong error %v", err)
		}
		_, _, _, _, ok := LookupSession("user")
		if ok != errors.Is(failure, transient) {
			t.Fatal("incorrect refresh failure retention")
		}
	}
}

func TestCredentialWriterProcess(t *testing.T) {
	profile := os.Getenv("BASALTIC_CREDENTIAL_WRITER")
	if profile == "" {
		return
	}
	for i := 0; i < 10; i++ {
		if err := StoreSession(profile, "access", "refresh", "endpoint", time.Now().Add(time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestConcurrentCredentialWritersPreserveProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.yaml")
	t.Setenv(EnvCredentialsFile, path)
	var wg sync.WaitGroup
	for _, profile := range []string{"one", "two", "three", "four", "five", "six"} {
		wg.Add(1)
		go func(profile string) {
			defer wg.Done()
			cmd := exec.Command(os.Args[0], "-test.run=^TestCredentialWriterProcess$")
			cmd.Env = append(os.Environ(), "BASALTIC_CREDENTIAL_WRITER="+profile)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Errorf("writer: %s (%v)", out, err)
			}
		}(profile)
	}
	wg.Wait()
	if len(loadCredentials(path).Sessions) != 6 {
		t.Fatal("concurrent writes lost profiles")
	}
}
