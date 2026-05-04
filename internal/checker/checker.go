package checker

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/yourusername/driftwatch/internal/config"
)

// Result holds the outcome of a single check.
type Result struct {
	Name    string
	Drifted bool
	Expected string
	Actual   string
	Err     error
}

// Checker evaluates a set of checks against their expected state.
type Checker struct {
	checks []config.Check
}

// New creates a new Checker from the provided check definitions.
func New(checks []config.Check) *Checker {
	return &Checker{checks: checks}
}

// RunAll executes all checks and returns a slice of Results.
func (c *Checker) RunAll() []Result {
	results := make([]Result, 0, len(c.checks))
	for _, chk := range c.checks {
		results = append(results, c.run(chk))
	}
	return results
}

func (c *Checker) run(chk config.Check) Result {
	switch chk.Type {
	case "file_hash":
		return checkFileHash(chk)
	case "env_var":
		return checkEnvVar(chk)
	default:
		return Result{
			Name: chk.Name,
			Err:  fmt.Errorf("unknown check type: %s", chk.Type),
		}
	}
}

func checkFileHash(chk config.Check) Result {
	data, err := os.ReadFile(chk.Target)
	if err != nil {
		return Result{Name: chk.Name, Err: fmt.Errorf("reading file: %w", err)}
	}
	actual := fmt.Sprintf("%x", sha256.Sum256(data))
	return Result{
		Name:     chk.Name,
		Drifted:  actual != chk.Expected,
		Expected: chk.Expected,
		Actual:   actual,
	}
}

func checkEnvVar(chk config.Check) Result {
	actual := os.Getenv(chk.Target)
	return Result{
		Name:     chk.Name,
		Drifted:  actual != chk.Expected,
		Expected: chk.Expected,
		Actual:   actual,
	}
}
