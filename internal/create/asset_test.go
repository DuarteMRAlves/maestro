package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	tests := map[string]struct {
		name     internal.AssetName
		image    internal.Image
		expected internal.Asset
	}{
		"all fields": {
			name:     createAssetName(t, "some-name"),
			image:    internal.NewImage("some-image"),
			expected: createAsset(t, "some-name", "some-image"),
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				storage := mock.AssetStorage{
					Assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				err := createFn(tc.name, tc.image)
				if err != nil {
					t.Fatalf("create error: %s", err)
				}

				if diff := cmp.Diff(1, len(storage.Assets)); diff != "" {
					t.Fatalf("number of assets mismatch:\n%s", diff)
				}

				asset, exists := storage.Assets[tc.expected.Name()]
				if !exists {
					t.Fatalf("created asset does not exist in storage")
				}
				cmpAsset(t, tc.expected, asset, "created asset")
			},
		)
	}
}

func TestCreateAsset_Err(t *testing.T) {
	tests := map[string]struct {
		name    internal.AssetName
		image   internal.Image
		isError error
	}{
		"empty name": {name: createAssetName(t, ""), isError: EmptyAssetName},
		"empty image": {
			name:    createAssetName(t, "some-name"),
			image:   internal.NewImage(""),
			isError: EmptyImageName,
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				storage := mock.AssetStorage{
					Assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				err := createFn(tc.name, tc.image)
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !errors.Is(err, tc.isError) {
					format := "Wrong error: expected %s, got %s"
					t.Fatalf(format, tc.isError, err)
				}
				if diff := cmp.Diff(0, len(storage.Assets)); diff != "" {
					t.Fatalf("number of assets mismatch:\n%s", diff)
				}
			},
		)
	}
}

func TestCreateAsset_AlreadyExists(t *testing.T) {
	name := "some-name"
	assetName := createAssetName(t, name)
	image1 := internal.NewImage("some-image-1")
	image2 := internal.NewImage("some-image-2")
	expected := internal.NewAsset(assetName, image1)
	storage := mock.AssetStorage{
		Assets: map[internal.AssetName]internal.Asset{},
	}

	createFn := Asset(storage)

	err := createFn(assetName, image1)
	if err != nil {
		t.Fatalf("first create error: %s", err)
	}
	if diff := cmp.Diff(1, len(storage.Assets)); diff != "" {
		t.Fatalf("first create number of assets mismatch:\n%s", diff)
	}
	asset, exists := storage.Assets[expected.Name()]
	if !exists {
		t.Fatalf("first created asset does not exist in storage")
	}
	cmpAsset(t, expected, asset, "first created asset")

	err = createFn(assetName, image2)
	if err == nil {
		t.Fatalf("expected create error but got none")
	}
	var alreadyExists *internal.AlreadyExists
	if !errors.As(err, &alreadyExists) {
		format := "Wrong error type: expected *internal.AlreadyExists, got %s"
		t.Fatalf(format, reflect.TypeOf(err))
	}
	expError := &internal.AlreadyExists{Type: "asset", Ident: name}
	if diff := cmp.Diff(expError, alreadyExists); diff != "" {
		t.Fatalf("error mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(1, len(storage.Assets)); diff != "" {
		t.Fatalf("second create number of assets mismatch:\n%s", diff)
	}
	asset, exists = storage.Assets[expected.Name()]
	if !exists {
		t.Fatalf("second created asset does not exist in storage")
	}
	cmpAsset(t, expected, asset, "second created asset")
}

func createAssetName(t *testing.T, assetName string) internal.AssetName {
	name, err := internal.NewAssetName(assetName)
	if err != nil {
		t.Fatalf("create asset name %s: %s", assetName, err)
	}
	return name
}

func createAsset(t *testing.T, assetName, image string) internal.Asset {
	name := createAssetName(t, assetName)
	img := internal.NewImage(image)
	return internal.NewAsset(name, img)
}

func cmpAsset(t *testing.T, x, y internal.Asset, msg string, args ...interface{}) {
	cmpOpts := cmp.AllowUnexported(
		internal.Asset{},
		internal.AssetName{},
		internal.Image{},
	)
	if diff := cmp.Diff(x, y, cmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
