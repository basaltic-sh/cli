// Package config reads and writes the CLI's two files.
//
// ~/.config/basaltic/config.yaml holds named profiles. It is the file people
// copy between machines and check into dotfile repositories:
//
//	default_profile: production
//	output: text
//	profiles:
//	  production:
//	    region: sa-saopaulo-1
//	    api_key: ACCESS_KEY_ID:SECRET_ACCESS_KEY
//	    account_id: acme
//
// ~/.config/basaltic/credentials.yaml holds cached access tokens, and is
// deliberately a separate file for that reason — a token must not travel with
// a config someone copies around. See credentials.go.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// EnvConfigFile overrides the config file location.
const EnvConfigFile = "BASALTIC_CONFIG_FILE"

// File is the on-disk shape of config.yaml.
type File struct {
	DefaultProfile string             `yaml:"default_profile,omitempty"`
	Output         string             `yaml:"output,omitempty"`
	Profiles       map[string]Profile `yaml:"profiles,omitempty"`
}

// Profile is one named environment.
type Profile struct {
	// Region is the region regional services are addressed in.
	Region string `yaml:"region,omitempty"`
	// APIKey is the service account key pair, written "ACCESS_KEY_ID:SECRET".
	// One field because it is one credential; people copy it around as one
	// string and splitting it in the file only invites half of it being
	// pasted.
	APIKey string `yaml:"api_key,omitempty"`
	// AccountID selects the account requests act on.
	AccountID string `yaml:"account_id,omitempty"`

	// Domain builds every endpoint from another domain, for reaching a
	// deployment that is not production. Undocumented for end users.
	Domain string `yaml:"domain,omitempty"`
	// Endpoints overrides individual services, keyed by service name.
	Endpoints map[string]string `yaml:"endpoints,omitempty"`
	// Insecure disables TLS verification for every host this profile talks
	// to. For development rigs with self-signed certificates.
	Insecure bool `yaml:"insecure,omitempty"`
}

// Credentials splits the profile's key pair.
func (p Profile) Credentials() (accessKeyID, secretAccessKey string, err error) {
	id, secret, found := strings.Cut(strings.TrimSpace(p.APIKey), ":")
	if !found || id == "" || secret == "" {
		return "", "", errors.New("api_key must be written as ACCESS_KEY_ID:SECRET_ACCESS_KEY")
	}
	return id, secret, nil
}

// Path returns the config file location.
func Path() (string, error) {
	if p := os.Getenv(EnvConfigFile); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "basaltic", "config.yaml"), nil
}

// Load reads the config file. A missing file is not an error: a CLI that
// refuses to run before it has been configured cannot show its own help.
func Load() (*File, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &File{Profiles: map[string]Profile{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return &f, nil
}

// Save writes the config file, creating its directory.
//
// 0600 because api_key is stored in plain text.
func Save(f *File) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return writeFileAtomic(path, data, 0o600)
}

// Resolve returns the profile to use, and its name.
//
// Precedence: the --profile flag, then BASALTIC_PROFILE, then
// default_profile, then a profile literally named "default".
func (f *File) Resolve(flagProfile string) (string, Profile) {
	for _, name := range []string{flagProfile, os.Getenv("BASALTIC_PROFILE"), f.DefaultProfile, "default"} {
		if name == "" {
			continue
		}
		if p, ok := f.Profiles[name]; ok {
			return name, p
		}
		// A profile asked for by name and absent is worth reporting rather
		// than silently falling through, but only when the user named it.
		if name == flagProfile || name == os.Getenv("BASALTIC_PROFILE") {
			return name, Profile{}
		}
	}
	return "default", Profile{}
}

// writeFileAtomic writes through a temporary file in the same directory, so a
// crash or a full disk cannot leave a half-written config behind.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
