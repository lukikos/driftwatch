// Package checker evaluates individual drift checks defined in configuration.
package checker

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/driftwatch/internal/config"
)

// Result holds the outcome of a single check.
type Result struct {
	CheckName string
	Drifted   bool
	Message   string
}

// Checker runs configured drift checks.
type Checker struct {
	checks []config.Check
}

// New creates a Checker from the provided checks slice.
func New(checks []config.Check) *Checker {
	return &Checker{checks: checks}
}

// Run executes all checks and returns their results.
func (c *Checker) Run() ([]Result, error) {
	var results []Result
	for _, chk := range c.checks {
		r, err := c.runOne(chk)
		if err != nil {
			return nil, fmt.Errorf("check %q: %w", chk.Name, err)
		}
		results = append(results, r)
	}
	return results, nil
}

func (c *Checker) runOne(chk config.Check) (Result, error) {
	switch chk.Type {
	case "file_hash":
		drifted, msg, err := checkFileHash(chk.Path, chk.Expected)
		if err != nil {
			return Result{}, err
		}
		return Result{CheckName: chk.Name, Drifted: drifted, Message: msg}, nil
	case "env_var":
		drifted, msg, err := checkEnvVar(chk.EnvVar, chk.Expected)
		if err != nil {
			return Result{}, err
		}
		return Result{CheckName: chk.Name, Drifted: drifted, Message: msg}, nil
	case "http_status":
		drifted, msg, err := checkHTTPStatus(chk.URL, chk.ExpectedStatus)
		if err != nil {
			return Result{}, err
		}
		return Result{CheckName: chk.Name, Drifted: drifted, Message: msg}, nil
	default:
		return Result{}, fmt.Errorf("unknown check type: %s", chk.Type)
	}
}

func checkFileHash(path, expected string) (bool, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, "", fmt.Errorf("hashing file: %w", err)
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if actual != expected {
		return true, fmt.Sprintf("expected hash %s, got %s", expected, actual), nil
	}
	return false, "", nil
}

func checkEnvVar(name, expected string) (bool, string, error) {
	if name == "" {
		return false, "", fmt.Errorf("env_var check requires a non-empty env_var field")
	}
	actual := os.Getenv(name)
	if actual != expected {
		return true, fmt.Sprintf("expected %q, got %q", expected, actual), nil
	}
	return false, "", nil
}
