package errdefs

import "errors"

const dummyErrMsg = "dummy error message"

var dummyErr = errors.New(dummyErrMsg)
