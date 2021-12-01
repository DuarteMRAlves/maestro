package resources

import (
	"fmt"
)

type Resource struct {
	Kind string
	Spec map[string]string
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource{Kind:%v,Spec:%v", r.Kind, r.Spec)
}

func copyResource(dst *Resource, src *Resource) {
	dst.Kind = src.Kind
	dst.Spec = make(map[string]string, len(src.Spec))
	for k, v := range src.Spec {
		dst.Spec[k] = v
	}
}
