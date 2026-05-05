// Package runner orchestrates periodic drift checks and dispatches
// webhook alerts when drift is detected.
package runner

import (
	"context"
	"log"
	"time"

	"github.com/example/driftwatch/internal/checker"
	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/webhook"
)

// Runner periodically evaluates all configured checks and sends alerts.
type Runner struct {
	cfg     *config.Config
	checker *checker.Checker
	hook    *webhook.Client
}

// New creates a Runner wired to the provided config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:     cfg,
		checker: checker.New(),
		hook:    webhook.New(cfg.WebhookURL),
	}
}

// Run blocks until ctx is cancelled, executing checks every cfg.Interval.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	log.Printf("driftwatch started — interval %s, %d check(s)",
		r.cfg.Interval, len(r.cfg.Checks))

	// Run once immediately before waiting for the first tick.
	r.runChecks(ctx)

	for {
		select {
		case <-ticker.C:
			r.runChecks(ctx)
		case <-ctx.Done():
			log.Println("driftwatch shutting down")
			return
		}
	}
}

// runChecks iterates over all configured checks and alerts on drift.
func (r *Runner) runChecks(ctx context.Context) {
	for _, chk := range r.cfg.Checks {
		drifted, detail, err := r.checker.Run(chk)
		if err != nil {
			log.Printf("[ERROR] check %q: %v", chk.Name, err)
			continue
		}
		if drifted {
			log.Printf("[DRIFT] %s: %s", chk.Name, detail)
			if err := r.hook.Send(ctx, chk.Name, detail); err != nil {
				log.Printf("[ERROR] webhook for %q: %v", chk.Name, err)
			}
		} else {
			log.Printf("[OK]    %s", chk.Name)
		}
	}
}
