// Package checker provides infrastructure drift detection for a variety of
// check types. This file documents the s3_bucket_access check type.
//
// # s3_bucket_access
//
// Performs an HTTP HEAD request against an S3-compatible bucket URL and
// compares the response status code to an expected value. No AWS credentials
// or SDK dependency is required — the check works with any publicly accessible
// or pre-signed S3 URL.
//
// Fields:
//
//	url             (required) Full URL of the S3 bucket or object to probe.
//	expected_status (optional) Expected HTTP status code string. Defaults to "200".
//
// Example YAML:
//
//	- name: public-assets-bucket
//	  type: s3_bucket_access
//	  fields:
//	    url: https://my-bucket.s3.amazonaws.com/
//	    expected_status: "200"
package checker
