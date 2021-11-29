package create

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Resource struct {
	Kind string
	Spec map[string]string
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource{Kind:%v,Spec:%v", r.Kind, r.Spec)
}

func ParseResources(files []string) ([]*Resource, error) {
	resources := make([]*Resource, 0)

	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read file %v: %v", f, err)
		}
		dataResources, err := UnmarshalResources(data)
		if err != nil {
			return nil, err
		}
		resources = append(resources, dataResources...)
	}

	return resources, nil
}

func UnmarshalResources(data []byte) ([]*Resource, error) {
	var (
		curr Resource
		err  error
	)

	reader := bytes.NewReader(data)

	dec := yaml.NewDecoder(reader)

	resources := make([]*Resource, 0)

	for {
		if err = dec.Decode(&curr); err != nil {
			break
		}

		r := &Resource{}
		copyResource(r, &curr)
		resources = append(resources, r)
	}

	if err == io.EOF {
		return resources, nil
	} else {
		return nil, errdefs.InternalWithMsg("unmarshal resource: %v", err)
	}
}

func copyResource(dst *Resource, src *Resource) {
	dst.Kind = src.Kind
	dst.Spec = make(map[string]string, len(src.Spec))
	for k, v := range src.Spec {
		dst.Spec[k] = v
	}
}
