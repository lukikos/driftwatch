// Package checker evaluates individual drift checks defined in configuration.
package checker

import (
	"fmt"

	"github.com/user/driftwatch/internal/config"
)

// Result holds the outcome of a single drift check.
type Result struct {
	CheckName string
	Drifted   bool
	Message   string
}

// Checker runs configured checks and returns results.
type Checker struct {
	checks []config.Check
}

// New creates a Checker from the provided check configurations.
func New(checks []config.Check) *Checker {
	return &Checker{checks: checks}
}

// RunAll executes all configured checks and returns their results.
func (c *Checker) RunAll() []Result {
	results := make([]Result, 0, len(c.checks))
	for _, chk := range c.checks {
		result := c.run(chk)
		results = append(results, result)
	}
	return results
}

// run dispatches a single check by type.
func (c *Checker) run(chk config.Check) Result {
	var (
		drifted bool
		msg     string
		err     error
	)

	switch chk.Type {
	case "env_var":
		drifted, msg, err = checkEnvVar(chk.Fields)
	case "file_hash":
		drifted, msg, err = checkFileHash(chk.Fields)
	case "http_status":
		drifted, msg, err = checkHTTPStatus(chk.Fields)
	case "process_running":
		drifted, msg, err = checkProcessRunning(chk.Fields)
	case "port_open":
		drifted, msg, err = checkPortOpen(chk.Fields)
	default:
		err = fmt.Errorf("unknown check type: %s", chk.Type)
	}

	if err != nil {
		return Result{CheckName: chk.Name, Drifted: true, Message: err.Error()}
	}
	return Result{CheckName: chk.Name, Drifted: drifted, Message: msg}
}
