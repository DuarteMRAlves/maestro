// Package types defines generic types that are used in the maestro API. These
// types are independent of the technology used to send the messages and will
// be received and returned by the InternalAPI for the maestro server.
//
// The types are annotated with the following references:
// * optional - the field may be omitted.
// * conflicts: <field> - the field should not be defined if <field> is specified.
package types
