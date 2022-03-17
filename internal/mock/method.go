package mock

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Method struct {
	MethodClientBuilder internal.UnaryClientBuilder
	In                  internal.MessageDesc
	Out                 internal.MessageDesc
}

func (m Method) ClientBuilder() internal.UnaryClientBuilder {
	return m.MethodClientBuilder
}

func (m Method) Input() internal.MessageDesc {
	return m.In
}

func (m Method) Output() internal.MessageDesc {
	return m.Out
}

type MethodLoader struct {
	Methods map[internal.MethodContext]internal.UnaryMethod
}

func (ml MethodLoader) Load(methodCtx internal.MethodContext) (
	internal.UnaryMethod,
	error,
) {
	m, exists := ml.Methods[methodCtx]
	if !exists {
		ident := fmt.Sprintf(
			"MethodContext{address: %s, service: %v, method: %v}",
			methodCtx.Address(),
			methodCtx.Service(),
			methodCtx.Method(),
		)
		err := &internal.NotFound{Type: "stage", Ident: ident}
		return nil, err
	}
	return m, nil
}
