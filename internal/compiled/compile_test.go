package compiled

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/spec"
	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input        *spec.Pipeline
		expected     *Pipeline
		methodLoader MethodLoaderFunc
	}{
		"required fields": {
			input: &spec.Pipeline{
				Name: "pipeline",
				Stages: []*spec.Stage{
					{
						Name:          "stage-1",
						MethodContext: spec.MethodContext{Address: "address-1"},
					},
					{
						Name:          "stage-2",
						MethodContext: spec.MethodContext{Address: "address-2"},
					},
					{
						Name:          "stage-3",
						MethodContext: spec.MethodContext{Address: "address-3"},
					},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						TargetStage: "stage-2",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						TargetStage: "stage-3",
					},
				},
			},
			expected: &Pipeline{
				name: createPipelineName("pipeline"),
				mode: internal.OfflineExecution,
				stages: stageGraph{
					createStageName("stage-1"): &Stage{
						name:    createStageName("stage-1"),
						address: internal.NewAddress("address-1"),
						method:  testStage1Method{},
						inputs:  []*internal.Link{},
						outputs: []*internal.Link{
							createLink(
								createLinkName("1-to-2"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.MessageField{},
								),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.MessageField{},
								),
							),
						},
					},
					createStageName("stage-2"): &Stage{
						name:    createStageName("stage-2"),
						address: internal.NewAddress("address-2"),
						method:  testStage2Method{},
						inputs: []*internal.Link{
							createLink(
								createLinkName("1-to-2"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.MessageField{},
								),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.MessageField{},
								),
							),
						},
						outputs: []*internal.Link{
							createLink(
								createLinkName("2-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.MessageField{},
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.MessageField{},
								),
							),
						},
					},
					createStageName("stage-3"): &Stage{
						name:    createStageName("stage-3"),
						address: internal.NewAddress("address-3"),
						method:  testStage3Method{},
						inputs: []*internal.Link{
							createLink(
								createLinkName("2-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.MessageField{},
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.MessageField{},
								),
							),
						},
						outputs: []*internal.Link{},
					},
				},
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				ctx3 := internal.NewMethodContext(
					internal.NewAddress("address-3"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
					ctx3: testStage3Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"all fields": {
			input: &spec.Pipeline{
				Name: "pipeline",
				Mode: spec.OnlineExecution,
				Stages: []*spec.Stage{
					{
						Name: "stage-1",
						MethodContext: spec.MethodContext{
							Address: "address-1",
							Service: "service-1",
							Method:  "method-1",
						},
					},
					{
						Name: "stage-2",
						MethodContext: spec.MethodContext{
							Address: "address-2",
							Service: "service-2",
							Method:  "method-2",
						},
					},
					{
						Name: "stage-3",
						MethodContext: spec.MethodContext{
							Address: "address-3",
							Service: "service-3",
							Method:  "method-3",
						},
					},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						SourceField: "field1",
						TargetStage: "stage-2",
						TargetField: "field2",
					},
					{
						Name:        "1-to-3",
						SourceStage: "stage-1",
						SourceField: "field1",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						SourceField: "field1",
						TargetStage: "stage-3",
						TargetField: "field2",
					},
				},
			},
			expected: &Pipeline{
				name: createPipelineName("pipeline"),
				mode: internal.OnlineExecution,
				stages: stageGraph{
					createStageName("stage-1"): &Stage{
						name:    createStageName("stage-1"),
						address: internal.NewAddress("address-1"),
						method:  testStage1Method{},
						inputs:  []*internal.Link{},
						outputs: []*internal.Link{
							createLink(
								createLinkName("1-to-2"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.NewMessageField("field2"),
								),
							),
							createLink(
								createLinkName("1-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.NewMessageField("field1"),
								),
							),
						},
					},
					createStageName("stage-2"): &Stage{
						name:    createStageName("stage-2"),
						address: internal.NewAddress("address-2"),
						method:  testStage2Method{},
						inputs: []*internal.Link{
							createLink(
								createLinkName("1-to-2"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.NewMessageField("field2"),
								),
							),
						},
						outputs: []*internal.Link{
							createLink(
								createLinkName("2-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.NewMessageField("field2"),
								),
							),
						},
					},
					createStageName("stage-3"): &Stage{
						name:    createStageName("stage-3"),
						address: internal.NewAddress("address-3"),
						method:  testStage3Method{},
						inputs: []*internal.Link{
							createLink(
								createLinkName("1-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-1"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.NewMessageField("field1"),
								),
							),
							createLink(
								createLinkName("2-to-3"),
								internal.NewLinkEndpoint(
									createStageName("stage-2"),
									internal.NewMessageField("field1"),
								),
								internal.NewLinkEndpoint(
									createStageName("stage-3"),
									internal.NewMessageField("field2"),
								),
							),
						},
						outputs: []*internal.Link{},
					},
				},
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.NewService("service-1"),
					internal.NewMethod("method-1"),
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.NewService("service-2"),
					internal.NewMethod("method-2"),
				)
				ctx3 := internal.NewMethodContext(
					internal.NewAddress("address-3"),
					internal.NewService("service-3"),
					internal.NewMethod("method-3"),
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
					ctx3: testStage3Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := NewContext(tc.methodLoader)
			output, err := New(ctx, tc.input)
			if err != nil {
				t.Fatalf("new error: %s", err)
			}
			cmpOpts := cmp.AllowUnexported(
				Pipeline{},
				internal.PipelineName{},
				internal.ExecutionMode{},
				Stage{},
				internal.StageName{},
				internal.Address{},
				internal.Service{},
				internal.Method{},
				internal.Link{},
				internal.LinkName{},
				internal.LinkEndpoint{},
				internal.MessageField{},
			)
			if diff := cmp.Diff(tc.expected, output, cmpOpts); diff != "" {
				t.Fatalf("output mismatch:\n%s", diff)
			}
		})
	}
}

func TestNewIsErr(t *testing.T) {
	tests := map[string]struct {
		input        *spec.Pipeline
		validateErr  func(err error) string
		methodLoader MethodLoaderFunc
	}{
		"empty pipeline name": {
			input: &spec.Pipeline{},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyPipelineName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptyPipelineName)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				t.Fatalf("No such method: %s", methodCtx)
				return nil, nil
			},
		},
		"empty stage name": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{MethodContext: spec.MethodContext{Address: "address"}},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyStageName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptyStageName)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				t.Fatalf("No such method: %s", methodCtx)
				return nil, nil
			},
		},
		"empty address": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyAddress) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptyAddress)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				t.Fatalf("No such method: %s", methodCtx)
				return nil, nil
			},
		},
		"empty link name": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{SourceStage: "stage-1", TargetStage: "stage-2"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyLinkName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptyLinkName)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"empty link source name": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{Name: "1-to-2", TargetStage: "stage-2"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptySourceName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptySourceName)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"empty link target name": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{Name: "1-to-2", SourceStage: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyTargetName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEmptyTargetName)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"equal link source and target": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{Name: "1-to-2", SourceStage: "stage-1", TargetStage: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEqualSourceAndTarget) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, err, errEqualSourceAndTarget)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"source does not exist": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{Name: "1-to-2", SourceStage: "stage-3", TargetStage: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *stageNotFound
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *stageNotFound, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &stageNotFound{name: "stage-3"}
				cmpOpts := cmp.AllowUnexported(stageNotFound{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"target does not exist": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{Name: "1-to-2", SourceStage: "stage-1", TargetStage: "stage-4"},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *stageNotFound
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *stageNotFound, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &stageNotFound{name: "stage-4"}
				cmpOpts := cmp.AllowUnexported(stageNotFound{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"new link set full message": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
					{Name: "stage-3", MethodContext: spec.MethodContext{Address: "address-3"}},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						TargetStage: "stage-2"},
					{
						Name:        "1-to-3",
						SourceStage: "stage-1",
						SourceField: "field1",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						TargetStage: "stage-3",
					},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *linkSetsFullMessage
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *linkSetsFullMessage, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &linkSetsFullMessage{name: "2-to-3"}
				cmpOpts := cmp.AllowUnexported(linkSetsFullMessage{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				ctx3 := internal.NewMethodContext(
					internal.NewAddress("address-3"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
					ctx3: testStage3Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"old link set full message": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
					{Name: "stage-3", MethodContext: spec.MethodContext{Address: "address-3"}},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						TargetStage: "stage-2"},
					{
						Name:        "1-to-3",
						SourceStage: "stage-1",
						TargetStage: "stage-3",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						SourceField: "field1",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *linkSetsFullMessage
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *linkSetsFullMessage, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &linkSetsFullMessage{name: "1-to-3"}
				cmpOpts := cmp.AllowUnexported(linkSetsFullMessage{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				ctx3 := internal.NewMethodContext(
					internal.NewAddress("address-3"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
					ctx3: testStage3Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"new and old links set same": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
					{Name: "stage-3", MethodContext: spec.MethodContext{Address: "address-3"}},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						TargetStage: "stage-2"},
					{
						Name:        "1-to-3",
						SourceStage: "stage-1",
						SourceField: "field2",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						SourceField: "field1",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *linksSetSameField
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *linksSetSameField, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &linksSetSameField{A: "2-to-3", B: "1-to-3", field: "field1"}
				cmpOpts := cmp.AllowUnexported(linksSetSameField{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				ctx3 := internal.NewMethodContext(
					internal.NewAddress("address-3"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
					ctx3: testStage3Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"incompatible message descriptor": {
			input: &spec.Pipeline{
				Name: "Pipeline",
				Stages: []*spec.Stage{
					{Name: "stage-1", MethodContext: spec.MethodContext{Address: "address-1"}},
					{Name: "stage-2", MethodContext: spec.MethodContext{Address: "address-2"}},
				},
				Links: []*spec.Link{
					{
						Name:        "1-to-2",
						SourceStage: "stage-1",
						SourceField: "field1",
						TargetStage: "stage-2",
					},
				},
			},
			validateErr: func(err error) string {
				var concreteErr *incompatibleMessageDesc
				if !errors.As(err, &concreteErr) {
					format := "Wrong error type: expected *incompatibleMessageDesc, got %s"
					return fmt.Sprintf(format, reflect.TypeOf(err))
				}
				expErr := &incompatibleMessageDesc{A: testInnerValDesc{}, B: testOuterValDesc{}}
				cmpOpts := cmp.AllowUnexported(incompatibleMessageDesc{})
				if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
					return fmt.Sprintf("error mismatch:\n%s", diff)
				}
				return ""
			},
			methodLoader: func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
				ctx1 := internal.NewMethodContext(
					internal.NewAddress("address-1"),
					internal.Service{},
					internal.Method{},
				)
				ctx2 := internal.NewMethodContext(
					internal.NewAddress("address-2"),
					internal.Service{},
					internal.Method{},
				)
				mapper := map[internal.MethodContext]internal.UnaryMethod{
					ctx1: testStage1Method{},
					ctx2: testStage2Method{},
				}
				s, ok := mapper[methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := NewContext(tc.methodLoader)
			output, err := New(ctx, tc.input)
			if err == nil {
				t.Fatalf("error expected but received nil")
			}
			if output != nil {
				t.Fatalf("expected nil output")
			}
			if msg := tc.validateErr(err); msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}

func createPipelineName(name string) internal.PipelineName {
	pipelineName, err := internal.NewPipelineName(name)
	if err != nil {
		panic(fmt.Sprintf("create pipeline name %s: %s", name, err))
	}
	return pipelineName
}

func createStageName(name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		panic(fmt.Sprintf("create stage name %s: %s", name, err))
	}
	return stageName
}

func createLinkName(name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		panic(fmt.Sprintf("create link name %s: %s", name, err))
	}
	return linkName
}

func createLink(
	name internal.LinkName, source, target internal.LinkEndpoint,
) *internal.Link {
	l := internal.NewLink(name, source, target)
	return &l
}

type testStage1Method struct{}

func (m testStage1Method) ClientBuilder() internal.UnaryClientBuilder {
	return nil
}

func (m testStage1Method) Input() internal.MessageDesc {
	return testEmptyDesc{}
}

func (m testStage1Method) Output() internal.MessageDesc {
	return testOuterValDesc{}
}

type testStage2Method struct{}

func (m testStage2Method) ClientBuilder() internal.UnaryClientBuilder {
	return nil
}

func (m testStage2Method) Input() internal.MessageDesc {
	return testOuterValDesc{}
}

func (m testStage2Method) Output() internal.MessageDesc {
	return testOuterValDesc{}
}

type testStage3Method struct{}

func (m testStage3Method) ClientBuilder() internal.UnaryClientBuilder {
	return nil
}

func (m testStage3Method) Input() internal.MessageDesc {
	return testOuterValDesc{}
}

func (m testStage3Method) Output() internal.MessageDesc {
	return testEmptyDesc{}
}

type testEmptyDesc struct{}

func (d testEmptyDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testEmptyDesc)
	return ok
}

func (d testEmptyDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return nil }
}

func (d testEmptyDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	panic("method get field should not be called for testEmptyDesc")
}

// Represents a descriptor of a message with two inner fields: field1 and field2.
// Each field is associated with a descriptor of type testInnerValDesc
type testOuterValDesc struct{}

func (d testOuterValDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testOuterValDesc)
	return ok
}

func (d testOuterValDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return nil }
}

func (d testOuterValDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	switch f.Unwrap() {
	case "field1", "field2":
		return testInnerValDesc{}, nil
	default:
		panic(fmt.Sprintf("Unknown field for testOuterValDesc: %s", f.Unwrap()))
	}
}

func (d testOuterValDesc) String() string {
	return "testOuterValDesc"
}

type testInnerValDesc struct{}

func (d testInnerValDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testInnerValDesc)
	return ok
}

func (d testInnerValDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return nil }
}

func (d testInnerValDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	panic("method get field should not be called for testInnerValDesc")
}

func (d testInnerValDesc) String() string {
	return "testInnerValDesc"
}
