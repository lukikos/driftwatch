// Package checker provides infrastructure drift detection for various check types.
//
// # Symlink Check
//
// The symlink check verifies that a symbolic link exists at a specified path
// and optionally resolves to an expected target.
//
// Configuration fields:
//
//	type: symlink
//	fields:
//	  path: /etc/nginx/sites-enabled/default   # required: path to the symlink
//	  expected_target: /etc/nginx/sites-available/default  # optional: resolved target
//
// Drift is reported when:
//   - The path does not exist
//   - The path exists but is not a symlink
//   - An expected_target is set and the symlink resolves to a different path
//
// Example driftwatch.yaml entry:
//
//	checks:
//	  - name: nginx-default-site-symlink
//	    type: symlink
//	    fields:
//	      path: /etc/nginx/sites-enabled/default
//	      expected_target: /etc/nginx/sites-available/default
package checker
