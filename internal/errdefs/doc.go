// Package errdefs implements the error handling mechanisms used throughout the
// code base.
//
// The package defines a set of interfaces for the several errors. These errors
// have a one to one correspondence with grpc errors, defined in
// https://grpc.github.io/grpc/core/md_doc_statuscodes.html.
// The errors should be used in the same conditions as the ones defined in the
// previous link.
//
// The error interfaces defined in this package are also here implemented, as
// well as a set of utility functions that should be used to create and
// manage any generated errors.
package errdefs
