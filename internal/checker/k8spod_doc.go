// Package checker provides infrastructure drift detection logic.
//
// # k8s_pod check
//
// The k8s_pod check type uses kubectl to verify that a Kubernetes pod
// matching a given name prefix exists in a namespace and has the expected
// status phase.
//
// Configuration fields:
//
//	pod_prefix  (required) - prefix of the pod name to match
//	namespace   (optional) - Kubernetes namespace, defaults to "default"
//	expected    (optional) - expected pod phase, defaults to "Running"
//
// Example:
//
//	checks:
//	  - name: api-pod-running
//	    type: k8s_pod
//	    fields:
//	      pod_prefix: api-server
//	      namespace: production
//	      expected: Running
package checker
