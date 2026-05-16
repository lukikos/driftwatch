// Package checker provides infrastructure drift detection checks.
//
// # ElastiCache Cluster Status Check
//
// The elasticache_cluster check queries an ElastiCache-compatible HTTP
// endpoint to verify that a cluster reports the expected status.
//
// Configuration parameters:
//
//	endpoint    (required) Base URL of the ElastiCache API or mock server.
//	cluster_id  (required) The ElastiCache cluster identifier to inspect.
//	expected    (optional) Expected cluster status string. Defaults to "available".
//
// Example YAML:
//
//	checks:
//	  - name: cache-cluster-status
//	    type: elasticache_cluster
//	    params:
//	      endpoint: "http://localhost:4566"
//	      cluster_id: "my-cache-cluster"
//	      expected: "available"
package checker
