package resources

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

func ParseFiles(files []string) ([]*Resource, error) {
	resources := make([]*Resource, 0)

	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errdefs.InvalidArgumentWithError(err)
		}
		dataResources, err := ParseBytes(data)
		if err != nil {
			return nil, err
		}
		resources = append(resources, dataResources...)
	}

	return resources, nil
}

func ParseBytes(data []byte) ([]*Resource, error) {
	var err error

	reader := bytes.NewReader(data)

	dec := yaml.NewDecoder(reader)

	resources := make([]*Resource, 0)

	for {
		r := &Resource{}
		if err = dec.Decode(&r); err != nil {
			break
		}
		resources = append(resources, r)
	}

	if err == io.EOF {
		return resources, nil
	} else {
		return nil, errdefs.InternalWithMsg("unmarshal resource: %v", err)
	}
}
