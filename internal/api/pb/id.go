package pb

import (
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

var emptyIdPb = protobuff.MarshalID(identifier.Empty())
