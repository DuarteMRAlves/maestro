// Package util offers common mechanisms for all other packages.
// It offers utility functions to validate the stage of the program.
// It has functions to verify if the received arguments follow certain conditions
// and return errors otherwise.
//
// The package offers uniform error handling, as it allows validations to be
// associated to the same error throughout the program.
//
// The package also contains methods to help with tests.
//
// The functions help set up independent servers with different ports
// so that tests can run in parallel safely.
// It also offers common rules for creating fields for the several
// domain objects.
package util
