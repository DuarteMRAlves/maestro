package compiled

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/spec"
	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input        *spec.Pipeline
		expected     *Pipeline
		methodLoader MethodLoaderFunc
	}{
		"linear specification": {
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
				name: PipelineName{val: "pipeline"},
				mode: OfflineExecution,
				stages: stageGraph{
					StageName{val: "stage-1:aux-source"}: &Stage{
						name:  StageName{val: "stage-1:aux-source"},
						sType: StageTypeSource,
						ictx: &InvocationContext{
							address:     Address("address-1"),
							unaryMethod: testLinearStage1Method{},
						},
						inputs: []*Link{},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
					},
					StageName{val: "stage-1"}: &Stage{
						name:  StageName{val: "stage-1"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-1"),
							unaryMethod: testLinearStage1Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "1-to-2"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-2"}},
							},
						},
					},
					StageName{val: "stage-2"}: &Stage{
						name:  StageName{val: "stage-2"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-2"),
							unaryMethod: testLinearStage2Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "1-to-2"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-2"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "2-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-2"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3"}},
							},
						},
					},
					StageName{val: "stage-3"}: &Stage{
						name:  StageName{val: "stage-3"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-3"),
							unaryMethod: testLinearStage3Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "2-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-2"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-sink-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-3"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3:aux-sink"}},
							},
						},
					},
					StageName{val: "stage-3:aux-sink"}: &Stage{
						name:  StageName{val: "stage-3:aux-sink"},
						sType: StageTypeSink,
						ictx: &InvocationContext{
							address:     Address("address-3"),
							unaryMethod: testLinearStage3Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-sink-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-3"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3:aux-sink"}},
							},
						},
						outputs: []*Link{},
					},
				},
			},
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := MethodContext{address: Address("address-1")}
				ctx2 := MethodContext{address: Address("address-2")}
				ctx3 := MethodContext{address: Address("address-3")}

				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
					ctx3: testLinearStage3Method{},
				}
				s, ok := mapper[*methodCtx]
				if !ok {
					panic(fmt.Sprintf("No such method: %s", methodCtx))
				}
				return s, nil
			},
		},
		"split and merge": {
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
						TargetStage: "stage-2",
					},
					{
						Name:        "1-to-3",
						SourceStage: "stage-1",
						TargetStage: "stage-3",
						TargetField: "field1",
					},
					{
						Name:        "2-to-3",
						SourceStage: "stage-2",
						TargetStage: "stage-3",
						TargetField: "field2",
					},
				},
			},
			expected: &Pipeline{
				name: PipelineName{val: "pipeline"},
				mode: OnlineExecution,
				stages: stageGraph{
					StageName{val: "stage-1:aux-source"}: &Stage{
						name:  StageName{val: "stage-1:aux-source"},
						sType: StageTypeSource,
						ictx: &InvocationContext{
							address:     Address("address-1"),
							service:     Service("service-1"),
							method:      Method("method-1"),
							unaryMethod: testSplitAndMergeStage1Method{},
						},
						inputs: []*Link{},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
					},
					StageName{val: "stage-1"}: &Stage{
						name:  StageName{val: "stage-1"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-1"),
							service:     Service("service-1"),
							method:      Method("method-1"),
							unaryMethod: testSplitAndMergeStage1Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{"stage-1:aux-split-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
							},
						},
					},
					StageName{val: "stage-1:aux-split"}: &Stage{
						name:  StageName{val: "stage-1:aux-split"},
						sType: StageTypeSplit,
						ictx: &InvocationContext{
							address:     Address("address-1"),
							service:     Service("service-1"),
							method:      Method("method-1"),
							unaryMethod: testSplitAndMergeStage1Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{"stage-1:aux-split-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "1-to-2"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-2"}},
							},
							{
								name:   LinkName{val: "1-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: MessageField{val: "field1"},
								},
							},
						},
					},
					StageName{val: "stage-2"}: &Stage{
						name:  StageName{val: "stage-2"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-2"),
							service:     Service("service-2"),
							method:      Method("method-2"),
							unaryMethod: testSplitAndMergeStage2Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "1-to-2"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-2"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "2-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-2"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: MessageField{val: "field2"},
								},
							},
						},
					},
					StageName{val: "stage-3:aux-merge"}: &Stage{
						name:  StageName{val: "stage-3:aux-merge"},
						sType: StageTypeMerge,
						ictx: &InvocationContext{
							address:     Address("address-3"),
							service:     Service("service-3"),
							method:      Method("method-3"),
							unaryMethod: testSplitAndMergeStage3Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "1-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: MessageField{val: "field1"},
								},
							},
							{
								name:   LinkName{val: "2-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-2"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: MessageField{val: "field2"},
								},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-merge-link"},
								source: &LinkEndpoint{stage: StageName{"stage-3:aux-merge"}},
								target: &LinkEndpoint{stage: StageName{"stage-3"}},
							},
						},
					},
					StageName{val: "stage-3"}: &Stage{
						name:  StageName{val: "stage-3"},
						sType: StageTypeUnary,
						ictx: &InvocationContext{
							address:     Address("address-3"),
							service:     Service("service-3"),
							method:      Method("method-3"),
							unaryMethod: testSplitAndMergeStage3Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-merge-link"},
								source: &LinkEndpoint{stage: StageName{"stage-3:aux-merge"}},
								target: &LinkEndpoint{stage: StageName{"stage-3"}},
							},
						},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-sink-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-3"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3:aux-sink"}},
							},
						},
					},
					StageName{val: "stage-3:aux-sink"}: &Stage{
						name:  StageName{val: "stage-3:aux-sink"},
						sType: StageTypeSink,
						ictx: &InvocationContext{
							address:     Address("address-3"),
							service:     Service("service-3"),
							method:      Method("method-3"),
							unaryMethod: testSplitAndMergeStage3Method{},
						},
						inputs: []*Link{
							{
								name:   LinkName{val: "stage-3:aux-sink-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-3"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-3:aux-sink"}},
							},
						},
						outputs: []*Link{},
					},
				},
			},
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "service-1", "method-1")
				ctx2 := NewMethodContext("address-2", "service-2", "method-2")
				ctx3 := NewMethodContext("address-3", "service-3", "method-3")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testSplitAndMergeStage1Method{},
					ctx2: testSplitAndMergeStage2Method{},
					ctx3: testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[*methodCtx]
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
				PipelineName{},
				ExecutionMode{},
				Stage{},
				StageName{},
				InvocationContext{},
				Link{},
				LinkName{},
				LinkEndpoint{},
				MessageField{},
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				ctx3 := NewMethodContext("address-3", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
					ctx3: testLinearStage3Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				ctx3 := NewMethodContext("address-3", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
					ctx3: testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[*methodCtx]
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
						TargetStage: "stage-2",
					},
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				ctx3 := NewMethodContext("address-3", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
					ctx3: testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[*methodCtx]
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
			methodLoader: func(methodCtx *MethodContext) (UnaryMethod, error) {
				ctx1 := NewMethodContext("address-1", "", "")
				ctx2 := NewMethodContext("address-2", "", "")
				mapper := map[MethodContext]UnaryMethod{
					ctx1: testLinearStage1Method{},
					ctx2: testLinearStage2Method{},
				}
				s, ok := mapper[*methodCtx]
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

type testLinearStage1Method struct{}

func (m testLinearStage1Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testLinearStage1Method) Input() MessageDesc {
	return testEmptyDesc{}
}

func (m testLinearStage1Method) Output() MessageDesc {
	return testOuterValDesc{}
}

type testLinearStage2Method struct{}

func (m testLinearStage2Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testLinearStage2Method) Input() MessageDesc {
	return testOuterValDesc{}
}

func (m testLinearStage2Method) Output() MessageDesc {
	return testOuterValDesc{}
}

type testLinearStage3Method struct{}

func (m testLinearStage3Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testLinearStage3Method) Input() MessageDesc {
	return testOuterValDesc{}
}

func (m testLinearStage3Method) Output() MessageDesc {
	return testEmptyDesc{}
}

type testSplitAndMergeStage1Method struct{}

func (m testSplitAndMergeStage1Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testSplitAndMergeStage1Method) Input() MessageDesc {
	return testEmptyDesc{}
}

func (m testSplitAndMergeStage1Method) Output() MessageDesc {
	return testInnerValDesc{}
}

type testSplitAndMergeStage2Method struct{}

func (m testSplitAndMergeStage2Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testSplitAndMergeStage2Method) Input() MessageDesc {
	return testInnerValDesc{}
}

func (m testSplitAndMergeStage2Method) Output() MessageDesc {
	return testInnerValDesc{}
}

type testSplitAndMergeStage3Method struct{}

func (m testSplitAndMergeStage3Method) ClientBuilder() UnaryClientBuilder {
	return nil
}

func (m testSplitAndMergeStage3Method) Input() MessageDesc {
	return testOuterValDesc{}
}

func (m testSplitAndMergeStage3Method) Output() MessageDesc {
	return testEmptyDesc{}
}

type testEmptyDesc struct{}

func (d testEmptyDesc) Compatible(other MessageDesc) bool {
	_, ok := other.(testEmptyDesc)
	return ok
}

func (d testEmptyDesc) EmptyGen() EmptyMessageGen {
	return func() Message { return nil }
}

func (d testEmptyDesc) GetField(f MessageField) (MessageDesc, error) {
	panic("method get field should not be called for testEmptyDesc")
}

// Represents a descriptor of a message with two inner fields: field1 and field2.
// Each field is associated with a descriptor of type testInnerValDesc
type testOuterValDesc struct{}

func (d testOuterValDesc) Compatible(other MessageDesc) bool {
	_, ok := other.(testOuterValDesc)
	return ok
}

func (d testOuterValDesc) EmptyGen() EmptyMessageGen {
	return func() Message { return nil }
}

func (d testOuterValDesc) GetField(f MessageField) (MessageDesc, error) {
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

func (d testInnerValDesc) Compatible(other MessageDesc) bool {
	_, ok := other.(testInnerValDesc)
	return ok
}

func (d testInnerValDesc) EmptyGen() EmptyMessageGen {
	return func() Message { return nil }
}

func (d testInnerValDesc) GetField(f MessageField) (MessageDesc, error) {
	panic("method get field should not be called for testInnerValDesc")
}

func (d testInnerValDesc) String() string {
	return "testInnerValDesc"
}
