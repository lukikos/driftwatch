// Package checker provides the mount_point check type.
//
// # Mount Point Check
//
// The mount_point check verifies whether a given filesystem path is
// currently mounted or unmounted. This is useful for detecting drift
// in infrastructure where specific volumes (e.g. NFS, EBS, tmpfs)
// are expected to be present.
//
// # Configuration Fields
//
//	- path (string, required): absolute path of the expected mount point
//	- expected (string, optional): "mounted" (default) or "unmounted"
//
// # Example
//
//	checks:
//	  - name: data-volume-mounted
//	    type: mount_point
//	    fields:
//	      path: /mnt/data
//	      expected: mounted
//
package checker
