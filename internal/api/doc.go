// Package api defines generic types that are used in the maestro API. These
// types are independent of the technology used to send the messages and will
// be received and returned by the InternalAPI for the maestro server.
//
// The types are annotated with the following references:
// * required - the field must be defined.
// * optional - the field may be omitted.
// * non-empty - the field should have at least on element.
// * unique - the field should have a unique value for each object.
// * conflicts: <field> - the field should not be defined if <field> is specified.
//
// This package also defines several messages, used to communicate with the
// maestro server. The messages follow the same annotation reference as the
// types.
package api
