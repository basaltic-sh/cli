package selfupdate

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// The update check is a background courtesy, not a feature of any command. It
// must never delay what the user asked for, never write to stdout, and never
// speak at all when the output is being piped somewhere.
//
// It notifies; it does not update. A CLI that replaces its own binary without
// being asked is a CLI that changes under a script's feet between two runs.

// EnvDisable turns the check off. NO_UPDATE_NOTIFIER is honoured too, since
// it is the conventional spelling and someone who sets it means all of them.
const EnvDisable = "BASALTIC_NO_UPDATE_CHECK"

// checkInterval is how often the latest version is looked up.
const checkInterval = 24 * time.Hour

// waitBudget is how long the exit path will wait for an in-flight check.
// Small on purpose: the check usually finished while the command was running,
// and if it did not, the answer keeps until next time.
const waitBudget = 700 * time.Millisecond

type state struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version"`
}

// Checker looks up the newest release in the background.
type Checker struct {
	Current string
	Client  *http.Client

	once   sync.Once
	done   chan struct{}
	cached string
}

// NewChecker returns a checker, or nil when checking is switched off.
//
// interactive should report whether the CLI is talking to a terminal: a notice
// printed into a pipeline is noise at best, and at worst lands in the middle
// of something being parsed.
func NewChecker(current string, hc *http.Client, interactive bool) *Checker {
	if !interactive || disabled() {
		return nil
	}
	// A development build has no meaningful version to compare, so there is
	// nothing to tell anyone.
	if canonical(current) == "" {
		return nil
	}
	return &Checker{Current: current, Client: hc, done: make(chan struct{})}
}

func disabled() bool {
	for _, k := range []string{EnvDisable, "NO_UPDATE_NOTIFIER"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" && v != "0" && v != "false" {
			return true
		}
	}
	// A CI run is not a person who can act on a notice.
	return os.Getenv("CI") != ""
}

// Start reads the cached answer and, when it is stale, refreshes it in the
// background. It returns immediately.
func (c *Checker) Start(ctx context.Context) {
	if c == nil {
		return
	}
	c.once.Do(func() {
		st := loadState()
		c.cached = st.LatestVersion

		if time.Since(st.CheckedAt) < checkInterval {
			close(c.done)
			return
		}
		go func() {
			defer close(c.done)
			// Detached from the command's context: cancelling the command
			// should not be reported as a failed update check, and this
			// writes nothing the command depends on.
			lookup, cancel := context.WithTimeout(context.Background(), apiTimeout)
			defer cancel()

			rel, err := Latest(lookup, c.Client)
			if err != nil {
				// Silence is correct here. A background courtesy that
				// complains about the network is worse than one that says
				// nothing, and the user did not ask for this.
				return
			}
			saveState(state{CheckedAt: time.Now(), LatestVersion: rel.Version()})
			c.cached = rel.Version()
		}()
	})
}

// Notice returns the message to show at exit, or "".
func (c *Checker) Notice() string {
	if c == nil {
		return ""
	}
	select {
	case <-c.done:
	case <-time.After(waitBudget):
		// Still running. Whatever it finds is cached for next time.
		return ""
	}
	if c.cached == "" || !IsNewer(c.Current, c.cached) {
		return ""
	}
	return "\nA newer basaltic is available: " + c.Current + " → " + c.cached +
		"\nRun `basaltic upgrade` to install it, or set " + EnvDisable + "=1 to stop checking.\n"
}

func statePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "basaltic", "update-check.json")
}

// loadState never fails: an unreadable cache is an empty one, which costs a
// lookup rather than breaking anything.
func loadState() state {
	p := statePath()
	if p == "" {
		return state{}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return state{}
	}
	var st state
	if err := json.Unmarshal(data, &st); err != nil {
		return state{}
	}
	return st
}

func saveState(st state) {
	p := statePath()
	if p == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return
	}
	data, err := json.Marshal(st)
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o600)
}
