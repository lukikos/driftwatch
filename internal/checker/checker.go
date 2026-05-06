// Package checker evaluates individual drift checks defined in the configuration.
// Each check has a type (file_hash, env_var, http_status, process) and a set of
// key/value parameters specific to that type.
package checker

import (
	"fmt"

	"github.com/user/driftwatch/internal/config"
)

// Result holds the outcome of a single drift check.
type Result struct {
	Name    string
	Drifted bool
	Message string
}

// Checker runs the configured checks and returns their results.
type Checker struct {
	checks []config.Check
}

// New creates a Checker from the provided check configurations.
func New(checks []config.Check) *Checker {
	return &Checker{checks: checks}
}

// RunAll executes every configured check and returns a slice of Results.
func (c *Checker) RunAll() []Result {
	results := make([]Result, 0, len(c.checks))
	for _, chk := range c.checks {
		r := c.run(chk)
		results = append(results, r)
	}
	return results
}

func (c *Checker) run(chk config.Check) Result {
	var (
		drifted bool
		msg     string
		err     error
	)

	switch chk.Type {
	case "file_hash":
		drifted, msg, err = checkFileHash(chk.Params)
	case "env_var":
		drifted, msg, err = checkEnvVar(chk.Params)
	case "http_status":
		drifted, msg, err = checkHTTPStatus(chk.Params)
	case "process":
		drifted, msg, err = checkProcessRunning(chk.Params)
	default:
		err = fmt.Errorf("unknown check type: %q", chk.Type)
	}

	if err != nil {
		return Result{
			Name:    chk.Name,
			Drifted: true,
			Message: fmt.Sprintf("check error: %v", err),
		}
	}

	return Result{
		Name:    chk.Name,
		Drifted: drifted,
		Message: msg,
	}
}
