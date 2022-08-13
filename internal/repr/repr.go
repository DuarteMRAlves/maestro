package repr

import (
	"fmt"
	"strings"

	"github.com/DuarteMRAlves/maestro/internal/api"
)

func Pipeline(p *api.Pipeline) string {
	var builder strings.Builder
	writePipeline(&builder, p)
	return builder.String()
}

func writePipeline(b *strings.Builder, p *api.Pipeline) {
	if p == nil {
		b.WriteString("nil")
		return
	}

	writeStringField(b, "Name", p.Name, 0)

	b.WriteString("\nExecutionMode: ")
	writeExecutionMode(b, p.Mode)

	b.WriteString("\nStages: ")
	writeStages(b, p.Stages...)

	b.WriteString("\nLinks: ")
	writeLinks(b, p.Links...)
}

func writeExecutionMode(b *strings.Builder, m api.ExecutionMode) {
	switch m {
	case api.OfflineExecution:
		b.WriteString("Offline")
	case api.OnlineExecution:
		b.WriteString("Online")
	default:
		b.WriteString("Unknown")
	}
}

func writeStages(b *strings.Builder, ss ...*api.Stage) {
	b.WriteRune('[')
	if len(ss) == 0 {
		b.WriteRune(']')
		return
	}
	indent := uint(1)
	b.WriteRune('\n')
	for _, s := range ss {
		writeIdent(b, indent)
		b.WriteString("{\n")
		writeStage(b, s, indent+1)
		b.WriteRune('\n')
		writeIdent(b, indent)
		b.WriteRune('}')
		b.WriteRune('\n')
	}
	b.WriteRune(']')
}

func writeStage(b *strings.Builder, s *api.Stage, indent uint) {
	writeStringField(b, "Name", s.Name, indent)
	b.WriteRune('\n')
	writeStringField(b, "Address", s.Address, indent)
	b.WriteRune('\n')
	writeStringField(b, "Service", s.Service, indent)
	b.WriteRune('\n')
	writeStringField(b, "Method", s.Method, indent)
}

func writeLinks(b *strings.Builder, ll ...*api.Link) {
	b.WriteRune('[')
	if len(ll) == 0 {
		b.WriteRune(']')
		return
	}
	indent := uint(1)
	b.WriteRune('\n')
	for _, l := range ll {
		writeIdent(b, indent)
		b.WriteString("{\n")
		writeLink(b, l, indent+1)
		b.WriteRune('\n')
		writeIdent(b, indent)
		b.WriteRune('}')
		b.WriteRune('\n')
	}
	b.WriteRune(']')
}

func writeLink(b *strings.Builder, l *api.Link, indent uint) {
	writeStringField(b, "Name", l.Name, indent)
	b.WriteRune('\n')
	writeStringField(b, "SourceStage", l.SourceStage, indent)
	b.WriteRune('\n')
	writeStringField(b, "SourceField", l.SourceField, indent)
	b.WriteRune('\n')
	writeStringField(b, "TargetStage", l.TargetStage, indent)
	b.WriteRune('\n')
	writeStringField(b, "TargetField", l.TargetField, indent)
	b.WriteRune('\n')
	writeUIntegerField(b, "NumEmptyMessages", l.NumEmptyMessages, indent)
}

func writeStringField(b *strings.Builder, field, val string, indent uint) {
	writeIdent(b, indent)
	v := fmt.Sprintf("%s: %q", field, val)
	b.WriteString(v)
}

func writeUIntegerField(b *strings.Builder, field string, val uint, indent uint) {
	writeIdent(b, indent)
	v := fmt.Sprintf("%s: %d", field, val)
	b.WriteString(v)
}

func writeIdent(b *strings.Builder, indent uint) {
	for i := uint(0); i < indent; i++ {
		b.WriteRune('\t')
	}
}
