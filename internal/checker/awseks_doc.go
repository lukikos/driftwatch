// Package checker provides infrastructure drift detection checks.
//
// # EKS Cluster Status Check
//
// The eks_cluster_status check type queries an EKS-compatible API endpoint
// to verify that a named cluster is in the expected lifecycle status.
//
// Required fields:
//   - endpoint:     base URL of the EKS-compatible API
//   - cluster_name: name of the cluster to inspect
//
// Optional fields:
//   - expected: desired cluster status (default: "ACTIVE")
//
// Example configuration:
//
//	- name: prod-eks-cluster
//	  type: eks_cluster_status
//	  fields:
//	    endpoint: https://eks.us-east-1.amazonaws.com
//	    cluster_name: prod-cluster
//	    expected: ACTIVE
package checker
