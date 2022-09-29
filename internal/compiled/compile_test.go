package compiled

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input    *api.Pipeline
		expected *Pipeline
		resolver method.ResolveFunc
	}{
		"linear specification": {
			input: &api.Pipeline{
				Name: "pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
					{Name: "stage-3", Address: "method-3"},
				},
				Links: []*api.Link{
					{Name: "1-to-2", SourceStage: "stage-1", TargetStage: "stage-2"},
					{Name: "2-to-3", SourceStage: "stage-2", TargetStage: "stage-3"},
				},
			},
			expected: &Pipeline{
				name: PipelineName{val: "pipeline"},
				stages: stageGraph{
					StageName{val: "stage-1:aux-source"}: &Stage{
						name:    StageName{val: "stage-1:aux-source"},
						sType:   StageTypeSource,
						address: "method-1/*/*",
						desc:    testLinearStage1Method{},
						inputs:  []*Link{},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
					},
					StageName{val: "stage-1"}: &Stage{
						name:    StageName{val: "stage-1"},
						sType:   StageTypeUnary,
						address: "method-1/*/*",
						desc:    testLinearStage1Method{},
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
						name:    StageName{val: "stage-2"},
						sType:   StageTypeUnary,
						address: "method-2/*/*",
						desc:    testLinearStage2Method{},
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
						name:    StageName{val: "stage-3"},
						sType:   StageTypeUnary,
						address: "method-3/*/*",
						desc:    testLinearStage3Method{},
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
						name:    StageName{val: "stage-3:aux-sink"},
						sType:   StageTypeSink,
						address: "method-3/*/*",
						desc:    testLinearStage3Method{},
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
					"method-3/*/*": testLinearStage3Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %q", address))
				}
				return s, nil
			},
		},
		"split and merge": {
			input: &api.Pipeline{
				Name: "pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
					{Name: "stage-3", Address: "method-3"},
				},
				Links: []*api.Link{
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
				stages: stageGraph{
					StageName{val: "stage-1:aux-source"}: &Stage{
						name:    StageName{val: "stage-1:aux-source"},
						sType:   StageTypeSource,
						address: "method-1/*/*",
						desc:    testSplitAndMergeStage1Method{},
						inputs:  []*Link{},
						outputs: []*Link{
							{
								name:   LinkName{val: "stage-1:aux-source-link"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-source"}},
								target: &LinkEndpoint{stage: StageName{val: "stage-1"}},
							},
						},
					},
					StageName{val: "stage-1"}: &Stage{
						name:    StageName{val: "stage-1"},
						sType:   StageTypeUnary,
						address: "method-1/*/*",
						desc:    testSplitAndMergeStage1Method{},
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
						name:    StageName{val: "stage-1:aux-split"},
						sType:   StageTypeSplit,
						address: "method-1/*/*",
						desc:    testSplitAndMergeStage1Method{},
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
									field: message.Field("field1"),
								},
							},
						},
					},
					StageName{val: "stage-2"}: &Stage{
						name:    StageName{val: "stage-2"},
						sType:   StageTypeUnary,
						address: "method-2/*/*",
						desc:    testSplitAndMergeStage2Method{},
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
									field: message.Field("field2"),
								},
							},
						},
					},
					StageName{val: "stage-3:aux-merge"}: &Stage{
						name:    StageName{val: "stage-3:aux-merge"},
						sType:   StageTypeMerge,
						address: "method-3/*/*",
						desc:    testSplitAndMergeStage3Method{},
						inputs: []*Link{
							{
								name:   LinkName{val: "1-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-1:aux-split"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: message.Field("field1"),
								},
							},
							{
								name:   LinkName{val: "2-to-3"},
								source: &LinkEndpoint{stage: StageName{val: "stage-2"}},
								target: &LinkEndpoint{
									stage: StageName{val: "stage-3:aux-merge"},
									field: message.Field("field2"),
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
						name:    StageName{val: "stage-3"},
						sType:   StageTypeUnary,
						address: "method-3/*/*",
						desc:    testSplitAndMergeStage3Method{},
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
						name:    StageName{val: "stage-3:aux-sink"},
						sType:   StageTypeSink,
						address: "method-3/*/*",
						desc:    testSplitAndMergeStage3Method{},
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testSplitAndMergeStage1Method{},
					"method-2/*/*": testSplitAndMergeStage2Method{},
					"method-3/*/*": testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := NewContext(tc.resolver)
			output, err := New(ctx, tc.input)
			if err != nil {
				t.Fatalf("new error: %s", err)
			}
			cmpOpts := cmp.AllowUnexported(
				Pipeline{},
				PipelineName{},
				Stage{},
				StageName{},
				Link{},
				LinkName{},
				LinkEndpoint{},
			)
			if diff := cmp.Diff(tc.expected, output, cmpOpts); diff != "" {
				t.Fatalf("output mismatch:\n%s", diff)
			}
		})
	}
}

func TestNewIsErr(t *testing.T) {
	tests := map[string]struct {
		input       *api.Pipeline
		validateErr func(err error) string
		resolver    method.ResolveFunc
	}{
		"empty pipeline name": {
			input: &api.Pipeline{},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyPipelineName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEmptyPipelineName, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				t.Fatalf("No such method: %s", address)
				return nil, nil
			},
		},
		"empty stage name": {
			input: &api.Pipeline{
				Name:   "Pipeline",
				Stages: []*api.Stage{{Address: "method"}},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyStageName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEmptyStageName, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				t.Fatalf("No such method: %s", address)
				return nil, nil
			},
		},
		"empty link name": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
					{SourceStage: "stage-1", TargetStage: "stage-2"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyLinkName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEmptyLinkName, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"empty link source name": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
					{Name: "1-to-2", TargetStage: "stage-2"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptySourceName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEmptySourceName, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"empty link target name": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
					{Name: "1-to-2", SourceStage: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEmptyTargetName) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEmptyTargetName, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"equal link source and target": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
					{Name: "1-to-2", SourceStage: "stage-1", TargetStage: "stage-1"},
				},
			},
			validateErr: func(err error) string {
				if !errors.Is(err, errEqualSourceAndTarget) {
					format := "error mismatch: expected %s, received %s"
					return fmt.Sprintf(format, errEqualSourceAndTarget, err)
				}
				return ""
			},
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"source does not exist": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"target does not exist": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"new link set full message": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
					{Name: "stage-3", Address: "method-3"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
					"method-3/*/*": testLinearStage3Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"old link set full message": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
					{Name: "stage-3", Address: "method-3"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
					"method-3/*/*": testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"new and old links set same": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
					{Name: "stage-3", Address: "method-3"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
					"method-3/*/*": testSplitAndMergeStage3Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
		"incompatible message descriptor": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{Name: "stage-1", Address: "method-1"},
					{Name: "stage-2", Address: "method-2"},
				},
				Links: []*api.Link{
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
			resolver: func(_ context.Context, address string) (method.Desc, error) {
				mapper := map[string]method.Desc{
					"method-1/*/*": testLinearStage1Method{},
					"method-2/*/*": testLinearStage2Method{},
				}
				s, ok := mapper[address]
				if !ok {
					panic(fmt.Sprintf("No such method: %v", address))
				}
				return s, nil
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := NewContext(tc.resolver)
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

func (m testLinearStage1Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testLinearStage1Method) Input() message.Type {
	return testEmptyDesc{}
}

func (m testLinearStage1Method) Output() message.Type {
	return testOuterValDesc{}
}

type testLinearStage2Method struct{}

func (m testLinearStage2Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testLinearStage2Method) Input() message.Type {
	return testOuterValDesc{}
}

func (m testLinearStage2Method) Output() message.Type {
	return testOuterValDesc{}
}

type testLinearStage3Method struct{}

func (m testLinearStage3Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testLinearStage3Method) Input() message.Type {
	return testOuterValDesc{}
}

func (m testLinearStage3Method) Output() message.Type {
	return testEmptyDesc{}
}

type testSplitAndMergeStage1Method struct{}

func (m testSplitAndMergeStage1Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testSplitAndMergeStage1Method) Input() message.Type {
	return testEmptyDesc{}
}

func (m testSplitAndMergeStage1Method) Output() message.Type {
	return testInnerValDesc{}
}

type testSplitAndMergeStage2Method struct{}

func (m testSplitAndMergeStage2Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testSplitAndMergeStage2Method) Input() message.Type {
	return testInnerValDesc{}
}

func (m testSplitAndMergeStage2Method) Output() message.Type {
	return testInnerValDesc{}
}

type testSplitAndMergeStage3Method struct{}

func (m testSplitAndMergeStage3Method) Dial() (method.Conn, error) {
	return nil, nil
}

func (m testSplitAndMergeStage3Method) Input() message.Type {
	return testOuterValDesc{}
}

func (m testSplitAndMergeStage3Method) Output() message.Type {
	return testEmptyDesc{}
}

type testEmptyDesc struct{}

func (d testEmptyDesc) Compatible(other message.Type) bool {
	_, ok := other.(testEmptyDesc)
	return ok
}

func (d testEmptyDesc) Subfield(f message.Field) (message.Type, error) {
	panic("method get field should not be called for testEmptyDesc")
}

func (d testEmptyDesc) Build() message.Instance { panic("called build method") }

// Represents a descriptor of a message with two inner fields: field1 and field2.
// Each field is associated with a descriptor of type testInnerValDesc
type testOuterValDesc struct{}

func (d testOuterValDesc) Compatible(other message.Type) bool {
	_, ok := other.(testOuterValDesc)
	return ok
}

func (d testOuterValDesc) Subfield(f message.Field) (message.Type, error) {
	switch f {
	case "field1", "field2":
		return testInnerValDesc{}, nil
	default:
		panic(fmt.Sprintf("Unknown field for testOuterValDesc: %s", string(f)))
	}
}

func (d testOuterValDesc) Build() message.Instance { panic("called build method") }

func (d testOuterValDesc) String() string {
	return "testOuterValDesc"
}

type testInnerValDesc struct{}

func (d testInnerValDesc) Compatible(other message.Type) bool {
	_, ok := other.(testInnerValDesc)
	return ok
}

func (d testInnerValDesc) Subfield(f message.Field) (message.Type, error) {
	panic("method get field should not be called for testInnerValDesc")
}

func (d testInnerValDesc) Build() message.Instance { panic("called build method") }

func (d testInnerValDesc) String() string {
	return "testInnerValDesc"
}
